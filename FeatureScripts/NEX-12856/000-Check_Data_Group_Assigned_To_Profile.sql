--multiline
CREATE PROCEDURE `Check_Data_Group_Assigned_To_Profile`(in dataGroupName text, MID text, chainName text, acquireName text, contactApplicationConfigs mediumtext, contactApplicationConfigShared mediumtext, contactKeyConfigs mediumtext, ctlsApplicationConfigs mediumtext, ctlsApplicationConfigShared mediumtext)
BEGIN
    SET @ContactApplicationConfigs = "contactApplicationConfigs";
    SET @ContactApplicationConfigShared = "contactApplicationConfigShared";
    SET @ContactKeyConfigs = "contactKeyConfigs";
    SET @CtlsApplicationConfigs = "ctlsApplicationConfigs";
    SET @CtlsApplicationConfigShared = "ctlsApplicationConfigShared";
    SET @profile_id = "";

    IF (MID IS NOT NULL) THEN
        SET @profile_id = (Select profile_id from profile_data pd where pd.data_element_id = (Select data_element_id from data_element where name = "merchantNo") AND
            pd.datavalue = MID);
    END IF;
    IF (chainName IS NOT NULL) THEN
        SET @profile_id = (Select profile_id from profile pd where binary pd.name = chainName AND pd.profile_type_id = (Select profile_type_id from profile_type where name = "chain"));
    END IF;
    IF (acquireName IS NOT NULL) THEN
        SET @profile_id = (Select profile_id from profile pd where binary pd.name = acquireName AND pd.profile_type_id = (Select profile_type_id from profile_type where name = "acquirer"));
    END IF;

    SET @dataGroupExists := EXISTS(Select * FROM profile_data_group WHERE profile_id = @profile_id AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = dataGroupName));

    # DataGroup is enabled then only updated the respective values.
    IF @dataGroupExists = 1 THEN
        # Set the ContactApplicationConfigs value
        IF (contactApplicationConfigs IS NOT NULL) THEN
            INSERT into profile_data (profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden, is_encrypted)
            values (@profile_id, (SELECT data_element_id FROM data_element WHERE name = @ContactApplicationConfigs AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = dataGroupName)), contactApplicationConfigs, 1, NOW(), 'system', NOW(), 'system', 1, 1, 0)
            ON DUPLICATE KEY
                UPDATE datavalue=VALUES(datavalue), updated_at = VALUES(updated_at), updated_by = VALUES(updated_by);
        END IF;

        # Set the ContactApplicationConfigShared value
        IF (contactApplicationConfigShared IS NOT NULL) THEN
            INSERT into profile_data (profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden, is_encrypted)
            values (@profile_id, (SELECT data_element_id FROM data_element WHERE name = @ContactApplicationConfigShared AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = dataGroupName)), contactApplicationConfigShared, 1, NOW(), 'system', NOW(), 'system', 1, 1, 0)
            ON DUPLICATE KEY
                UPDATE datavalue=VALUES(datavalue), updated_at = VALUES(updated_at), updated_by = VALUES(updated_by);
        END IF;

        # Set the ContactKeyConfigs value
        IF (contactKeyConfigs IS NOT NULL) THEN
            INSERT into profile_data (profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden, is_encrypted)
            values (@profile_id, (SELECT data_element_id FROM data_element WHERE name = @ContactKeyConfigs AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = dataGroupName)), contactKeyConfigs, 1, NOW(), 'system', NOW(), 'system', 1, 1, 0)
            ON DUPLICATE KEY
                UPDATE datavalue=VALUES(datavalue), updated_at = VALUES(updated_at), updated_by = VALUES(updated_by);
        END IF;

        # Set the CtlsApplicationConfigs value
        IF (ctlsApplicationConfigs IS NOT NULL) THEN
            INSERT into profile_data (profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden, is_encrypted)
            values (@profile_id, (SELECT data_element_id FROM data_element WHERE name = @CtlsApplicationConfigs AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = dataGroupName)), ctlsApplicationConfigs, 1, NOW(), 'system', NOW(), 'system', 1, 1, 0)
            ON DUPLICATE KEY
                UPDATE datavalue=VALUES(datavalue), updated_at = VALUES(updated_at), updated_by = VALUES(updated_by);
        END IF;

        # Set the CtlsApplicationConfigShared value
        IF (ctlsApplicationConfigShared IS NOT NULL) THEN
            INSERT into profile_data (profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden, is_encrypted)
            values (@profile_id, (SELECT data_element_id FROM data_element WHERE name = @CtlsApplicationConfigShared AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = dataGroupName)), ctlsApplicationConfigShared, 1, NOW(), 'system', NOW(), 'system', 1, 1, 0)
            ON DUPLICATE KEY
                UPDATE datavalue=VALUES(datavalue), updated_at = VALUES(updated_at), updated_by = VALUES(updated_by);
        END IF;
    END IF;
END