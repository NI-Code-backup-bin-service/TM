INSERT IGNORE INTO data_element(data_group_id, name, datatype, is_allow_empty, version, updated_at, updated_by, created_at, created_by, validation_expression, validation_message, front_end_validate, `unique`, options, displayname_en, is_encrypted, is_password, sort_order_in_group, tooltip, file_max_size, file_min_ratio, file_max_ratio) VALUES ((SELECT data_group_id FROM data_group WHERE name = 'receipt'),'printClearPan', 'BOOLEAN', 0, 1, NOW(), 'system', NOW(), 'system', NULL, NULL, 0, 0, NULL, 'Print Clear PAN', 0, 0, 4, 'Enables/Disables printing of the clear PAN on receipts (affects merchant copies only)', NULL, NULL, NULL);
INSERT IGNORE INTO data_element_locations_data_element(location_id, data_element_id) VALUES ((SELECT location_id FROM data_element_locations WHERE location_name = 'site_configuration'), (SELECT data_element_id FROM data_element WHERE name = 'printClearPan' AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'receipt')));
INSERT IGNORE INTO profile_data (profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden, is_encrypted, not_overridable) values (1, (SELECT data_element_id FROM data_element WHERE name = 'printClearPan' AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'receipt')), 'false', 1, NOW(), 'NISuper', NOW(), 'NISuper', 1, 1, 0, 0) ON DUPLICATE KEY UPDATE datavalue=datavalue, updated_at = updated_at;