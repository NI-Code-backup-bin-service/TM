INSERT IGNORE INTO `data_element` (`data_group_id`, `name`, `datatype`, `is_allow_empty`, `version`, `updated_by`, `created_by`, `max_length`, `validation_expression`, `validation_message`, `front_end_validate`, `unique`, `options`, `displayname_en`, `is_encrypted`, `is_password`) VALUES ((SELECT `data_group_id` FROM data_group WHERE `name` = 'endOfDay'), 'xReportNameOverride', 'STRING', '1', '1', 'system', 'system', '12', '^.{0,12}$', 'Must be no longer than 12 characters', '1', '0', '', 'X Report Name Override', '0', '0');
INSERT IGNORE INTO `data_element` (`data_group_id`, `name`, `datatype`, `is_allow_empty`, `version`, `updated_by`, `created_by`, `max_length`, `validation_expression`, `validation_message`, `front_end_validate`, `unique`, `options`, `displayname_en`, `is_encrypted`, `is_password`) VALUES ((SELECT `data_group_id` FROM data_group WHERE `name` = 'endOfDay'), 'zReportNameOverride', 'STRING', '1', '1', 'system', 'system', '12', '^.{0,12}$', 'Must be no longer than 12 characters', '1', '0', '', 'Z Report Name Override', '0', '0');