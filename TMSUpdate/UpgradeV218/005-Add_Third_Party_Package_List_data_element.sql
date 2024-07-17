INSERT IGNORE INTO data_element(data_group_id, name, datatype, is_allow_empty, version, updated_at, updated_by, created_at, created_by, validation_expression, validation_message, front_end_validate, `unique`, options, displayname_en, is_encrypted, is_password, sort_order_in_group, tooltip,tid_overridable) VALUES ((SELECT data_group_id FROM data_group WHERE name = 'thirdParty'), 'thirdPartyPackageList', 'JSON', 1, 1, NOW(), 'system', NOW(), 'system', NULL, NULL , 0, 0, '', 'ThirdParty PackageList', 0, 0, 3, '', 1);
INSERT IGNORE INTO data_element_locations_data_element(location_id, data_element_id) VALUES ((SELECT location_id FROM data_element_locations WHERE location_name = 'tid_override'), (SELECT data_element_id FROM data_element WHERE name = 'thirdPartyPackageList' AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'thirdParty')));
INSERT IGNORE INTO profile_data(profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden,is_encrypted,not_overridable) VALUES (1, (SELECT data_element_id FROM data_element WHERE `name` = 'thirdPartyPackageList' LIMIT 1), '', 1, NOW(), 'system', NOW(), 'system', 1, 0, 0, 1);