--multiline
CREATE PROCEDURE `store_profile_data`(IN p_profile_ident int, IN p_data_element int, IN p_datavalue MEDIUMTEXT, IN p_updated_by varchar(255), IN p_approved int, IN p_overriden int, IN p_is_encrypted tinyint(1))
BEGIN
    DECLARE v_currentVersion int;
    SET v_currentVersion = (SELECT MAX(version) FROM profile_data WHERE `data_element_id` = p_data_element AND `profile_id` = p_profile_ident);

    IF v_currentVersion IS NULL THEN
        set v_currentVersion = 0;
END IF;

    IF (SELECT de.datatype FROM data_element de WHERE de.data_element_id = p_data_element) = 'BOOLEAN' THEN
SELECT LOWER(p_datavalue) INTO p_datavalue;
END IF;

INSERT INTO profile_data (profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden, is_encrypted)
VALUES (p_profile_ident, p_data_element, p_datavalue, v_currentVersion + 1, current_timestamp, p_updated_by, current_timestamp, p_updated_by, p_approved, p_overriden, p_is_encrypted);
END