delete from approvals where data_element_id=(select data_element_id from data_element where name='autoCutOver' and data_group_id=(select data_group_id from data_group where name='endOfDay'));
delete from data_element_locations_data_element  where data_element_id=(select data_element_id from data_element where name='autoCutOver' and data_group_id=(select data_group_id from data_group where name='endOfDay'));
delete from profile_data  where data_element_id=(select data_element_id from data_element where name='autoCutOver' and data_group_id=(select data_group_id from data_group where name='endOfDay'));
delete from data_element where name='autoCutOver' and data_group_id=(select data_group_id from data_group where name='endOfDay');