#DataGroup
INSERT IGNORE INTO data_group (name, version, updated_at, updated_by, created_at, created_by, displayname_en) VALUES ('balanceInquiry', 1, NOW(), 'system', NOW(), 'system', 'Balance Inquiry');
#Add to active modules
UPDATE data_element SET options = CONCAT(options, '|balanceInquiry') WHERE name = 'active' AND options not like '%balanceInquiry%';
#Add to PIN Restricted Modules
UPDATE data_element SET options = CONCAT(options, '|balanceInquiry') WHERE name = 'PINRestrictedModules' AND options not like '%balanceInquiry%';
#Add the data element
INSERT IGNORE INTO data_element (data_group_id, name, datatype, is_allow_empty, version, updated_at, updated_by, created_at, created_by,  front_end_validate, `unique`, options, displayname_en, is_encrypted, is_password, sort_order_in_group) VALUES ((SELECT data_group_id FROM data_group WHERE data_group.name = 'balanceInquiry'), 'binRanges', 'JSON', 1, 1, NOW(), 'system', NOW(), 'system', 0, 0, '', 'Bin Ranges', 0, 0, 1);
#Add the new data element to the correct location
INSERT IGNORE INTO data_element_locations_data_element (location_id, data_element_id) VALUES ((SELECT location_id FROM data_element_locations WHERE location_name = 'site_configuration'), (SELECT data_element_id FROM data_element WHERE name = 'binRanges' AND data_group_id = (SELECT data_group_id FROM data_group WHERE data_group.name = 'balanceInquiry')));
