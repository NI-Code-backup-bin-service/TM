INSERT INTO data_group (`name`, `version`, `updated_at`, `updated_by`, `created_at`, `created_by`, `displayname_en`) VALUES('paymentServices', 1, NOW(), 'system', NOW(), 'system', 'Payment Services');
INSERT INTO data_element (`data_group_id`, `name`, `datatype`, `is_allow_empty`, `version`, `updated_at`, `updated_by`, `created_at`, `created_by`, `max_length`, `validation_expression`, `validation_message`, `front_end_validate`, `unique`, `options`, `displayname_en`, `is_encrypted`, `is_password`, `sort_order_in_group`, `required_at_site_level`, `tooltip`, `file_max_size`, `file_min_ratio`, `file_max_ratio`, `tid_overridable`, `is_read_only_at_creation`) VALUES ((SELECT data_group_id from data_group WHERE `name` = 'paymentServices' LIMIT 1), 'paymentServiceEnabled', 'BOOLEAN', 0, 1, NOW(), 'System', NOW(), 'System', NULL, NULL, NULL, 0, 0, '', 'Enabled', 0, 0, 1, 0, 'Enable/Disable Payment Services', NULL, NULL, NULL, 0, 0);
INSERT IGNORE INTO data_element_locations_data_element(location_id, data_element_id) VALUES ((SELECT location_id FROM data_element_locations WHERE location_name = 'site_configuration'), (SELECT data_element_id FROM data_element WHERE name = 'paymentServiceEnabled' AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'paymentServices')));
INSERT INTO data_element (`data_group_id`, `name`, `datatype`, `is_allow_empty`, `version`, `updated_at`, `updated_by`, `created_at`, `created_by`, `max_length`, `validation_expression`, `validation_message`, `front_end_validate`, `unique`, `options`, `displayname_en`, `is_encrypted`, `is_password`, `sort_order_in_group`, `required_at_site_level`, `tooltip`, `file_max_size`, `file_min_ratio`, `file_max_ratio`, `tid_overridable`, `is_read_only_at_creation`) VALUES ((SELECT data_group_id from data_group WHERE `name` = 'paymentServices' LIMIT 1), 'paymentAuthType', 'STRING', 0, 1, NOW(), 'System', NOW(), 'System', NULL, NULL, NULL, 0, 0, 'Single-Auth|Multi-Auth', 'Authentication Type', 0, 0, 1, 0, 'Please select service authentication type.', NULL, NULL, NULL, 0, 0);
INSERT IGNORE INTO data_element_locations_data_element(location_id, data_element_id) VALUES ((SELECT location_id FROM data_element_locations WHERE location_name = 'site_configuration'), (SELECT data_element_id FROM data_element WHERE name = 'paymentAuthType' AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'paymentServices')));
INSERT INTO data_element (`data_group_id`, `name`, `datatype`, `is_allow_empty`, `version`, `updated_at`, `updated_by`, `created_at`, `created_by`, `max_length`, `validation_expression`, `validation_message`, `front_end_validate`, `unique`, `options`, `displayname_en`, `is_encrypted`, `is_password`, `sort_order_in_group`, `required_at_site_level`, `tooltip`, `file_max_size`, `file_min_ratio`, `file_max_ratio`, `tid_overridable`, `is_read_only_at_creation`) VALUES ((SELECT data_group_id from data_group WHERE `name` = 'paymentServices' LIMIT 1), 'paymentServiceGroup', 'STRING', 0, 1, NOW(), 'System', NOW(), 'System', NULL, NULL, NULL, 0, 0, '', 'Service Group', 0, 0, 1, 0, 'Please select service group.', NULL, NULL, NULL, 0, 0);
INSERT IGNORE INTO data_element_locations_data_element(location_id, data_element_id) VALUES ((SELECT location_id FROM data_element_locations WHERE location_name = 'site_configuration'), (SELECT data_element_id FROM data_element WHERE name = 'paymentServiceGroup' AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'paymentServices')));
INSERT INTO data_element (`data_group_id`, `name`, `datatype`, `is_allow_empty`, `version`, `updated_at`, `updated_by`, `created_at`, `created_by`, `max_length`, `validation_expression`, `validation_message`, `front_end_validate`, `unique`, `options`, `displayname_en`, `is_encrypted`, `is_password`, `sort_order_in_group`, `required_at_site_level`, `tooltip`, `file_max_size`, `file_min_ratio`, `file_max_ratio`, `tid_overridable`, `is_read_only_at_creation`) VALUES ((SELECT data_group_id from data_group WHERE `name` = 'paymentServices' LIMIT 1), 'paymentServicesConfigs', 'JSON', 1, 1, NOW(), 'System', NOW(), 'System', NULL, NULL, NULL, 0, 0, '', 'Services', 0, 0, 1, 0, 'Please configure TID and MID of the services.', NULL, NULL, NULL, 1, 0);
INSERT IGNORE INTO data_element_locations_data_element(location_id, data_element_id) VALUES ((SELECT location_id FROM data_element_locations WHERE location_name = 'tid_override'), (SELECT data_element_id FROM data_element WHERE name = 'paymentServicesConfigs' AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'paymentServices')));
INSERT INTO data_element (`data_group_id`, `name`, `datatype`, `is_allow_empty`, `version`, `updated_at`, `updated_by`, `created_at`, `created_by`, `max_length`, `validation_expression`, `validation_message`, `front_end_validate`, `unique`, `options`, `displayname_en`, `is_encrypted`, `is_password`, `sort_order_in_group`, `required_at_site_level`, `tooltip`, `file_max_size`, `file_min_ratio`, `file_max_ratio`, `tid_overridable`, `is_read_only_at_creation`) VALUES ((SELECT data_group_id from data_group WHERE `name` = 'paymentServices' LIMIT 1), 'paymentMode', 'STRING', 1, 1, NOW(), 'System', NOW(), 'System', NULL, NULL, NULL, 0, 0, 'Single|Multiple', 'Payment Mode', 0, 0, 1, 0, 'Please select payment mode.', NULL, NULL, NULL, 1, 0);
INSERT IGNORE INTO data_element_locations_data_element(location_id, data_element_id) VALUES ((SELECT location_id FROM data_element_locations WHERE location_name = 'tid_override'), (SELECT data_element_id FROM data_element WHERE name = 'paymentMode' AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'paymentServices')));