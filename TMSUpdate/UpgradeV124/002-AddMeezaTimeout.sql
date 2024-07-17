INSERT IGNORE INTO data_element(data_group_id, name, datatype, is_allow_empty, version, updated_at, updated_by, created_at, created_by, validation_expression, validation_message, front_end_validate, `unique`, options, displayname_en, is_encrypted, is_password, sort_order_in_group) VALUES ((SELECT data_group_id FROM data_group WHERE name = 'meezaQR'), 'timeoutDuration', 'STRING', 0, 1, NOW(), 'system', NOW(), 'system', '^[0-9]*$', 'Timeout duration must be numeric and contain no decimal places', 1, 0, '', 'Timeout Duration (seconds)', 0, 0, 4);
INSERT IGNORE INTO data_element_locations_data_element(location_id, data_element_id) VALUES ((SELECT location_id FROM data_element_locations WHERE location_name = 'site_configuration'), (SELECT data_element_id FROM data_element WHERE name = 'timeoutDuration'));