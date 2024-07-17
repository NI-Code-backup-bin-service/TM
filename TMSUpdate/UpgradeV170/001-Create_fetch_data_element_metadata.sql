--multiline
CREATE PROCEDURE `fetch_data_element_metadata`(IN dataElementId INT)
BEGIN
    SELECT
        data_element_id,
        name,
        datatype,
        is_allow_empty,
        max_length,
        validation_expression,
        validation_message,
        front_end_validate,
        `unique`,
        `options`,
        displayname_en,
        IFNULL(is_password, 0),
        is_encrypted,
        tooltip
    FROM data_element
    WHERE data_element_id = dataElementId;
END