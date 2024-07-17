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
            tooltip,
            file_max_size,
            file_min_ratio,
            file_max_ratio,
            is_read_only_at_creation,
            required_at_acquirer_level,
            required_at_chain_level
        FROM data_element
        WHERE data_element_id = dataElementId;
    END;