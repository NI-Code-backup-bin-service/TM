--multiline
CREATE PROCEDURE get_site_tids(IN site_id int)
BEGIN
SELECT t.tid_id,
       t.serial,
       t.PIN,
       t.ExpiryDate,
       t.reset_pin,
       t.reset_pin_expiry_date,
       t.ActivationDate,
       (SELECT pd.datavalue FROM profile_data pd WHERE pd.data_element_id = 1 and  pd.profile_id = tp2.profile_id) as 'merchant_id',
        (SELECT COUNT(*) FROM tid_user_override tuo WHERE tuo.tid_id = t.tid_id) as 'userOverrides',
        IFNULL(ts.tid_profile_id,0) as 'tidProfileID',
        (SELECT COUNT(*) FROM velocity_limits tvl WHERE tvl.tid_id = t.tid_id) as 'fraudOverride',
        IFNULL(dg.data_group_id,0) as 'data_group_id',
        IFNULL(dg.displayname_en,'') as 'displayname_en',
        IFNULL(de.data_element_id ,0) as 'data_element_id',
        IFNULL(de.name ,'')  as 'data_element_name',
        IFNULL(dg.name ,'')  as 'data_group_name',
        IFNULL(de.tid_overridable ,false)  as 'tid_overridable',
        IFNULL(de.datatype ,'') as 'datatype',
       de.tooltip as 'displayname_en',
       pd.datavalue,
       IFNULL(de.is_allow_empty ,false) as 'is_allow_empty',
       de.max_length,
       de.validation_expression,
       de.validation_message,
       IFNULL(de.front_end_validate, 0) as 'front_end_validate',
       de.options as 'options',
       IFNULL(de.displayname_en ,'') as `display_name`,
       IFNULL(de.is_password, 0) AS is_password,
       pd.is_encrypted AS is_encrypted,
       IFNULL(de.sort_order_in_group ,false)  as 'sort_order_in_group',
       IFNULL(de.is_read_only_at_creation  ,false) as 'is_read_only_at_creation',
       IFNULL(de.required_at_acquirer_level ,false) as 'required_at_acquirer_level',
       IFNULL(de.required_at_chain_level ,false)  as 'required_at_chain_level'
FROM tid t
         LEFT JOIN tid_site ts ON ts.tid_id = t.tid_id
         LEFT JOIN site s on s.site_id = ts.site_id
         LEFT JOIN (site_profiles tp2
    join profile p2 on p2.profile_id = tp2.profile_id
    join profile_type pt2 on pt2.profile_type_id = p2.profile_type_id and pt2.priority = 2) on tp2.site_id = s.site_id
         LEFT JOIN profile_data_group pdg ON pdg.profile_id = ts.tid_profile_id
         LEFT join data_group dg on dg.data_group_id = pdg.data_group_id
         LEFT join data_element de on de.data_group_id = dg.data_group_id
         left join profile_data pd on pd.data_element_id = de.data_element_id and pd.profile_id = pdg.profile_id
         left join profile p ON p.profile_id = pd.profile_id
    AND pd.version = (SELECT MAX(d.version) FROM profile_data d WHERE d.data_element_id = 1 AND d.profile_id = p2.profile_id AND d.approved = 1)
WHERE ts.site_id = site_id;
END;