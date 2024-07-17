--multiline;
CREATE PROCEDURE `get_available_modules`()
BEGIN
	SELECT `options` FROM data_element WHERE name = 'active';
END