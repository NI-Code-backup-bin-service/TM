--multiline;
CREATE PROCEDURE bulk_tid_profile_data_update(IN profileID int,
                                              IN dataElementID int,
                                              IN newdataValue MEDIUMTEXT,
                                              IN updated_by_user varchar(255),
                                              IN is_value_encrypted BOOLEAN)
BEGIN
    SET @dataElementExists = 0;
    SET @datagroupExists =  -1;
    SET @datagroupId = (select data_group_id from data_element where data_element_id = dataElementID);
    SELECT IFNULL(profile_id, 0) into @dataElementExists from profile_data WHERE profile_id = profileID AND data_element_id = dataElementID;
    SELECT IFNULL(profile_id, -1) into @datagroupExists from profile_data_group WHERE profile_id = profileID AND data_group_id = @datagroupId;

    IF (@dataElementExists > 0)
        THEN
        UPDATE profile_data pd
        SET pd.datavalue = newdataValue, pd.updated_at=NOW(), pd.updated_by = updated_by_user, pd.is_encrypted = is_value_encrypted
        WHERE pd.profile_id = profileID AND pd.data_element_id = dataElementID;
    ELSE
        INSERT INTO profile_data (profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden, is_encrypted)
        VALUES (profileID, dataElementID, newdataValue, 1, current_timestamp, updated_by_user, current_timestamp, updated_by_user, 1, 1, is_value_encrypted);
    END IF;
    IF (@datagroupExists < 0)
        THEN
        INSERT INTO profile_data_group(
            profile_id,data_group_id,version,updated_at,updated_by,created_at,created_by)
        VALUES (
            profileID,
            @datagroupId,
            1,
            current_timestamp,
            updated_by_user,
            current_timestamp,
            updated_by_user);
    END IF;
END;