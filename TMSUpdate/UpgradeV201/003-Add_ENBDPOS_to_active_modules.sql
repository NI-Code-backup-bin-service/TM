UPDATE data_element SET options = CONCAT(options, '|ENBD') WHERE name = 'active' AND options not LIKE '%ENBD%';