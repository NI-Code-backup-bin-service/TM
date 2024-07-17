# Profile to set data on
SET @p_profileId = 3;
#
# data to set on profile
SET @v_NolEnabled = 'true';
SET @v_AID = '784000';
SET @v_NolTaxiBin = '990098';
SET @v_NolMerBin = '978432';
SET @v_DriverId = '6372879';
SET @v_MaxCardBalLimitType1 = '500000';
SET @v_MaxCardBalLimitType2 = '100000';
SET @v_MaxCardBalLimitType4 = '500000';
SELECT data_group_id INTO @v_nolDataGroupId FROM data_group WHERE name = 'nol';
# Enable the nol
INSERT IGNORE INTO profile_data_group(profile_id, data_group_id, version, updated_at, updated_by, created_at, created_by) VALUES (@p_profileId, @v_nolDataGroupId, 1, NOW(), 'system', NOW(), 'system');
# Set the data
INSERT IGNORE INTO profile_data (profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden, is_encrypted, not_overridable) SELECT @p_profileId, de.data_element_id, CASE WHEN de.name = 'AID' THEN @v_AID WHEN de.name = 'driverId' THEN @v_DriverId WHEN de.name = 'maxCardBalLimitType1' THEN @v_MaxCardBalLimitType1 WHEN de.name = 'maxCardBalLimitType2' THEN @v_MaxCardBalLimitType2 WHEN de.name = 'maxCardBalLimitType4' THEN @v_MaxCardBalLimitType4 WHEN de.name = 'nolEnabled' THEN @v_NolEnabled WHEN de.name = 'nolMerBin' THEN @v_NolMerBin WHEN de.name = 'nolTaxiBin' THEN @v_NolTaxiBin END, 1, NOW(), 'system', NOW(), 'system', 1, 0, 0, 0 FROM data_element de WHERE de.data_group_id = @v_nolDataGroupId;