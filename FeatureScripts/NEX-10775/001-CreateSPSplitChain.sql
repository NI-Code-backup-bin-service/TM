--multiline
CREATE PROCEDURE split_chain(IN merchantID varchar(255), IN chainProfileID INT, IN newChainProfileID INT)
    sp : BEGIN
            DECLARE chainType mediumtext;
            DECLARE profileID, siteID INT;

            SET @chainType= (SELECT name FROM profile_type
                                WHERE profile_type_id = (SELECT profile_type_id FROM profile
                                                            WHERE profile_id = chainProfileID));

            IF (@chainType != 'chain') THEN
                SELECT 'Chain Profile ID is not valid' As StatusInfo, merchantID, chainProfileID, newChainProfileID;
                LEAVE sp;
            END IF;

            SET @chainType = (SELECT name FROM profile_type
                                WHERE profile_type_id = (SELECT profile_type_id FROM profile
                                                            WHERE profile_id = newChainProfileID));

            IF (@chainType != 'chain') THEN
                SELECT 'New Chain Profile ID is not valid' As StatusInfo, merchantID, chainProfileID, newChainProfileID;
                LEAVE sp;
            END IF;

            SET @profileID = (SELECT profile_id FROM profile_data
                                WHERE data_element_id = (SELECT data_element_id FROM data_element
                                                            WHERE name='merchantNo'
                                                            AND data_group_id=(SELECT data_element_id
                                                                                FROM data_group
                                                                                WHERE name='store'))
                                AND datavalue=merchantID);

            IF (@profileID = 0 OR @profileID IS NULL) THEN
                SELECT 'Merchant ID is not valid' As StatusInfo, merchantID, chainProfileID, newChainProfileID;
                LEAVE sp;
            END IF;

            SET @siteID = (SELECT site_id FROM site_profiles where profile_id=@profileID);

            IF ((SELECT COUNT(*) FROM site_profiles where profile_id=chainProfileID AND site_id=@siteID) = 0) THEN
                SELECT 'merchant ID does not belongs to Chain Profile ID.-->' As StatusInfo, merchantID, chainProfileID, newChainProfileID;
                LEAVE sp;
            END IF;

            UPDATE site_profiles SET profile_id=newChainProfileID WHERE profile_id=chainProfileID AND site_id=@siteID;

            SELECT 'Successfully moved site to new chain' As StatusInfo, merchantID, chainProfileID, newChainProfileID;
END