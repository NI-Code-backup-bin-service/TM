-- Set this to the profile (Site, Acquirer, Chain) you wish to update
SET @profile_id = 4;
-- Add Stored Procedure for this script
CREATE PROCEDURE `Check_Data_Group_Assigned_To_Profile_With_Merchant_UAPM_Name`(in profileIdValue int, dataGroupName text, dataElementName text, dataElementValue text) BEGIN SET @dataGroupExists := EXISTS(Select * FROM profile_data_group WHERE profile_id = profileIdValue AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = dataGroupName)); IF @dataGroupExists = 1 THEN INSERT into profile_data (profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden, is_encrypted) values (profileIdValue, (SELECT data_element_id FROM data_element WHERE name = dataElementName AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = dataGroupName)), dataElementValue, 1, NOW(), 'NISuper', NOW(), 'NISuper', 1, 0, 0) ON DUPLICATE KEY UPDATE datavalue=VALUES(datavalue), updated_at = VALUES(updated_at); END IF; END
-- Merchant UAPM Name
CALL Check_Data_Group_Assigned_To_Profile_With_Merchant_UAPM_Name(@profile_id, 'pullpayments', 'merchant_endpoint_mapping', 'Testing here For by build UAPM')
-- Delete Stored Procedure
drop procedure if exists Check_Data_Group_Assigned_To_Profile_With_Merchant_UAPM_Name;