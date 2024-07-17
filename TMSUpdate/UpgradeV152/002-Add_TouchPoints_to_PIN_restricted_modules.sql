UPDATE data_element SET options = CONCAT(options, '|touchpointRedeem|touchpointVoid|touchpointCheckBalance') WHERE name = 'PINRestrictedModules' AND options not like '%touchPoint%';
