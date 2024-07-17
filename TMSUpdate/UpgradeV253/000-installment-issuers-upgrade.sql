DELETE FROM data_element_locations_data_element del WHERE del.location_id = (SELECT location_id FROM data_element_locations WHERE location_name = 'site_configuration') AND del.data_element_id = (SELECT data_element_id FROM data_element de WHERE de.name = 'EPPTenor' AND de.data_group_id=(SELECT dg.data_group_id FROM data_group dg WHERE dg.name='instalments'));
DELETE FROM approvals ap WHERE ap.data_element_id = (SELECT data_element_id FROM data_element de WHERE de.data_group_id = (SELECT dg.data_group_id FROM data_group dg WHERE dg.name='instalments') AND de.name='EPPTenor');
DELETE FROM profile_data pd WHERE pd.data_element_id = (SELECT data_element_id FROM data_element de WHERE de.data_group_id = (SELECT dg.data_group_id FROM data_group dg WHERE dg.name='instalments') AND de.name='EPPTenor');
DELETE FROM data_element de WHERE de.data_group_id = (SELECT dg.data_group_id FROM data_group dg WHERE dg.name='instalments') AND de.name='EPPTenor';
INSERT IGNORE INTO data_element (`data_group_id`, `name`, `datatype`, `is_allow_empty`, `version`, `updated_at`, `updated_by`, `created_at`, `created_by`, `max_length`, `validation_expression`, `validation_message`, `front_end_validate`, `unique`, `options`, `displayname_en`, `is_encrypted`, `is_password`, `sort_order_in_group`, `required_at_site_level`, `tooltip`, `file_max_size`, `file_min_ratio`, `file_max_ratio`, `tid_overridable`, `is_read_only_at_creation`) VALUES ((SELECT data_group_id from data_group WHERE `name` = 'instalments' LIMIT 1), 'minAmount', 'STRING', 0, 1, NOW(), 'System', NOW(), 'System', NULL, '^[0-9]+(\.[0-9]{0,2})?$', 'Please ensure that the minimum amount is numeric and includes exactly two decimal places.', 0, 0, '', 'Minimum Amount', 0, 0, 5, 0, 'Minimum eligibility amount for installment plans.', NULL, NULL, NULL, 0, 0);
INSERT IGNORE INTO data_element_locations_data_element(location_id, data_element_id) VALUES ((SELECT location_id FROM data_element_locations WHERE location_name = 'site_configuration'), (SELECT data_element_id FROM data_element WHERE name = 'minAmount' AND data_group_id=(SELECT dg.data_group_id FROM data_group dg WHERE dg.name='instalments')));