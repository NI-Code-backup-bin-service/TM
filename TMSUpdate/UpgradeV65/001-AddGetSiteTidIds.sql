--multiline;
CREATE PROCEDURE `get_site_tid_ids`(IN siteID INT)
BEGIN
    SELECT t.tid_id
    FROM tid t
             LEFT JOIN tid_site ts ON ts.tid_id = t.tid_id
    WHERE ts.site_id = siteID;
END