UPDATE data_element SET options = CONCAT(options, '|EnocCard') WHERE name = 'active' AND options not LIKE '%EnocCard%';