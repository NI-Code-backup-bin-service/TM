--multiline
create procedure get_profile_data_for_tab_by_profile_id(IN p_profile_id int, IN p_tab_name text)
BEGIN
    SELECT
        JSON_ARRAYAGG(data_groups.JsonObject)
    FROM (
             SELECT
                 JSON_OBJECT(
                         'DataGroupID', dg.data_group_id,
                         'DataGroup', dg.name,
                         'DisplayName', dg.displayname_en,
                         'DataElements',
                         JSON_ARRAYAGG(
                                 JSON_OBJECT(
                                         'ElementId', de.data_element_id,
                                         'Name', de.name,
                                         'Type', de.datatype,
                                         'IsAllowedEmpty', if(de.is_allow_empty = 0 OR de.is_allow_empty = '' OR de.is_allow_empty IS NULL, CAST(false AS JSON), CAST(true AS JSON)),
                                         'DataValue', if(pd.datavalue = '' OR pd.datavalue IS NULL, get_parent_datavalue(de.data_element_id,  p.profile_id), pd.datavalue),
                                         'MaxLength', de.max_length,
                                         'ValidationExpression', de.validation_expression,
                                         'ValidationMessage', de.validation_message,
                                         'FrontEndValidate', if(de.front_end_validate = 0 OR de.front_end_validate = '', CAST(false AS JSON), CAST(true AS JSON)),
                                         'Unique', if(de.`unique` = 0 OR de.`unique` = '' OR de.`unique` IS NULL, CAST(false AS JSON), CAST(true AS JSON)),
                                         'Overriden',  CAST(false AS JSON),#if(pd.overriden = 0 OR pd.overriden = '' OR pd.overriden IS NULL, CAST(true AS JSON), CAST(false AS JSON)),
                                         'DisplayName', de.displayname_en,
                                         'IsPassword', if(de.is_password = 0 OR de.is_password = '' OR de.is_password IS NULL, CAST(false AS JSON), CAST(true AS JSON)),
                                         'IsEncrypted', if(de.is_encrypted = 0 OR de.is_encrypted = '' OR de.is_encrypted IS NULL, CAST(false AS JSON), CAST(true AS JSON)),
                                         'SortOrderInGroup', de.sort_order_in_group
                                     )
                             )
                     ) As "JsonObject"
             FROM profile p
                      INNER JOIN profile_type pt ON
                     p.profile_type_id = pt.profile_type_id
                      INNER JOIN data_element_locations del
                                 ON pt.profile_type_id = del.profile_type_id
                      LEFT JOIN data_element_locations_data_element delde
                                ON del.location_id = delde.location_id
                      INNER JOIN data_element de ON
                     delde.data_element_id = de.data_element_id
                      INNER JOIN data_group dg ON
                     de.data_group_id = dg.data_group_id
                      LEFT JOIN profile_data pd on
                         pd.data_element_id = de.data_element_id
                     AND
                         pd.profile_id = p.profile_id
             WHERE
                     p.profile_id = p_profile_id
               AND
                     del.location_name = p_tab_name
             GROUP BY dg.data_group_id) data_groups;
END;