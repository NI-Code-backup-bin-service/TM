UPDATE data_element SET options = REPLACE(options, '|terraPayVoid', '') WHERE name IN ('PINRestrictedModules', 'supervisorOnly') AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'modules');