--multiline
INSERT IGNORE INTO profile_data(profile_id,data_element_id,datavalue,
 version,
 updated_at,
 updated_by,
 created_at,
 created_by,
 approved,
 overriden,
 is_encrypted)
SELECT profile_id,
       (SELECT data_element_id FROM data_element WHERE name = 'onlineRefundSchemes'),
       '["MASTER","VISA"]',
       1,
       CURRENT_TIMESTAMP,
       'system',
       CURRENT_TIMESTAMP,
       'system',
       1,
       0,
       0
FROM (SELECT DISTINCT(profile_id) FROM profile_data WHERE data_element_id= (SELECT data_element_id FROM data_element WHERE name = 'onlineRefund') AND datavalue = 'true' AND
profile_id IN (SELECT DISTINCT(profile_id) FROM profile WHERE profile_type_id IN (SELECT DISTINCT(profile_type_id) FROM profile_type WHERE name IN ('acquirer','chain','site'))))
AS subquery;