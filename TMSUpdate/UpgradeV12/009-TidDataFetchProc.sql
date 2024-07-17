--multiline;
CREATE PROCEDURE `tid_data_fetch`(IN TID int, IN lastChecked BIGINT)
BEGIN
	DECLARE siteId INT;
	SET @siteId = (SELECT site_id from tid_site WHERE tid_id = TID);
    
    CALL site_data_fetch(@siteId, lastChecked);

END