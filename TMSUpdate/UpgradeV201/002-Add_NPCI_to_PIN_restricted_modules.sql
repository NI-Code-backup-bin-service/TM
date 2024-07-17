UPDATE data_element SET options = CONCAT(options, '|NPCI') WHERE name = 'PINRestrictedModules' AND options not like '%NPCI%';
