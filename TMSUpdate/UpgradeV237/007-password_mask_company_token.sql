UPDATE data_element SET is_password = 1 WHERE name = 'companyToken' AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'dpoMomo');