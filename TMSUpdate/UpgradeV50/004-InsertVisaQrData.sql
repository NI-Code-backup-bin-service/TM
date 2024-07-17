UPDATE data_element SET `options` = CONCAT(`options`, '|', 'visaQr') WHERE `name` = 'active';
INSERT INTO data_group (data_group_id, `name`, version, updated_at, updated_by, created_at, created_by, displayname_en) VALUES (0, 'visaQr', 1, CURRENT_TIMESTAMP, 'system', CURRENT_TIMESTAMP, 'system', 'VISA QR');
INSERT INTO data_element (data_element_id, data_group_id, `name`, datatype, is_allow_empty, version, updated_at, updated_by, created_at, created_by, max_length, validation_expression, validation_message, front_end_validate, `unique`, `options`, displayname_en, is_encrypted, is_password, sort_order_in_group) VALUES (0, (SELECT data_group_id from data_group where name = 'visaQr'), 'mpan', 'STRING', 0, 1, CURRENT_TIMESTAMP, 'system', CURRENT_TIMESTAMP, 'system', 16, NULL, "Must be 16 characters long", 0, 0, '', 'MPAN', 0, 0, 1), (0, (SELECT data_group_id from data_group where name = 'visaQr'), 'categoryCode', 'STRING', 0, 1, CURRENT_TIMESTAMP, 'system', CURRENT_TIMESTAMP, 'system', NULL, NULL, NULL, 0, 0, '', 'Category Code', 0, 0, 2);