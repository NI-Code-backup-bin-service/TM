UPDATE data_element SET options = CONCAT(options, '|queuedTransaction') WHERE name = 'PINRestrictedModules' AND options not like '%queuedTransaction%' AND data_group_id = (select data_group_id from data_group where name = 'modules');