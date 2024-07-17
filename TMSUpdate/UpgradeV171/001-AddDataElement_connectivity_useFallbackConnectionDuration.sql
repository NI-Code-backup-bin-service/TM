INSERT IGNORE INTO data_element(data_group_id, name, datatype, is_allow_empty, version, updated_at, updated_by, created_at, created_by, validation_expression, validation_message, front_end_validate, `unique`, options, displayname_en, is_encrypted, is_password, sort_order_in_group, tooltip) VALUES ((SELECT data_group_id FROM data_group WHERE name = 'connectivity'), 'useFallbackConnectionDuration', 'INTEGER', 0, 1, NOW(), 'system', NOW(), 'system', '^[0-9]{1,9}$', 'Must be numeric and contain no decimal places', 1, 0, '', 'Use Fallback Connection Duration (minutes)', 0, 0, 3, 'Determines the duration in minutes after which the PED will return to the primary connection');
INSERT IGNORE INTO data_element_locations_data_element(location_id, data_element_id) VALUES ((SELECT location_id FROM data_element_locations WHERE location_name = 'site_configuration'), (SELECT data_element_id FROM data_element WHERE name = 'useFallbackConnectionDuration' AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'connectivity')));
INSERT IGNORE INTO profile_data (profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden) VALUE (1, (SELECT data_element_id FROM data_element WHERE name = 'useFallbackConnectionDuration' AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'connectivity')), '180', 1, NOW(), 'system', NOW(), 'system', 1, 0);