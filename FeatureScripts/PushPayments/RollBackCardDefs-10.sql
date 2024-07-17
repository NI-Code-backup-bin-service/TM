# Temp table to store the data.
DROP TEMPORARY TABLE IF EXISTS visa_qr_temp;
CREATE TEMPORARY TABLE visa_qr_temp (merchantNo VARCHAR(12), siteId INT(11), siteProfileId INT(11), currentActiveModulesValue TEXT, mpan VARCHAR(16), categoryCode VARCHAR(255));

SET @siteProfileTypeId = (SELECT profile_type_id FROM profile_type WHERE `name` = 'site');
SET @merchantNoDataElementId = (SELECT data_element_id FROM data_element WHERE `name` = 'merchantNo');
SET @cardDefinitionsDataElementId = 21;

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

# Remove the card definitions entry from the provided site.
DELETE FROM profile_data
WHERE data_element_id = @cardDefinitionsDataElementId
AND profile_id IN (SELECT siteProfileId FROM visa_qr_temp);

# Also remove the card definitions entry from the approvals/history otherwise it will look weird as it suggests that this change should be the current setting.
DELETE FROM approvals
WHERE data_element_id = @cardDefinitionsDataElementId
AND profile_id IN (SELECT siteProfileId FROM visa_qr_temp);

DROP TEMPORARY TABLE IF EXISTS visa_qr_temp;