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
VALUES  ('200600041354', NULL, NULL, NULL, '4035560076550610', '5411'),
		('200600039952', NULL, NULL, NULL, '4035560051611960', '5814'),
		('200600039689', NULL, NULL, NULL, '4035560005978640', '5814'),
		('200600039507', NULL, NULL, NULL, '4035560087192950', '5013'),
		('200600039226', NULL, NULL, NULL, '4035560073215880', '5814'),
		('200600038665', NULL, NULL, NULL, '4035560061258620', '5814'),
		('200600038533', NULL, NULL, NULL, '4035560092619490', '5814'),
		('200600038483', NULL, NULL, NULL, '4035560060594060', '8099'),
		('200600037980', NULL, NULL, NULL, '4035560087094450', '7299'),
		('200600037972', NULL, NULL, NULL, '4035560079266020', '4812'),
		('200600037667', NULL, NULL, NULL, '4035560086639640', '5814'),
		('200600037493', NULL, NULL, NULL, '4035560051427060', '5691'),
		('200600037279', NULL, NULL, NULL, '4035560055902860', '5995'),
		('200600036446', NULL, NULL, NULL, '4035560086206220', '5193'),
		('200600035703', NULL, NULL, NULL, '4035560077829670', '5411'),
		('200600035505', NULL, NULL, NULL, '4035560019783550', '5814'),
		('200600034821', NULL, NULL, NULL, '4035560051118900', '5944'),
		('200600034540', NULL, NULL, NULL, '4035560029871420', '5814'),
		('200600034268', NULL, NULL, NULL, '4035560054618270', '5814'),
		('200600034011', NULL, NULL, NULL, '4035560026775710', '5814'),
		('200600033997', NULL, NULL, NULL, '4035560040078370', '5814'),
		('200600033609', NULL, NULL, NULL, '4035560001022210', '5814'),
		('200600033013', NULL, NULL, NULL, '4035560070079650', '5814'),
		('200600031710', NULL, NULL, NULL, '4035560072909410', '5814'),
		('200600031645', NULL, NULL, NULL, '4035560042733700', '5814'),
		('200600030969', NULL, NULL, NULL, '4035560017414870', '5814'),
		('200600030613', NULL, NULL, NULL, '4035560071209570', '5814'),
		('200600030381', NULL, NULL, NULL, '4035560040155690', '5814'),
		('200600030290', NULL, NULL, NULL, '4035560059213580', '5814'),
		('200600028989', NULL, NULL, NULL, '4035560090239440', '5814'),
		('200600027858', NULL, NULL, NULL, '4035560096627980', '7299'),
		('200600027445', NULL, NULL, NULL, '4035560089642930', '7298'),
		('200600026785', NULL, NULL, NULL, '4035560022395100', '5814'),
		('200600026454', NULL, NULL, NULL, '4035560094480190', '5814'),
		('200600026173', NULL, NULL, NULL, '4035560063148390', '5814'),
		('200600025894', NULL, NULL, NULL, '4035560003424560', '5411'),
		('200600025217', NULL, NULL, NULL, '4035560085175490', '5411'),
		('200600024251', NULL, NULL, NULL, '4035560093241390', '7349'),
		('200600022958', NULL, NULL, NULL, '4035560075451640', '5193'),
		('200600022917', NULL, NULL, NULL, '4035560022705540', '5411'),
		('200600022404', NULL, NULL, NULL, '4035560044096380', '5814'),
		('200600021919', NULL, NULL, NULL, '4035560075386180', '5814'),
		('200600021521', NULL, NULL, NULL, '4035560006826650', '5814'),
		('200600020614', NULL, NULL, NULL, '4035560078730480', '5814'),
		('200600020465', NULL, NULL, NULL, '4035560053929830', '5814'),
		('200600020325', NULL, NULL, NULL, '4035560084723410', '5814'),
		('200600020234', NULL, NULL, NULL, '4035560023610830', '5814'),
		('200600019491', NULL, NULL, NULL, '4035560074902890', '5814'),
		('200600018600', NULL, NULL, NULL, '4035560081112340', '5814'),
		('200600017917', NULL, NULL, NULL, '4035560094471220', '5411'),
		('200600017719', NULL, NULL, NULL, '4035560061759040', '5814'),
		('200600017347', NULL, NULL, NULL, '4035560043877600', '5977'),
		('200600017230', NULL, NULL, NULL, '4035560007503480', '5814'),
		('200600017164', NULL, NULL, NULL, '4035560097489310', '5814'),
		('200600016554', NULL, NULL, NULL, '4035560079135060', '5814'),
		('200600015911', NULL, NULL, NULL, '4035560057803940', '8050'),
		('200600015358', NULL, NULL, NULL, '4035560040871620', '5814'),
		('200600013940', NULL, NULL, NULL, '4035560086025100', '5814'),
		('200600013304', NULL, NULL, NULL, '4035560041723360', '5814'),
		('200600013262', NULL, NULL, NULL, '4035560093632960', '7298'),
		('200600013098', NULL, NULL, NULL, '4035560082691870', '5814'),
		('200600012967', NULL, NULL, NULL, '4035560098437140', '5411'),
		('200600012868', NULL, NULL, NULL, '4035560036595790', '5411'),
		('200600012645', NULL, NULL, NULL, '4035560004999970', '5814'),
		('200600009922', NULL, NULL, NULL, '4035560039146120', '5814'),
		('200600008171', NULL, NULL, NULL, '4035560043328430', '5814'),
		('200600007702', NULL, NULL, NULL, '4035560045165150', '5411'),
		('200600007603', NULL, NULL, NULL, '4035560002155860', '5814'),
		('200600007520', NULL, NULL, NULL, '4035560016766960', '5251');

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
AND data_element_id = @activeDataElementId
AND datavalue NOT LIKE '%visaQr%';

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
WHERE slu.site_id IN (SELECT siteId FROM visa_qr_temp)
AND Modules NOT LIKE '%visaQrSale%';

DROP TEMPORARY TABLE IF EXISTS visa_qr_temp;