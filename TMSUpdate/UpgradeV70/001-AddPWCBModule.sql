UPDATE data_element SET options = CONCAT(options, '|PWCB') WHERE name = 'active' AND options not like '%PWCB%';
UPDATE data_element SET options = CONCAT(options, '|PWCB') WHERE name = 'PINRestrictedModules' AND options not like '%PWCB%';