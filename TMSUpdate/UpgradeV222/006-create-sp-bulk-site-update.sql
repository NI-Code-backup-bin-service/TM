--multiline;
CREATE PROCEDURE bulk_site_profile_data_update(IN profileID int,
                                 IN dataElementID int,
                                 IN newdataValue MEDIUMTEXT,
                                 IN updated_by_user varchar(255),
                                 IN is_value_encrypted BOOLEAN)
BEGIN

DECLARE siteId ,overr INT;

IF EXISTS (SELECT profile_data_id FROM profile_data WHERE profile_id = profileID AND data_element_id = dataElementID) THEN
    SET @siteId = (select site_id from site_profiles where profile_id = profileID);
    SET @overr = (SELECT overriden FROM profile_data WHERE profile_id = profileID AND data_element_id = dataElementID);
    IF (newDataValue = "NULL") THEN
        IF (@overr = 1) THEN
            DELETE FROM profile_data
            WHERE profile_id = profileID
              AND data_element_id = dataElementID;
            UPDATE site s
            LEFT JOIN site_profiles sp
                ON sp.site_id = s.site_id
            SET s.updated_at = NOW()
            WHERE sp.profile_id = profileID;
        ELSE
            UPDATE profile_data
            SET datavalue = "", updated_at=NOW(), updated_by = updated_by_user
            WHERE profile_id = profileID
              AND data_element_id = dataElementID;
        END IF;
    ELSE
    UPDATE profile_data pd
    SET pd.datavalue = newdataValue, pd.updated_at=NOW(), pd.updated_by = updated_by_user, pd.is_encrypted = is_value_encrypted
    WHERE pd.profile_id = profileID
      AND pd.data_element_id = dataElementID;
    END IF;

ELSE
    INSERT INTO profile_data (profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden, is_encrypted)
    VALUES (profileID,
            dataElementID,
            newdataValue,
            1,
            current_timestamp,
            updated_by_user,
            current_timestamp,
            updated_by_user,
            1,
            1,
            is_value_encrypted);
END IF;

END