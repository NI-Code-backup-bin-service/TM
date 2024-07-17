ALTER TABLE bulk_approvals ADD change_type int DEFAULT NULL;
update bulk_approvals set change_type=2;