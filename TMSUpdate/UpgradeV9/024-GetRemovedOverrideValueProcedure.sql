--multiline;
CREATE PROCEDURE `get_removed_override_value`(IN siteId INT, IN elementId INT)
BEGIN
	SELECT
    datavalue
    from site_data sd
    WHERE sd.site_id = siteId
    AND sd.data_element_id = elementId
    AND priority > 2
    ORDER BY priority
    LIMIT 1; -- ignore site / tid values
END