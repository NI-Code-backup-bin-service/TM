INSERT IGNORE INTO data_group (name, version, updated_at, updated_by, created_at, created_by, displayname_en) VALUES ('souhoola', '1', now(), 'system', now(), 'system', 'Souhoola');
UPDATE data_element SET options = CONCAT(options,'|souhoola') WHERE name = 'active' AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'Modules');
INSERT IGNORE INTO data_element(data_group_id, name, datatype, is_allow_empty, version, updated_at, updated_by, created_at, created_by, front_end_validate, `unique`, options, displayname_en, is_encrypted, is_password, sort_order_in_group, tooltip, tid_overridable) VALUES ((SELECT data_group_id FROM data_group WHERE name = 'souhoola'), 'username', 'STRING', 0, 1, NOW(), 'system', NOW(), 'system', 0, 1, '', 'User Name', 0, 0, 31, 'Please provide required User Name for Souhoola transaction login', 0);
INSERT IGNORE INTO data_element_locations_data_element(location_id, data_element_id) VALUES ((SELECT location_id FROM data_element_locations WHERE location_name = 'acquirer_configuration'), (SELECT data_element_id FROM data_element WHERE name = 'username' AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'souhoola')));
INSERT IGNORE INTO data_element_locations_data_element(location_id, data_element_id) VALUES ((SELECT location_id FROM data_element_locations WHERE location_name = 'chain_configuration'), (SELECT data_element_id FROM data_element WHERE name = 'username' AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'souhoola')));
INSERT IGNORE INTO data_element_locations_data_element(location_id, data_element_id) VALUES ((SELECT location_id FROM data_element_locations WHERE location_name = 'site_configuration'), (SELECT data_element_id FROM data_element WHERE name = 'username' AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'souhoola')));
INSERT IGNORE INTO data_element(data_group_id, name, datatype, is_allow_empty, version, updated_at, updated_by, created_at, created_by, front_end_validate, `unique`, options, displayname_en, is_encrypted, is_password, sort_order_in_group, tooltip, tid_overridable) VALUES ((SELECT data_group_id FROM data_group WHERE name = 'souhoola'), 'password', 'STRING', 0, 1, NOW(), 'system', NOW(), 'system', 0, 1, '', 'Password', 0, 1, 32, 'Please provide required Password for Souhoola transaction login', 0);
INSERT IGNORE INTO data_element_locations_data_element(location_id, data_element_id) VALUES ((SELECT location_id FROM data_element_locations WHERE location_name = 'acquirer_configuration'), (SELECT data_element_id FROM data_element WHERE name = 'password' AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'souhoola')));
INSERT IGNORE INTO data_element_locations_data_element(location_id, data_element_id) VALUES ((SELECT location_id FROM data_element_locations WHERE location_name = 'chain_configuration'), (SELECT data_element_id FROM data_element WHERE name = 'password' AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'souhoola')));
INSERT IGNORE INTO data_element_locations_data_element(location_id, data_element_id) VALUES ((SELECT location_id FROM data_element_locations WHERE location_name = 'site_configuration'), (SELECT data_element_id FROM data_element WHERE name = 'password' AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'souhoola')));		
