DELETE FROM profile_data WHERE (data_element_id = (SELECT data_element_id FROM data_element WHERE name = 'disableMagstripeCtls'));
DELETE FROM data_element_locations_data_element  WHERE (data_element_id = (SELECT data_element_id FROM data_element WHERE name = 'disableMagstripeCtls'));
DELETE FROM data_element WHERE (name = 'disableMagstripeCtls');