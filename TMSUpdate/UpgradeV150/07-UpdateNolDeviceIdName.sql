UPDATE data_element SET displayname_en = 'Nol Device ID' where name = 'nolDeviceId' AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'nol');
UPDATE profile_data SET datavalue = '' WHERE profile_id = (SELECT profile_id FROM NextGen_TMS.profile WHERE name='global') AND data_element_id = (SELECT data_element_id FROM data_element WHERE name = 'nolDeviceId');