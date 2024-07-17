# Temp table to store the data.
DROP TEMPORARY TABLE IF EXISTS visa_qr_temp;
CREATE TEMPORARY TABLE visa_qr_temp (merchantNo VARCHAR(12), siteId INT(11), siteProfileId INT(11), currentActiveModulesValue TEXT, mpan VARCHAR(16), categoryCode VARCHAR(255));

SET @siteProfileTypeId = (SELECT profile_type_id FROM profile_type WHERE `name` = 'site');
SET @chainProfileTypeId = (SELECT profile_type_id FROM profile_type WHERE `name` = 'chain');
SET @acquirerProfileTypeId = (SELECT profile_type_id FROM profile_type WHERE `name` = 'acquirer');
SET @globalProfileTypeId = (SELECT profile_type_id FROM profile_type WHERE `name` = 'global');

SET @merchantNoDataElementId = (SELECT data_element_id FROM data_element WHERE `name` = 'merchantNo');
SET @modulesDataGroupId = (SELECT data_group_id FROM data_group WHERE `name` = 'modules');
SET @activeDataElementId = (SELECT data_element_id FROM data_element WHERE `name` = 'active' AND data_group_id = @modulesDataGroupId);
SET @visaQrDataGroupId = (SELECT data_group_id FROM data_group WHERE `name` = 'visaQr');
SET @mpanDataElementId = (SELECT data_element_id FROM data_element WHERE `name` = 'mpan' AND data_group_id = @visaQrDataGroupId);
SET @categoryCodeDataElementId = (SELECT data_element_id FROM data_element WHERE `name` = 'categoryCode' AND data_group_id = @visaQrDataGroupId);

# Insert what we've been given.
INSERT INTO visa_qr_temp (merchantNo, siteId, siteProfileId, currentActiveModulesValue, mpan, categoryCode)
VALUES  ('200600002455', NULL, NULL, NULL, '4035560093518730', '5621'),
		('200600002604', NULL, NULL, NULL, '4035560061738583', '5532'),
		('200600002927', NULL, NULL, NULL, '4035560011386244', '7311'),
		('200600004287', NULL, NULL, NULL, '4035560012701110', '5814'),
		('200600005177', NULL, NULL, NULL, '4035560005465756', '8099'),
		('200600006290', NULL, NULL, NULL, '4035560002104739', '7549'),
		('200600006621', NULL, NULL, NULL, '4035560011462540', '5947'),
		('200600007264', NULL, NULL, NULL, '4035560065860284', '5099'),
		('200600007280', NULL, NULL, NULL, '4035560054473685', '5441'),
		('200600007538', NULL, NULL, NULL, '4035560042843262', '7297'),
		('200600008098', NULL, NULL, NULL, '4035560087576751', '5661'),
		('200600010367', NULL, NULL, NULL, '4035560015151909', '7298'),
		('200600011480', NULL, NULL, NULL, '4035560013589472', '5039'),
		('200600012223', NULL, NULL, NULL, '4035560046585026', '7298'),
		('200600012272', NULL, NULL, NULL, '4035560054022219', '7299'),
		('200600012363', NULL, NULL, NULL, '4035560091124416', '7298'),
		('200600012843', NULL, NULL, NULL, '4035560092437486', '5977'),
		('200600013304', NULL, NULL, NULL, '4035560041723366', '5814'),
		('200600014286', NULL, NULL, NULL, '4035560091674535', '6513'),
		('200600014377', NULL, NULL, NULL, '4035560068338197', '5814'),
		('200600014443', NULL, NULL, NULL, '4035560063611895', '7298'),
		('200600014930', NULL, NULL, NULL, '4035560071489615', '7298'),
		('200600014948', NULL, NULL, NULL, '4035560089343754', '5814'),
		('200600015325', NULL, NULL, NULL, '4035560031503216', '7298'),
		('200600015713', NULL, NULL, NULL, '4035560060908500', '5814'),
		('200600015796', NULL, NULL, NULL, '4035560093180853', '7298'),
		('200600015911', NULL, NULL, NULL, '4035560057803946', '8050'),
		('200600016133', NULL, NULL, NULL, '4035560087632919', '7298'),
		('200600016489', NULL, NULL, NULL, '4035560052688318', '7298'),
		('200600017198', NULL, NULL, NULL, '4035560058335195', '5533'),
		('200600017230', NULL, NULL, NULL, '4035560007503489', '5814'),
		('200600017719', NULL, NULL, NULL, '4035560061759043', '5814'),
		('200600017735', NULL, NULL, NULL, '4035560005120013', '5912'),
		('200600017933', NULL, NULL, NULL, '4035560018951990', '5691'),
		('200600018105', NULL, NULL, NULL, '4035560070198662', '4816'),
		('200600018758', NULL, NULL, NULL, '4035560014367449', '7523'),
		('200600018998', NULL, NULL, NULL, '4035560080665288', '7298'),
		('200600020135', NULL, NULL, NULL, '4035560049572351', '5814'),
		('200600020606', NULL, NULL, NULL, '4035560053820738', '5193'),
		('200600021919', NULL, NULL, NULL, '4035560075386189', '5814'),
		('200600022339', NULL, NULL, NULL, '4035560021833508', '7298'),
		('200600022941', NULL, NULL, NULL, '4035560044700882', '7297'),
		('200600023675', NULL, NULL, NULL, '4035560095419911', '5712'),
		('200600023766', NULL, NULL, NULL, '4035560089584167', '5814');

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

# Work out the current lowest level of active modules value from site, chain, acquirer, global.
UPDATE visa_qr_temp vqr
SET currentActiveModulesValue = (SELECT pd.datavalue
                                 FROM profile_data pd
                                 JOIN site_profiles sp ON pd.profile_id = sp.profile_id
                                 JOIN `profile` p ON sp.profile_id = p.profile_id AND p.profile_type_id = @siteProfileTypeId
                                 WHERE sp.site_id = vqr.siteId
                                 AND pd.data_element_id = @activeDataElementId)
WHERE currentActiveModulesValue IS NULL;

UPDATE visa_qr_temp vqr
SET currentActiveModulesValue = (SELECT pd.datavalue
                                 FROM profile_data pd
                                 JOIN site_profiles sp ON pd.profile_id = sp.profile_id
                                 JOIN `profile` p ON sp.profile_id = p.profile_id AND p.profile_type_id = @chainProfileTypeId
                                 WHERE sp.site_id = vqr.siteId
                                 AND pd.data_element_id = @activeDataElementId)
WHERE currentActiveModulesValue IS NULL;

UPDATE visa_qr_temp vqr
SET currentActiveModulesValue = (SELECT pd.datavalue
                                 FROM profile_data pd
                                 JOIN site_profiles sp ON pd.profile_id = sp.profile_id
                                 JOIN `profile` p ON sp.profile_id = p.profile_id AND p.profile_type_id = @acquirerProfileTypeId
                                 WHERE sp.site_id = vqr.siteId
                                 AND pd.data_element_id = @activeDataElementId)
WHERE currentActiveModulesValue IS NULL;

UPDATE visa_qr_temp vqr
SET currentActiveModulesValue = (SELECT pd.datavalue
                                 FROM profile_data pd
                                 JOIN site_profiles sp ON pd.profile_id = sp.profile_id
                                 JOIN `profile` p ON sp.profile_id = p.profile_id AND p.profile_type_id = @globalProfileTypeId
                                 WHERE sp.site_id = vqr.siteId
                                 AND pd.data_element_id = @activeDataElementId)
WHERE currentActiveModulesValue IS NULL;

# Ensure there is a data group entry for "modules" at site level so that we can select visaQr from the "active" element for each site.
INSERT INTO profile_data_group (profile_data_group_id, profile_Id, data_group_id, version, updated_at, updated_by, created_at, created_by)
SELECT 0, vqr.siteProfileId, @modulesDataGroupId, 1, CURRENT_TIMESTAMP, 'system', CURRENT_TIMESTAMP, 'system'
FROM visa_qr_temp vqr
WHERE NOT EXISTS (SELECT 1 FROM profile_data_group WHERE profile_id = vqr.siteProfileId AND data_group_id = @modulesDataGroupId);

# Add "visaQr" to any existing active modules.
UPDATE profile_data pd
SET datavalue = REPLACE(pd.datavalue, ']', ',"visaQr"]'),
    updated_at = CURRENT_TIMESTAMP,
    updated_by = 'system'
WHERE pd.profile_id IN (SELECT siteProfileId FROM visa_qr_temp)
AND data_element_id = @activeDataElementId;

# Add a new active modules element for anyone who doesn't have one.
INSERT INTO profile_data (profile_data_id, profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden, is_encrypted)
SELECT 0, vqr.siteProfileId, @activeDataElementId,  REPLACE(vqr.currentActiveModulesValue, ']', ',"visaQr"]'), 1, CURRENT_TIMESTAMP, 'system', current_timestamp, 'system', 1, 0, NULL
FROM visa_qr_temp vqr
WHERE NOT EXISTS (SELECT 1 FROM profile_data WHERE profile_id = vqr.siteProfileId AND data_element_id = @activeDataElementId);

# Now enable the VISA QR data group.
INSERT INTO profile_data_group (profile_data_group_id, profile_Id, data_group_id, version, updated_at, updated_by, created_at, created_by)
SELECT 0, vqr.siteProfileId, @visaQrDataGroupId, 1, CURRENT_TIMESTAMP, 'system', CURRENT_TIMESTAMP, 'system'
FROM visa_qr_temp vqr
WHERE NOT EXISTS (SELECT 1 FROM profile_data_group WHERE profile_id = vqr.siteProfileId AND data_group_id = @visaQrDataGroupId);

# And add mpan and categoryCode fields for the VISA QR options.
INSERT INTO profile_data (profile_data_id, profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden, is_encrypted)
SELECT 0, vqr.siteProfileId, @mpanDataElementId, vqr.mpan, 1, CURRENT_TIMESTAMP, 'system', current_timestamp, 'system', 1, 0, NULL
FROM visa_qr_temp vqr
WHERE NOT EXISTS (SELECT 1 FROM profile_data WHERE profile_id = vqr.siteProfileId AND data_element_id = @mpanDataElementId);

INSERT INTO profile_data (profile_data_id, profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden, is_encrypted)
SELECT 0, vqr.siteProfileId, @categoryCodeDataElementId, vqr.categoryCode, 1, CURRENT_TIMESTAMP, 'system', current_timestamp, 'system', 1, 0, NULL
FROM visa_qr_temp vqr
WHERE NOT EXISTS (SELECT 1 FROM profile_data WHERE profile_id = vqr.siteProfileId AND data_element_id = @categoryCodeDataElementId);

# Finally, ensure all site users have visaQrSale and visaQrRefund enabled.
UPDATE site_level_users slu
SET Modules = CONCAT(slu.Modules, ',', 'visaQrSale', ',', 'visaQrRefund')
WHERE slu.site_id IN (SELECT siteId FROM visa_qr_temp);

DROP TEMPORARY TABLE IF EXISTS visa_qr_temp;