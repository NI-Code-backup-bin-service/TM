UPDATE txn_limit_types SET limit_type = 'Number of Transactions (DAILY)' WHERE limit_type_id = 1;
UPDATE txn_limit_types SET limit_type = 'Number of Transactions (BATCH)' WHERE limit_type_id = 2;
UPDATE txn_limit_types SET limit_type = 'Single Transaction Amount' WHERE limit_type_id = 3;