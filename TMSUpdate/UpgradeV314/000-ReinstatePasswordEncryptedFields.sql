update data_element set is_password='1', is_encrypted = '1' where name='mpan' and data_group_id=(select data_group_id from data_group where name ='visaQr');
update data_element set is_password='1', is_encrypted = '1' where name='mpan' and data_group_id=(select data_group_id from data_group where name ='mastercardQr');
