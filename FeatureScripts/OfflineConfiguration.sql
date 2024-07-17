INSERT IGNORE INTO data_group (`name`, version, updated_at, updated_by, created_at, created_by, displayname_en) values ('offline', 1, NOW(),'system', NOW(), 'system', 'offline');
INSERT IGNORE INTO data_element (data_group_id,`name`,datatype,is_allow_empty,version,updated_at,updated_by,created_at,created_by, max_length, validation_expression, validation_message, front_end_validate, `unique`,`options`)values ((SELECT data_group_id from data_group WHERE `name` = 'offline' LIMIT 1),'available','BOOLEAN',1,1,NOW(),'system',NOW(),'system',NULL,NULL,NULL,0,0,'');
INSERT IGNORE INTO data_element (data_group_id,`name`,datatype,is_allow_empty,version,updated_at,updated_by,created_at,created_by, max_length, validation_expression, validation_message, front_end_validate, `unique`,`options`)values ((SELECT data_group_id from data_group WHERE `name` = 'offline' LIMIT 1),'schemes','JSON',1,1,NOW(),'system',NOW(),'system',NULL,NULL,NULL,0,0,'');
INSERT IGNORE INTO data_element (data_group_id,`name`,datatype,is_allow_empty,version,updated_at,updated_by,created_at,created_by, max_length, validation_expression, validation_message, front_end_validate, `unique`,`options`, displayname_en, is_encrypted, is_password, sort_order_in_group) values ((SELECT data_group_id from data_group WHERE `name` = 'offline' LIMIT 1),'totalLimit','LONG',1,1,NOW(),'system',NOW(),'system',NULL,"^[0]|([1-9]\\d\\d|[1-9]\\d{3,})$","Denomination must be 100 (1 AED) or higher",1,0,'',"",0,0,1) on duplicate key update datatype = "LONG", name = "totalLimit", validation_expression = "^[0]|([1-9]\\d\\d|[1-9]\\d{3,})$", validation_message = "Denomination must be either 0 (no limit) or 100 (1 AED) or higher", front_end_validate = 1, `unique` = 1, `options` = NULL, displayname_en = "", is_encrypted = 0, is_password = 0, sort_order_in_group = 1;
INSERT IGNORE INTO data_element (data_group_id,`name`,datatype,is_allow_empty,version,updated_at,updated_by,created_at,created_by, max_length, validation_expression, validation_message, front_end_validate, `unique`,`options`, displayname_en, is_encrypted, is_password, sort_order_in_group) values ((SELECT data_group_id from data_group WHERE `name` = 'offline' LIMIT 1),'uploadFrequency','INTEGER',1,1,NOW(),'system',NOW(),'system',NULL,"^[0-9]*$","Value must be numeric",1,0,'',"",0,0,2) on duplicate key update datatype = "INTEGER", name = "uploadFrequency", validation_expression = "^[0-9]*$", validation_message = "Value must be numeric", front_end_validate = 1, `unique` = 0, `options` = NULL, displayname_en = "", is_encrypted = 0, is_password = 0, sort_order_in_group = 2;
INSERT IGNORE INTO data_element (data_group_id,`name`,datatype,is_allow_empty,version,updated_at,updated_by,created_at,created_by, max_length, validation_expression, validation_message, front_end_validate, `unique`,`options`, displayname_en, is_encrypted, is_password, sort_order_in_group) values ((SELECT data_group_id from data_group WHERE `name` = 'offline' LIMIT 1),'offlineDuration','LONG',1,1,NOW(),'system',NOW(),'system',NULL,"^[0-9]*$","Value must be a numeric",1,0,'',"",0,0,3) on duplicate key update datatype = "LONG", name = "offlineDuration", validation_expression = "^[0-9]*$", validation_message = "Value must be a numeric", front_end_validate = 1, `unique` = 0, `options` = NULL, displayname_en = "", is_encrypted = 0, is_password = 0, sort_order_in_group = 3;
