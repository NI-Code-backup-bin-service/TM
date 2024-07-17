--multiline
DELETE FROM profile_data WHERE data_element_id = (SELECT data_element_id FROM data_element WHERE name = 'onlineRefundSchemes')
AND profile_id IN (SELECT DISTINCT(profile_id) FROM profile WHERE profile_type_id IN (SELECT DISTINCT(profile_type_id) FROM profile_type WHERE name='tid'));