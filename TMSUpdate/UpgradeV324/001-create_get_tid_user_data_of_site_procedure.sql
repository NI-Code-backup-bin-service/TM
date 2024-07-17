-- --multiline;
CREATE PROCEDURE `get_tid_user_data_of_site`(IN siteId int)
BEGIN
    SELECT tuo.tid_user_id, tuo.tid_id, tuo.Username, tuo.PIN, tuo.Modules, tuo.is_encrypted
    FROM tid_user_override tuo
    JOIN tid_site ts ON ts.tid_id=tuo.tid_id
    WHERE ts.site_id = siteId ORDER BY tuo.tid_id, tuo.Username;
END;

