--multiline;
CREATE PROCEDURE profile_data_fetch(IN profileId int)
BEGIN
select pdg.profile_id,
       dg.data_group_id,
       dg.displayname_en,
       de.data_element_id,
       de.name as 'data_element_name',
        dg.name as 'data_group',
        de.tid_overridable,
       de.datatype,
       de.tooltip,
       pd.datavalue,
       de.is_allow_empty,
       de.max_length,
       de.validation_expression,
       de.validation_message,
       de.front_end_validate,
       de.`options` as 'options',
        de.`displayname_en` as `display_name`,
       IFNULL(de.is_password, 0) AS is_password,
       pd.is_encrypted AS is_encrypted,
       de.sort_order_in_group,
       de.is_read_only_at_creation,
       de.required_at_acquirer_level,
       de.required_at_chain_level
from profile_data_group pdg
         join data_group dg on dg.data_group_id = pdg.data_group_id
         join data_element de on de.data_group_id = dg.data_group_id
         left join profile_data pd on pd.data_element_id = de.data_element_id and pd.profile_id = pdg.profile_id
         left join profile p ON p.profile_id = pd.profile_id
where pdg.profile_id = profileId order by de.sort_order_in_group;
END;