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
VALUES  ('200600088025', NULL, NULL, NULL, '4035560080111470', '5200'),
		('200600086482', NULL, NULL, NULL, '4035560049664420', '7298'),
		('200600086219', NULL, NULL, NULL, '4035560005848670', '5947'),
		('200600084479', NULL, NULL, NULL, '4035560058853700', '5814'),
		('200600084222', NULL, NULL, NULL, '4035560059993620', '5945'),
		('200600083414', NULL, NULL, NULL, '4035560092384950', '5411'),
		('200600082200', NULL, NULL, NULL, '4035560042795060', '5411'),
		('200600080675', NULL, NULL, NULL, '4035560014636240', '5411'),
		('200600080584', NULL, NULL, NULL, '4035560038155400', '5814'),
		('200600080279', NULL, NULL, NULL, '4035560043825230', '5411'),
		('200600080139', NULL, NULL, NULL, '4035560078640210', '5814'),
		('200600080006', NULL, NULL, NULL, '4035560015627550', '5411'),
		('200600079412', NULL, NULL, NULL, '4035560032246810', '5995'),
		('200600078984', NULL, NULL, NULL, '4035560054966050', '5814'),
		('200600077564', NULL, NULL, NULL, '4035560019679880', '5912'),
		('200600077408', NULL, NULL, NULL, '4035560014115700', '7379'),
		('200600076962', NULL, NULL, NULL, '4035560030081000', '5814'),
		('200600075303', NULL, NULL, NULL, '4035560043222030', '5697'),
		('200600074264', NULL, NULL, NULL, '4035560056398300', '5814'),
		('200600073480', NULL, NULL, NULL, '4035560042046310', '5814'),
		('200600073399', NULL, NULL, NULL, '4035560089730360', '5814'),
		('200600072953', NULL, NULL, NULL, '4035560083717440', '5411'),
		('200600072938', NULL, NULL, NULL, '4035560067402790', '7339'),
		('200600072524', NULL, NULL, NULL, '4035560093052720', '5511'),
		('200600072441', NULL, NULL, NULL, '4035560025712940', '5814'),
		('200600072417', NULL, NULL, NULL, '4035560021509030', '5814'),
		('200600072094', NULL, NULL, NULL, '4035560099131640', '7211'),
		('200600071385', NULL, NULL, NULL, '4035560091634560', '8911'),
		('200600071351', NULL, NULL, NULL, '4035560086122490', '5814'),
		('200600071138', NULL, NULL, NULL, '4035560005436350', '7221'),
		('200600071013', NULL, NULL, NULL, '4035560005178590', '5411'),
		('200600070569', NULL, NULL, NULL, '4035560091286790', '8021'),
		('200600070536', NULL, NULL, NULL, '4035560093632820', '5814'),
		('200600070411', NULL, NULL, NULL, '4035560018306720', '5814'),
		('200600070403', NULL, NULL, NULL, '4035560086075030', '4722'),
		('200600070130', NULL, NULL, NULL, '4035560037449150', '5131'),
		('200600070098', NULL, NULL, NULL, '4035560054478440', '5814'),
		('200600070015', NULL, NULL, NULL, '4035560052072950', '5814'),
		('200600069983', NULL, NULL, NULL, '4035560006426570', '5814'),
		('200600069660', NULL, NULL, NULL, '4035560076814100', '5978'),
		('200600069256', NULL, NULL, NULL, '4035560071529800', '5411'),
		('200600069165', NULL, NULL, NULL, '4035560096160330', '5331'),
		('200600069058', NULL, NULL, NULL, '4035560057735180', '5411'),
		('200600068852', NULL, NULL, NULL, '4035560015282000', '5814'),
		('200600068845', NULL, NULL, NULL, '4035560044979510', '7210'),
		('200600068829', NULL, NULL, NULL, '4035560054658780', '5511'),
		('200600068761', NULL, NULL, NULL, '4035560011647450', '5411'),
		('200600068621', NULL, NULL, NULL, '4035560074652670', '5977'),
		('200600068597', NULL, NULL, NULL, '4035560049019090', '5814'),
		('200600068472', NULL, NULL, NULL, '4035560013890250', '5411'),
		('200600068449', NULL, NULL, NULL, '4035560065129540', '7538'),
		('200600068423', NULL, NULL, NULL, '4035560087116120', '7210'),
		('200600068407', NULL, NULL, NULL, '4035560024770300', '4816'),
		('200600068324', NULL, NULL, NULL, '4035560024506160', '5912'),
		('200600068126', NULL, NULL, NULL, '4035560047618520', '5697'),
		('200600068084', NULL, NULL, NULL, '4035560084919620', '5411'),
		('200600067920', NULL, NULL, NULL, '4035560043402310', '4111'),
		('200600067896', NULL, NULL, NULL, '4035560000516650', '7299'),
		('200600067839', NULL, NULL, NULL, '4035560082901540', '5814'),
		('200600067805', NULL, NULL, NULL, '4035560094175650', '5814'),
		('200600067714', NULL, NULL, NULL, '4035560069076410', '7379'),
		('200600067698', NULL, NULL, NULL, '4035560009312110', '7298'),
		('200600067672', NULL, NULL, NULL, '4035560049083720', '5814'),
		('200600067607', NULL, NULL, NULL, '4035560092863910', '5111'),
		('200600067524', NULL, NULL, NULL, '4035560055719020', '5814'),
		('200600067516', NULL, NULL, NULL, '4035560057972200', '7298'),
		('200600067425', NULL, NULL, NULL, '4035560078000120', '4111'),
		('200600067375', NULL, NULL, NULL, '4035560013151470', '5814'),
		('200600067268', NULL, NULL, NULL, '4035560011122560', '5814'),
		('200600067193', NULL, NULL, NULL, '4035560074162580', '5511'),
		('200600067128', NULL, NULL, NULL, '4035560004308730', '7542'),
		('200600067110', NULL, NULL, NULL, '4035560042042980', '4111'),
		('200600067102', NULL, NULL, NULL, '4035560062560260', '7394'),
		('200600067078', NULL, NULL, NULL, '4035560032017570', '4722'),
		('200600067011', NULL, NULL, NULL, '4035560065770910', '7392'),
		('200600066963', NULL, NULL, NULL, '4035560061860310', '5814'),
		('200600066955', NULL, NULL, NULL, '4035560086677320', '5814'),
		('200600066815', NULL, NULL, NULL, '4035560004240320', '5814'),
		('200600066690', NULL, NULL, NULL, '4035560024428380', '5511'),
		('200600066658', NULL, NULL, NULL, '4035560033106240', '7542'),
		('200600066567', NULL, NULL, NULL, '4035560092226140', '5697'),
		('200600066559', NULL, NULL, NULL, '4035560033867200', '5399'),
		('200600066484', NULL, NULL, NULL, '4035560045229170', '5065'),
		('200600066419', NULL, NULL, NULL, '4035560086550940', '5814'),
		('200600066401', NULL, NULL, NULL, '4035560060968160', '5814'),
		('200600066336', NULL, NULL, NULL, '4035560081814790', '5111'),
		('200600066294', NULL, NULL, NULL, '4035560018724300', '5814'),
		('200600066260', NULL, NULL, NULL, '4035560095153150', '5533'),
		('200600066252', NULL, NULL, NULL, '4035560080409610', '4812'),
		('200600065965', NULL, NULL, NULL, '4035560073532160', '5193'),
		('200600065916', NULL, NULL, NULL, '4035560064550520', '5814'),
		('200600065858', NULL, NULL, NULL, '4035560068414540', '5131'),
		('200600065841', NULL, NULL, NULL, '4035560031138040', '5511'),
		('200600065809', NULL, NULL, NULL, '4035560074493080', '7339'),
		('200600065775', NULL, NULL, NULL, '4035560006305730', '8111'),
		('200600065759', NULL, NULL, NULL, '4035560039196170', '5814'),
		('200600065718', NULL, NULL, NULL, '4035560080700960', '7298'),
		('200600065601', NULL, NULL, NULL, '4035560007291780', '5411'),
		('200600065569', NULL, NULL, NULL, '4035560069111470', '5046'),
		('200600065452', NULL, NULL, NULL, '4035560024858240', '5411'),
		('200600065247', NULL, NULL, NULL, '4035560051837200', '5814'),
		('200600065213', NULL, NULL, NULL, '4035560082362740', '5814'),
		('200600065163', NULL, NULL, NULL, '4035560045741890', '5814'),
		('200600065064', NULL, NULL, NULL, '4035560041739600', '5411'),
		('200600065007', NULL, NULL, NULL, '4035560087515320', '7298'),
		('200600064950', NULL, NULL, NULL, '4035560065127600', '5814'),
		('200600064935', NULL, NULL, NULL, '4035560059028150', '5814'),
		('200600064927', NULL, NULL, NULL, '4035560033383630', '5814'),
		('200600064836', NULL, NULL, NULL, '4035560057586600', '5814'),
		('200600064760', NULL, NULL, NULL, '4035560068858530', '4900'),
		('200600064752', NULL, NULL, NULL, '4035560067554060', '4722'),
		('200600064703', NULL, NULL, NULL, '4035560067115710', '5411'),
		('200600064638', NULL, NULL, NULL, '4035560091988480', '5912'),
		('200600064570', NULL, NULL, NULL, '4035560003291270', '8099'),
		('200600064562', NULL, NULL, NULL, '4035560052029930', '5814'),
		('200600064547', NULL, NULL, NULL, '4035560090574480', '5271'),
		('200600064422', NULL, NULL, NULL, '4035560003053040', '5814'),
		('200600064372', NULL, NULL, NULL, '4035560071597850', '5814'),
		('200600064364', NULL, NULL, NULL, '4035560080389600', '5411'),
		('200600064331', NULL, NULL, NULL, '4035560044364790', '5511'),
		('200600064299', NULL, NULL, NULL, '4035560012247760', '5814'),
		('200600064281', NULL, NULL, NULL, '4035560093859110', '5814'),
		('200600064273', NULL, NULL, NULL, '4035560097585740', '7216'),
		('200600064158', NULL, NULL, NULL, '4035560047658420', '5814'),
		('200600064042', NULL, NULL, NULL, '4035560046968090', '5411'),
		('200600064034', NULL, NULL, NULL, '4035560068987380', '5814'),
		('200600064000', NULL, NULL, NULL, '4035560093364690', '5411'),
		('200600063846', NULL, NULL, NULL, '4035560079137080', '5411'),
		('200600063762', NULL, NULL, NULL, '4035560038568400', '8049'),
		('200600063689', NULL, NULL, NULL, '4035560021396370', '5814'),
		('200600063606', NULL, NULL, NULL, '4035560010711640', '5814'),
		('200600063291', NULL, NULL, NULL, '4035560048032870', '7699'),
		('200600063168', NULL, NULL, NULL, '4035560011360340', '4900'),
		('200600063036', NULL, NULL, NULL, '4035560082382770', '5411'),
		('200600062848', NULL, NULL, NULL, '4035560017357890', '5814'),
		('200600062764', NULL, NULL, NULL, '4035560040498940', '7298'),
		('200600062608', NULL, NULL, NULL, '4035560020588370', '5814'),
		('200600062541', NULL, NULL, NULL, '4035560008100250', '5814'),
		('200600062467', NULL, NULL, NULL, '4035560047343600', '7392'),
		('200600062343', NULL, NULL, NULL, '4035560016687190', '5814'),
		('200600062285', NULL, NULL, NULL, '4035560056890450', '5814'),
		('200600062046', NULL, NULL, NULL, '4035560070680730', '5814'),
		('200600062004', NULL, NULL, NULL, '4035560000402300', '5814'),
		('200600061972', NULL, NULL, NULL, '4035560011287330', '5814'),
		('200600061915', NULL, NULL, NULL, '4035560084310300', '5912'),
		('200600061758', NULL, NULL, NULL, '4035560040219210', '5814'),
		('200600061584', NULL, NULL, NULL, '4035560067979800', '5039'),
		('200600061345', NULL, NULL, NULL, '4035560073449710', '5814'),
		('200600061253', NULL, NULL, NULL, '4035560085246750', '7210'),
		('200600061204', NULL, NULL, NULL, '4035560028301740', '5814'),
		('200600061139', NULL, NULL, NULL, '4035560023951470', '5814'),
		('200600061121', NULL, NULL, NULL, '4035560028464660', '5999'),
		('200600061014', NULL, NULL, NULL, '4035560078249950', '5814'),
		('200600060891', NULL, NULL, NULL, '4035560088714200', '5511'),
		('200600060529', NULL, NULL, NULL, '4035560082311160', '5814'),
		('200600060495', NULL, NULL, NULL, '4035560025471920', '5814'),
		('200600060297', NULL, NULL, NULL, '4035560010475480', '5814'),
		('200600060271', NULL, NULL, NULL, '4035560093143650', '5814'),
		('200600060230', NULL, NULL, NULL, '4035560020302460', '5814'),
		('200600060172', NULL, NULL, NULL, '4035560059642630', '5814'),
		('200600060115', NULL, NULL, NULL, '4035560052813670', '4111'),
		('200600060008', NULL, NULL, NULL, '4035560056874050', '4812'),
		('200600059729', NULL, NULL, NULL, '4035560006978680', '5912'),
		('200600059661', NULL, NULL, NULL, '4035560022117270', '5814'),
		('200600059372', NULL, NULL, NULL, '4035560047923710', '5814'),
		('200600059299', NULL, NULL, NULL, '4035560051639900', '5945'),
		('200600059265', NULL, NULL, NULL, '4035560049915320', '6513'),
		('200600059018', NULL, NULL, NULL, '4035560090072160', '4111'),
		('200600058911', NULL, NULL, NULL, '4035560030074240', '5814'),
		('200600058903', NULL, NULL, NULL, '4035560094328770', '7298'),
		('200600058887', NULL, NULL, NULL, '4035560094036940', '5814'),
		('200600058812', NULL, NULL, NULL, '4035560024001640', '8299'),
		('200600058465', NULL, NULL, NULL, '4035560042293480', '4111'),
		('200600058457', NULL, NULL, NULL, '4035560014836860', '6513'),
		('200600058226', NULL, NULL, NULL, '4035560045319870', '5411'),
		('200600058135', NULL, NULL, NULL, '4035560074707330', '5814'),
		('200600058077', NULL, NULL, NULL, '4035560084380760', '5814'),
		('200600057970', NULL, NULL, NULL, '4035560076421120', '5411'),
		('200600057855', NULL, NULL, NULL, '4035560072120610', '5411'),
		('200600057848', NULL, NULL, NULL, '4035560043559670', '5814'),
		('200600057830', NULL, NULL, NULL, '4035560066942130', '7339'),
		('200600057772', NULL, NULL, NULL, '4035560037575920', '5814'),
		('200600057681', NULL, NULL, NULL, '4035560072142140', '5699'),
		('200600057673', NULL, NULL, NULL, '4035560016496270', '5814'),
		('200600057640', NULL, NULL, NULL, '4035560053821180', '5198'),
		('200600057566', NULL, NULL, NULL, '4035560007930320', '5814'),
		('200600057533', NULL, NULL, NULL, '4035560036699880', '5814'),
		('200600057186', NULL, NULL, NULL, '4035560065465200', '5945'),
		('200600056949', NULL, NULL, NULL, '4035560076855900', '5814'),
		('200600056667', NULL, NULL, NULL, '4035560004062500', '5814'),
		('200600056527', NULL, NULL, NULL, '4035560066119840', '5814'),
		('200600056394', NULL, NULL, NULL, '4035560098424770', '5814'),
		('200600056279', NULL, NULL, NULL, '4035560069874530', '5814'),
		('200600056220', NULL, NULL, NULL, '4035560038342820', '5814'),
		('200600055974', NULL, NULL, NULL, '4035560053491170', '7210'),
		('200600055271', NULL, NULL, NULL, '4035560078030690', '5411'),
		('200600055222', NULL, NULL, NULL, '4035560033241360', '7298'),
		('200600055156', NULL, NULL, NULL, '4035560002701210', '5814'),
		('200600054746', NULL, NULL, NULL, '4035560091999810', '5912'),
		('200600054514', NULL, NULL, NULL, '4035560070027110', '5411'),
		('200600054498', NULL, NULL, NULL, '4035560099149650', '5977'),
		('200600054423', NULL, NULL, NULL, '4035560075208170', '5814'),
		('200600054290', NULL, NULL, NULL, '4035560040557840', '4722'),
		('200600053862', NULL, NULL, NULL, '4035560012754050', '5912'),
		('200600053847', NULL, NULL, NULL, '4035560049409850', '5814'),
		('200600053623', NULL, NULL, NULL, '4035560064576060', '5411'),
		('200600052773', NULL, NULL, NULL, '4035560044434940', '7296'),
		('200600051825', NULL, NULL, NULL, '4035560023162310', '7298'),
		('200600051742', NULL, NULL, NULL, '4035560052426000', '5814'),
		('200600051064', NULL, NULL, NULL, '4035560036151720', '5814'),
		('200600042477', NULL, NULL, NULL, '4035560062298050', '7298'),
		('200600042451', NULL, NULL, NULL, '4035560063467330', '5814');

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