UPDATE `data_element` SET `is_password` = '1' WHERE (`name` = 'superPIN');
UPDATE `data_element` SET `is_password` = '1' WHERE (`name` = 'alipayPrivateKey');
UPDATE `approvals` SET `is_password` = 1 WHERE data_element_id IN (SELECT GROUP_CONCAT(`data_element_id`) FROM data_element WHERE `is_password` = 1);