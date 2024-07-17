UPDATE data_element SET options = CONCAT(options, '|ENBD') WHERE name = 'PINRestrictedModules' AND options not like '%ENBD%';
