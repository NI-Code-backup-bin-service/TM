--multiline;
CREATE PROCEDURE `get_element_id_by_name`(IN element_name varchar(255))
BEGIN
	SELECT data_element_id FROM data_element WHERE `name` = element_name;	 
END