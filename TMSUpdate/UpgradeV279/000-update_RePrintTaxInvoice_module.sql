UPDATE data_element SET options = CONCAT(options,'|RePrintTaxInvoice') WHERE name = 'active' AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'Modules') AND options not like '%RePrintTaxInvoice%';