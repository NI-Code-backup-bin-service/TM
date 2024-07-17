--multiline
CREATE PROCEDURE SET_LAST_APK_DOWNLOAD(IN tid int, update_date DATETIME)
BEGIN
    UPDATE tid t SET t.last_apk_download = update_date
    WHERE t.tid_id = tid;
END;