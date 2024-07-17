insert ignore into data_element (data_group_id, name, datatype, is_allow_empty, version, updated_at, updated_by, created_at, created_by, max_length, validation_expression, validation_message, front_end_validate, `unique`, options) values (1, "eposTopText", "STRING", 1, 1, NOW(), "system", NOW(), "system", null, null, null, 0, 0, "" );
insert ignore into data_element (data_group_id, name, datatype, is_allow_empty, version, updated_at, updated_by, created_at, created_by, max_length, validation_expression, validation_message, front_end_validate, `unique`, options) values (1, "eposBottomText", "STRING", 1, 1, NOW(), "system", NOW(), "system", null, null, null, 0, 0, "" );
insert ignore into profile_data(profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden) values(1, (select data_element_id from data_element where name = "eposTopText"), "Ready for Transaction", 1, NOW(), "system", NOW(), "system", 1, 0);