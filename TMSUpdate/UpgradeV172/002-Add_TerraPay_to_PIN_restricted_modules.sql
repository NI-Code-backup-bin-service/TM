UPDATE data_element SET options = CONCAT(options, '|terraPay') WHERE name = 'PINRestrictedModules' AND options not like '%terraPay%';
