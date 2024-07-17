insert ignore into data_group (`name`, version, updated_at, updated_by, created_at, created_by)  values ('upi', 1, NOW(), 'system', NOW(), 'system');
insert ignore into data_element (data_group_id, name, datatype, is_allow_empty, version, updated_at, updated_by, created_at, created_by, max_length, validation_expression, validation_message, front_end_validate, `unique`, options) values ((SELECT data_group_id from data_group WHERE `name` = 'upi' LIMIT 1), 'acquirerIIN', 'STRING', 1, 1, NOW(), "system", NOW(), "system", null, null, null, 0, 0, "" );
insert ignore into data_element (data_group_id, name, datatype, is_allow_empty, version, updated_at, updated_by, created_at, created_by, max_length, validation_expression, validation_message, front_end_validate, `unique`, options) values ((SELECT data_group_id from data_group WHERE `name` = 'upi' LIMIT 1), 'forwardingIIN', 'STRING', 1, 1, NOW(), "system", NOW(), "system", null, null, null, 0, 0, "" );
insert ignore into data_element (data_group_id, name, datatype, is_allow_empty, version, updated_at, updated_by, created_at, created_by, max_length, validation_expression, validation_message, front_end_validate, `unique`, options) values ((SELECT data_group_id from data_group WHERE `name` = 'upi' LIMIT 1), 'categoryCode', 'STRING', 1, 1, NOW(), "system", NOW(), "system", null, null, null, 0, 0, "" );
insert ignore into data_element (data_group_id, name, datatype, is_allow_empty, version, updated_at, updated_by, created_at, created_by, max_length, validation_expression, validation_message, front_end_validate, `unique`, options) values ((SELECT data_group_id from data_group WHERE `name` = 'upi' LIMIT 1), 'countryCode', 'STRING', 1, 1, NOW(), "system", NOW(), "system", null, null, null, 0, 0, "" );
