UPDATE data_element SET tooltip = 'Interface that transaction retrieval uses.' WHERE name = 'interface' AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'transactionRetrieval');
UPDATE data_element SET tooltip = 'Parameters that the customer will be prompted for when going through the TR flow.' WHERE name = 'parameterDefinitions' AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'transactionRetrieval');
UPDATE data_element SET tooltip = 'Customer receipt template to be used for custom TR fields.' WHERE name = 'customerReceiptTemplate' AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'transactionRetrieval');
UPDATE data_element SET tooltip = 'Merchant receipt template to be used for custom TR fields.' WHERE name = 'merchantReceiptTemplate' AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'transactionRetrieval');
UPDATE data_element SET tooltip = 'Site ID send in the transaction retrieval request.' WHERE name = 'siteID' AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'transactionRetrieval');
/* VPS */
UPDATE data_element SET tooltip = 'Enables VPS.' WHERE name = 'enabled' AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'vps');
UPDATE data_element SET tooltip = 'Upper financial limit of where VPS is enabled.' WHERE name = 'limit' AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'vps');
/* QPS */
UPDATE data_element SET tooltip = 'Enables QPS.' WHERE name = 'enabled' AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'qps');
UPDATE data_element SET tooltip = 'Upper financial limit of where QPS is enabled.' WHERE name = 'limit' AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'qps');