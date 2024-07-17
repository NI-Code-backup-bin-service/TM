UPDATE profile_data SET datavalue = 15 WHERE data_element_id = 61 AND IFNULL(datavalue, 0) > 15;
UPDATE profile_data SET datavalue = 1 WHERE data_element_id = 61 AND IFNULL(datavalue, 0) = 0;