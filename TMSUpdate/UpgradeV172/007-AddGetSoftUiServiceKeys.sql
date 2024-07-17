--multiline;
CREATE PROCEDURE `get_softui_service_keys` (IN service_name VARCHAR(45))
BEGIN
	SELECT key_name, key_value FROM softUi AS su WHERE su.service_name = service_name;
END