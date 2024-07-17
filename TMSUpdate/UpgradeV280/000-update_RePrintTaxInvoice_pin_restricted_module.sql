UPDATE data_element SET options = CONCAT(options, '|EnocCard|PrintTaxInvoice|RePrintTaxInvoice') WHERE name = 'PINRestrictedModules' AND options not like '%RePrintTaxInvoice%';
