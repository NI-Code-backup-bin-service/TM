--multiline
CREATE PROCEDURE `updateMpans`()
BEGIN
	# Temp table to store the data.
	DROP TEMPORARY TABLE IF EXISTS visa_qr_temp;
	CREATE TEMPORARY TABLE visa_qr_temp (merchantNo VARCHAR(12), siteId INT(11), siteProfileId INT(11), currentActiveModulesValue TEXT, mpan VARCHAR(16), categoryCode VARCHAR(255));

	SET @siteProfileTypeId = (SELECT profile_type_id FROM profile_type WHERE `name` = 'site');

	SET @merchantNoDataElementId = (SELECT data_element_id FROM data_element WHERE `name` = 'merchantNo');
	SET @visaQrDataGroupId = (SELECT data_group_id FROM data_group WHERE `name` = 'visaQr');
	SET @mpanDataElementId = (SELECT data_element_id FROM data_element WHERE `name` = 'mpan' AND data_group_id = @visaQrDataGroupId);

	# Insert what we've been given.
	INSERT INTO visa_qr_temp (merchantNo, siteId, siteProfileId)
	VALUES  ('999888777999', NULL, NULL);

	# Work out site IDs from merchant number data element entries.
	UPDATE visa_qr_temp
	SET siteId = (SELECT sp.site_id
				  FROM site_profiles sp
				  JOIN `profile` p ON sp.profile_id = p.profile_id AND p.profile_type_id = @siteProfileTypeId
				  JOIN profile_data pd ON p.profile_id = pd.profile_id AND pd.data_element_id = @merchantNoDataElementId
				  WHERE pd.datavalue = merchantNo
				  ORDER BY sp.updated_at DESC
				  LIMIT 1);
				  
	# Work out profile IDs from site IDs.
	UPDATE visa_qr_temp
	SET siteProfileId = (SELECT sp.profile_id
						 FROM site_profiles sp
						 JOIN `profile` p ON sp.profile_id = p.profile_id AND p.profile_type_id = @siteProfileTypeId
						 WHERE site_id = siteId
						 ORDER BY sp.updated_at DESC
						 LIMIT 1);

	# Log and then remove any sites which don't actually exist.
	SELECT * FROM visa_qr_temp WHERE siteId IS NULL OR siteProfileId IS NULL;
	DELETE FROM visa_qr_temp WHERE siteId IS NULL OR siteProfileId IS NULL;

	# Update the timestamp of the MPAN field so that config download will see this as new.
	UPDATE profile_data pd
	SET updated_at = CURRENT_TIMESTAMP,
		updated_by = 'system'
	WHERE pd.profile_id IN (SELECT siteProfileId FROM visa_qr_temp)
	AND data_element_id = @mpanDataElementId;

	DROP TEMPORARY TABLE IF EXISTS visa_qr_temp;
END