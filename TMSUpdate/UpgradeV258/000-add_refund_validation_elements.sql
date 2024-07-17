INSERT INTO data_element (`data_group_id`, `name`, `datatype`, `is_allow_empty`, `version`, `updated_at`, `updated_by`, `created_at`, `created_by`, `max_length`, `validation_expression`, `validation_message`, `front_end_validate`, `unique`, `options`, `displayname_en`, `is_encrypted`, `is_password`, `sort_order_in_group`, `required_at_site_level`, `tooltip`, `file_max_size`, `file_min_ratio`, `file_max_ratio`, `tid_overridable`, `is_read_only_at_creation`) VALUES ((SELECT data_group_id from data_group WHERE `name` = 'core' LIMIT 1), 'refundValidationFlag', 'STRING', 1, 1, NOW(), 'System', NOW(), 'System', NULL, NULL, NULL, 0, 0, 'N|M|C', 'Refund Validation Flag', 0, 0, 27, 0, 'Refund Validation Flag stands for N- Refund Validation is not active| M- Refund validation is active at Merchant Level| C - Refund validation is active at Chain Level', NULL, NULL, NULL, 1, 0);
INSERT IGNORE INTO data_element_locations_data_element(location_id, data_element_id) VALUES ((SELECT location_id FROM data_element_locations WHERE location_name = 'site_configuration'), (SELECT data_element_id FROM data_element WHERE name = 'refundValidationFlag' AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'core')));