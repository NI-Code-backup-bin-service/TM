INSERT IGNORE INTO profile_data(profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden) VALUES (1, (SELECT data_element_id FROM data_element WHERE `name` = 'thirdPartyPackageName' LIMIT 1), "", 1, NOW(), "system", NOW(), "system", 1, 0);