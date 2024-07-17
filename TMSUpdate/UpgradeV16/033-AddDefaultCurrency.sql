INSERT IGNORE INTO data_element (data_group_id,`name`,datatype,is_allow_empty,version,updated_at,updated_by,created_at,created_by, max_length, validation_expression, validation_message, front_end_validate, `unique`,`options`,displayname_en) VALUES ((SELECT data_group_id FROM data_group WHERE `name` = 'core' LIMIT 1),'defaultCurrency','STRING',1,1,NOW(),'system',NOW(),'system',NULL,NULL,NULL,0,0,'UAE Dirham|CFA Franc BCEAO|Egyptian Pound|Ghana Cedi|US Dollar','Default Currency');
INSERT IGNORE INTO profile_data(profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden) VALUES (1, (SELECT data_element_id FROM data_element WHERE `name` = 'defaultCurrency' LIMIT 1), "", 1, NOW(), "system", NOW(), "system", 1, 0);