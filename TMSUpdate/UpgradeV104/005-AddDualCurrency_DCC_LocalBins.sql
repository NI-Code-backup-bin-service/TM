INSERT IGNORE INTO data_element(data_group_id, name, datatype, is_allow_empty, version, updated_at, updated_by, created_at, created_by, validation_expression, validation_message, front_end_validate, `unique`, options, displayname_en, is_encrypted, is_password, sort_order_in_group) VALUES ((SELECT data_group_id FROM data_group WHERE name = 'dualCurrency'), 'dualCurrencyDccLocalBins', 'JSON', 1, 1, NOW(), 'system', NOW(), 'system', NULL, NULL, 0, 0, '', 'DCC Local Bins', 0, 0, 11);
INSERT IGNORE INTO data_element_locations_data_element(location_id, data_element_id) VALUES ((SELECT location_id FROM data_element_locations WHERE location_name = 'site_configuration'), (SELECT data_element_id FROM data_element WHERE name = 'dualCurrencyDccLocalBins'));