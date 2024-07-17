--multiline
UPDATE profile_data
SET datavalue = 'DF1006999999999999 DF1106000000000000 DF1206000000050000 DF130100 DF190113 DF220100 DF2606999999999999'
where profile_id = (SELECT profile_id from profile where name = 'global')
and data_element_id = (SELECT data_element_id from data_element where name = 'ctlsApplicationConfigShared'
and data_group_id = (SELECT data_group_id from data_group where name = 'dualCurrency'));