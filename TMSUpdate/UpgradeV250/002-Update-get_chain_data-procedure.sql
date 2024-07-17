--multiline;
CREATE PROCEDURE get_chain_data(IN profileId int, IN acquirerID int)
BEGIN
SELECT
    de.data_group_id,
    dg.name as 'data_group',
    dg.displayname_en,
    de.data_element_id,
    de.name as 'data_element_name',
    de.tooltip,
    source,
    datavalue,
    cd.overriden,
    de.datatype,
    de.is_allow_empty,
    de.max_length,
    de.validation_expression,
    de.validation_message,
    de.front_end_validate,
    de.`options` as 'options',
    de.sort_order_in_group,
    de.`displayname_en` as `display_name`,
    cd.is_encrypted,
    de.is_password,
    IFNULL(cd.not_overridable, 0),
    de.required_at_acquirer_level,
    de.required_at_chain_level
FROM chain_data cd
         JOIN data_element de ON cd.data_element_id = de.data_element_id
         JOIN data_group dg ON dg.data_group_id = de.data_group_id
         JOIN profile_data_group pdg ON pdg.data_group_id = dg.data_group_id AND (pdg.profile_id = profileId OR pdg.profile_id = acquirerId)
WHERE cd.profile_id = profileId;
END;