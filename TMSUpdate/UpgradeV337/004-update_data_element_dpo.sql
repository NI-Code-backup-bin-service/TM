UPDATE data_element SET options = REPLACE(options, '|dpo','|dpoMomoSale|dpoMomoRefund') WHERE name = 'active' AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'modules') AND options LIKE '%dpo%' LIMIT 1;