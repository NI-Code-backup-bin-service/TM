UPDATE data_element SET datatype = 'STRING', validation_expression = '^[0-9a-fA-Z]{0,6}$' WHERE name = 'nolDeviceId' AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'nol')