--multiline
CREATE PROCEDURE `get_element_value`(
    IN profile_id int,
    IN element_id int
)
BEGIN
    SELECT pd.datavalue, pd.is_encrypted, (EXISTS (SELECT de.data_element_id FROM data_element de WHERE de.`is_password` = 1 AND de.`data_element_id` = element_id)) AS isPassword
    FROM profile_data pd
    WHERE pd.data_element_id = element_id AND
            pd.profile_id = profile_id
    ORDER BY pd.version DESC
    LIMIT 1;
END