UPDATE data_element SET options = CONCAT(options, '|xls') WHERE name = 'PINRestrictedModules' AND options not like '%xls%';
