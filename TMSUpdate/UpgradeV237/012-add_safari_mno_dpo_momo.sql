UPDATE data_element SET options = CONCAT(options,'|SafaricomC2B') WHERE name = 'mno' AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'dpoMomo');