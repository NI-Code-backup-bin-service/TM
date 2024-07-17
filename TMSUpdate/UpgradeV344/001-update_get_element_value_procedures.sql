--multiline;
CREATE PROCEDURE `get_element_value`(IN profileID int, IN elementID int)
BEGIN
SELECT
    IFNULL((SELECT pd.datavalue
        FROM profile_data pd
        WHERE pd.profile_id = profileID
        AND pd.data_element_id = elementID),''),
    de.is_encrypted,
    de.is_password
    FROM data_element de
    WHERE de.data_element_id = elementID;
END;