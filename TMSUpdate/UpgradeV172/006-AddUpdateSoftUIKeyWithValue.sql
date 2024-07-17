--multiline;
CREATE PROCEDURE `update_softui_key_with_value` (IN service_name VARCHAR(45), IN key_name VARCHAR(45), IN key_value TEXT)
BEGIN
	REPLACE INTO softUi (`service_name`, `key_name`, `key_value`) VALUES (service_name, key_name, key_value);
END