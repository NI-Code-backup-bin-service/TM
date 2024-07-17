UPDATE data_element SET tid_overridable=0 WHERE name = 'onlineRefundSchemes' AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'core');
