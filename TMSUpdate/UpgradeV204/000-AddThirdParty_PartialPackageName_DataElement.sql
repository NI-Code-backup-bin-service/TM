INSERT IGNORE INTO data_element(data_group_id, name, datatype, is_allow_empty, version, updated_at, updated_by, created_at, created_by, max_length, validation_expression, validation_message, front_end_validate, `unique`, options, displayname_en, is_encrypted, is_password, sort_order_in_group, required_at_site_level, tooltip, tid_overridable) VALUES ((SELECT data_group_id FROM data_group WHERE name = 'thirdParty'), 'partialPackageName', 'STRING', 0, 1, NOW(), 'system', NOW(), 'system', 35, '^[\\S](.{0,34})(\\S)$', 'Must be up to 35 characters long, not have leading or trailing spaces and not blank', 1, 0, '', 'Partial Merchant Package Name', 0, 0, 4, 1, 'Must be up to 35 characters long, not have leading or trailing spaces and not blank',1);
INSERT IGNORE INTO data_element_locations_data_element(location_id, data_element_id) VALUES ((SELECT location_id FROM data_element_locations WHERE location_name = 'site_configuration'), (SELECT data_element_id FROM data_element WHERE name = 'partialPackageName' AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'thirdParty')));