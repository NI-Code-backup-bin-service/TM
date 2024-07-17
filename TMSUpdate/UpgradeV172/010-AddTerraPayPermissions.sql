UPDATE data_element AS de SET de.options = REPLACE(de.options, '|terraPay', '|terraPaySale|terraPayRefund|terraPayVoid') WHERE de.name = "PINRestrictedModules" AND de.options LIKE '%|terraPay%' AND de.options NOT LIKE '%|terraPaySale|terraPayRefund|terraPayVoid%';
UPDATE data_element SET options = CONCAT(`options`, '|terraPaySale|terraPayRefund|terraPayVoid') WHERE name = 'supervisorOnly';