--multiline;
CREATE PROCEDURE `get_softui_service_key` (IN service_name VARCHAR(45), IN key_name VARCHAR(45))
BEGIN
	SELECT key_value FROM softUi AS su WHERE su.service_name = service_name AND su.key_name = key_name;
END