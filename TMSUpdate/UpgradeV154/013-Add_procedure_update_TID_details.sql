--multiline
CREATE PROCEDURE update_TID_Details(IN update_APK DATETIME, IN update_FW varchar(20), IN update_SW varchar(20), IN last_txn DATETIME, IN tid int)
BEGIN
    UPDATE tid t SET t.last_apk_download = update_APK,
                     t.firmware_version = update_FW,
                     t.software_version = update_SW,
                     t.last_transaction_time = last_txn
    WHERE t.tid_id = tid;
END