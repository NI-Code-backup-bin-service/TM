UPDATE data_element SET is_allow_empty = 1 WHERE name = 'ctlsCvmLimit' AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'dualCurrency');
UPDATE data_element SET is_allow_empty = 1 WHERE name = 'ctlsTxnLimit' AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'dualCurrency');