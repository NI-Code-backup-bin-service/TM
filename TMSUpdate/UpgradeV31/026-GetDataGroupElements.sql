--multiline
CREATE PROCEDURE `get_data_group_elements`(IN dataGroupId INT)
BEGIN
    select
        dg.data_group_id,
        dg.name as 'data_group',
        e.data_element_id,
        e.name,
        e.datatype,
        e.max_length,
        e.validation_expression,
        e.validation_message,
        e.front_end_validate,
        e.options,
        e.is_password,
        e.is_encrypted
    from data_group dg
             left join data_element e ON e.data_group_id = dg.data_group_id
    WHERE dg.data_group_id = dataGroupId;
END