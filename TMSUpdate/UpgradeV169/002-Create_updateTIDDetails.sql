--multiline
CREATE PROCEDURE update_TID_Details(IN update_APK varchar(50), IN update_FW varchar(20), IN update_SW varchar(20), IN last_txn DATETIME, IN tid int, IN ip_address varchar(50), IN ip_addresses TEXT, IN sim_card_serial_number varchar(25))
BEGIN
    IF update_APK = "" THEN
        UPDATE tid t SET t.firmware_version = update_FW,
                         t.software_version = update_SW,
                         t.last_transaction_time = last_txn,
                         t.ip_address = ip_address,
                         t.ip_addresses = ip_addresses,
                         t.sim_card_serial_number = sim_card_serial_number
        WHERE t.tid_id = tid;
    ELSE
        UPDATE tid t SET t.last_apk_download = update_APK,
                         t.firmware_version = update_FW,
                         t.software_version = update_SW,
                         t.last_transaction_time = last_txn,
                         t.ip_address = ip_address,
                         t.ip_addresses = ip_addresses,
                         t.sim_card_serial_number = sim_card_serial_number
        WHERE t.tid_id = tid;
    END IF;
END