INSERT IGNORE INTO data_element(data_group_id, name, datatype, is_allow_empty, version, updated_at, updated_by, created_at, created_by, front_end_validate, `unique`, options, displayname_en, is_encrypted, is_password, sort_order_in_group) VALUES ((SELECT data_group_id FROM data_group WHERE name = 'core'), 'wifiEnabled', 'BOOLEAN', 0, 1, NOW(), 'system', NOW(), 'system', 0, 0, '', 'WiFi Enabled', 0, 0, 23);
INSERT IGNORE INTO data_element_locations_data_element(location_id, data_element_id) VALUES ((SELECT location_id FROM data_element_locations WHERE location_name = 'site_configuration'), (SELECT data_element_id FROM data_element WHERE name = 'wifiEnabled'));
INSERT IGNORE into profile_data (profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden, is_encrypted) values (1, (SELECT data_element_id FROM data_element WHERE name = 'wifiEnabled'), 'false', 1, NOW(), 'NISuper', NOW(), 'NISuper', 1, 0, 0);
INSERT IGNORE INTO data_element(data_group_id, name, datatype, is_allow_empty, version, updated_at, updated_by, created_at, created_by, front_end_validate, `unique`, options, displayname_en, is_encrypted, is_password, sort_order_in_group) VALUES ((SELECT data_group_id FROM data_group WHERE name = 'userMgmt'), 'wifiPIN', 'STRING', 0, 1, NOW(), 'system', NOW(), 'system', 0, 0, '', 'WiFi PIN', 0, 1, 23);
INSERT IGNORE INTO data_element_locations_data_element(location_id, data_element_id) VALUES ((SELECT location_id FROM data_element_locations WHERE location_name = 'site_configuration'), (SELECT data_element_id FROM data_element WHERE name = 'wifiPIN'));