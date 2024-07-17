--multiline;
CREATE PROCEDURE `get_velocity_transactions`()
BEGIN
	SELECT txn_type, display_name FROM txn_types;
END