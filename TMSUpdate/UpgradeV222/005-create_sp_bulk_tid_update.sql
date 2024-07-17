--multiline;
CREATE PROCEDURE bulk_tid_profile_data_update(IN profileID int,
    IN dataElementID int,
    IN newdataValue MEDIUMTEXT,
    IN updated_by_user varchar(255),
    IN is_value_encrypted BOOLEAN)
BEGIN
    UPDATE profile_data pd
    SET pd.datavalue = newdataValue, pd.updated_at=NOW(), pd.updated_by = updated_by_user, pd.is_encrypted = is_value_encrypted
    WHERE pd.profile_id = profileID
      AND pd.data_element_id = dataElementID;
END