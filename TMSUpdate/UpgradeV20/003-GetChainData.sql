--multiline
CREATE PROCEDURE get_chain_data(
 IN $profile_id int
)
BEGIN
    SELECT
        de.data_group_id,
        dg.name as data_group,
        de.data_element_id,
        de.name,
        source,
        datavalue,
        overriden,
        de.datatype,
        de.is_allow_empty,
        de.max_length,
        de.validation_expression,
        de.validation_message,
        de.front_end_validate,
        de.options,
        de.displayname_en
    FROM chain_data cd
    JOIN data_element de
    ON cd.data_element_id = de.data_element_id
    JOIN data_group dg
    ON dg.data_group_id = de.data_group_id
    WHERE cd.profile_id = $profile_id;
END