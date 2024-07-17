INSERT INTO data_group (name, version, updated_at, updated_by, created_at, created_by, displayname_en) VALUES ('connectivity', '1', '2020-07-22 14:12:39', 'system', '2020-07-22 14:12:39', 'system', 'Connectivity');
INSERT IGNORE INTO data_element(data_group_id, name, datatype, is_allow_empty, version, updated_at, updated_by, created_at, created_by, front_end_validate, `unique`, options, displayname_en, is_encrypted, is_password, sort_order_in_group) VALUES ((SELECT data_group_id FROM data_group WHERE name = 'connectivity'), 'channelTimeout', 'INTEGER', 1, 1, NOW(), 'system', NOW(), 'system', 0, 0, '', 'Channel Idle Timeout Value (seconds)', 0, 0, 1);
INSERT IGNORE INTO data_element_locations_data_element(location_id, data_element_id) VALUES ((SELECT location_id FROM data_element_locations WHERE location_name = 'site_configuration'), (SELECT data_element_id FROM data_element WHERE name = 'channelTimeout'));
INSERT IGNORE INTO profile_data(profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden) VALUES (1, (SELECT data_element_id FROM data_element WHERE `name` = 'channelTimeout' LIMIT 1), 110, 1, NOW(), "system", NOW(), "system", 1, 0);