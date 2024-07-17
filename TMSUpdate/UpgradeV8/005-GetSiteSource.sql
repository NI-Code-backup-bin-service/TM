--multiline
CREATE PROCEDURE `get_site_element_source`(IN siteId int, IN elementId int)
BEGIN
	SELECT 
		level
    FROM site_data
    WHERE site_id = siteId
    AND data_element_id = elementId
    ORDER BY priority, version DESC
    LIMIT 1;
END