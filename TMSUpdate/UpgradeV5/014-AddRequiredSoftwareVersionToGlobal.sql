insert into data_element(data_element_id, data_group_id, name, datatype, is_allow_empty, version, updated_at, updated_by) values ('65', '9', 'RequiredSoftwareVersion', 'STRING', '1', '1', now(), 'system');
insert into profile_data(profile_data_id, profile_id, data_element_id, datavalue, version, updated_at, updated_by) values ('336', '1', '65', '114', '1', now(), 'system');
insert into profile_data_group(profile_data_group_id, profile_id, data_group_id, version, updated_at, updated_by) values ('17', '1', '9', '1', now(), 'system');