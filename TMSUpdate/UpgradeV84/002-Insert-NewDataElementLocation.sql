INSERT IGNORE INTO data_element_locations (profile_type_id, location_name, location_display_name)VALUES ((SELECT profile_type_id FROM profile_type WHERE name = 'tid'),'tid_override', 'Override');
