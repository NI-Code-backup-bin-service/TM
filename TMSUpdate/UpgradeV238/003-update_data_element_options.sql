UPDATE data_element SET options = CONCAT(options,'|pushtopay') WHERE name = 'mode' AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'modules') AND `options` NOT LIKE '%pushtopay%' LIMIT 1;