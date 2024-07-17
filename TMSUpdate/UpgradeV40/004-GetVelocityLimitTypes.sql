--multiline;
CREATE PROCEDURE `get_velocity_limit_types`()
BEGIN
	SELECT limit_type FROM txn_limit_types;
END