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
VALUES ('001130032014', NULL, NULL, NULL, '4035560086385584', '5411'),
       ('001347640047', NULL, NULL, NULL, '4035560033658265', '5814'),
       ('001354210015', NULL, NULL, NULL, '4035560023345378', '5499'),
       ('001371260019', NULL, NULL, NULL, '4035560034773899', '5533'),
       ('200600009377', NULL, NULL, NULL, '4035560082458385', '5411'),
       ('200600012967', NULL, NULL, NULL, '4035560098437142', '5411'),
       ('200600022917', NULL, NULL, NULL, '4035560022705549', '5411'),
       ('200600026173', NULL, NULL, NULL, '4035560063148393', '5814'),
       ('200600074264', NULL, NULL, NULL, '4035560056398302', '5814');

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