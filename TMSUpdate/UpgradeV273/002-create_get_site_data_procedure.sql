--multiline;
CREATE PROCEDURE `get_site_data`(IN profileID INT, IN siteID INT)
BEGIN
    SELECT distinct e.site_id,
        e.data_group_id,
        dg.name,
        dg.displayname_en,
        e.data_element_id,
        e.name,
        de.tooltip,
        sd.level,
        sd.priority,
        sd.datavalue,
        sd.overriden,
        sd.not_overridable,
        e.datatype,
        e.is_allow_empty,
        e.max_length,
        e.validation_expression,
        e.validation_message,
        e.front_end_validate,
        e.options,
        e.display_name,
        e.sort_order_in_group,
        IFNULL(e.is_password, 0),
        sd.is_encrypted,
        e.file_max_size,
        e.file_min_ratio,
        e.file_max_ratio,
        de.tid_overridable,
        de.is_read_only_at_creation,
        de.required_at_acquirer_level,
        de.required_at_chain_level
    FROM site_data_elements e
    JOIN data_element de
        ON de.data_element_id = e.data_element_id
    INNER JOIN data_group dg
        ON dg.data_group_id = e.data_group_id
    LEFT JOIN site_data sd
        ON sd.site_id = e.site_id
        AND sd.data_element_id = e.data_element_id
    WHERE e.site_id = siteID
    AND e.data_group_id IN (SELECT data_group_id FROM profile_data_group where profile_id = profileID)
    AND e.location_name IN ('site_configuration','tid_override') order by e.sort_order_in_group;
END;