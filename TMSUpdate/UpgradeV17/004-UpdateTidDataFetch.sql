--multiline
CREATE PROCEDURE `tid_data_fetch`(IN TID int, IN lastChecked BIGINT)
BEGIN
	DECLARE siteId INT;
    DECLARE lastOverrideRemoval DATETIME;
    
	SET @siteId = (SELECT site_id from tid_site WHERE tid_id = TID);
    
    IF ((SELECT updated_at FROM tid_site WHERE tid_id = tid) > FROM_UNIXTIME(lastChecked)) OR 
		((SELECT updated_at FROM site WHERE site_id = @siteId) > FROM_UNIXTIME(lastChecked))
    THEN
        CALL site_data_fetch(@siteId, 0);
    ELSE
		CALL site_data_fetch(@siteId, lastChecked);
	END IF;
END