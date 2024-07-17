--multiline;
CREATE PROCEDURE get_data_elements_for_new_site(IN acquirer_id int, IN chain_id int)
BEGIN
SELECT dg.data_group_id,
       dg.name AS 'data_group',
        dg.displayname_en,
       de.data_element_id,
       de.tooltip,
       de.name,
       IFNULL(pt.name, 'site') AS `source`,
       pd.datavalue,
       de.datatype,
       de.is_allow_empty,
       de.options,
       de.displayname_en,
       IFNULL(de.is_password, 0),
       pd.is_encrypted,
       de.sort_order_in_group,
       de.required_at_site_level,
       IFNULL(pd.not_overridable, 0),
       de.is_read_only_at_creation
FROM data_group dg
         LEFT JOIN profile_data_group pdg
                   ON pdg.profile_id IN (acquirer_id, chain_id, 1) AND dg.data_group_id = pdg.data_group_id
         LEFT JOIN profile p ON pdg.profile_id = p.profile_id
         LEFT JOIN profile_type pt ON p.profile_type_id = pt.profile_type_id
         INNER JOIN data_element de ON de.data_group_id = dg.data_group_id
         INNER JOIN data_element_locations_data_element delde ON de.data_element_id = delde.data_element_id
         INNER JOIN data_element_locations del ON delde.location_id = del.location_id
         LEFT JOIN profile_data pd ON pd.data_element_id = de.data_element_id AND pd.profile_id = pdg.profile_id
WHERE del.location_name = 'site_configuration' order by de.sort_order_in_group;
END;