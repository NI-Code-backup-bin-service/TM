UPDATE data_element SET options = REPLACE(options, '|NIIPP', '|IPP')  WHERE name = 'active';