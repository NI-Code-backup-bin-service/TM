UPDATE data_element SET sort_order_in_group = 31,tooltip = 'Online refund enabled. Please select the scheme from "Online Refund Schemes" parameter to enable it scheme wise.' WHERE name = 'onlineRefund' AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'core') LIMIT 1;
INSERT IGNORE INTO data_element(data_group_id, name, datatype, is_allow_empty, version, updated_at, updated_by, created_at, created_by, validation_expression, validation_message, front_end_validate, `unique`, options, displayname_en, is_encrypted, is_password, sort_order_in_group, tooltip, tid_overridable) VALUES ((SELECT data_group_id FROM data_group WHERE name = 'core'), 'onlineRefundSchemes', 'JSON', 1, 1, NOW(), 'system', NOW(), 'system', NULL, NULL , 0, 0, 'MASTER|VISA|DINERS', 'Online Refund Schemes', 0, 0, 32, 'Select Online Refund Schemes', 1);
INSERT IGNORE INTO data_element_locations_data_element(location_id, data_element_id) VALUES ((SELECT location_id FROM data_element_locations WHERE location_name = 'site_configuration'), (SELECT data_element_id FROM data_element WHERE name = 'onlineRefundSchemes'));