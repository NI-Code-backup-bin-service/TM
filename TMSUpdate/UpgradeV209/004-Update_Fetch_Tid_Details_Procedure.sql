--multiline
CREATE PROCEDURE `fetch_tid_details`(IN tid INT)
BEGIN
    SELECT
        t.software_version,
        t.firmware_version,
        t.last_transaction_time,
        t.last_checked_time,
        t.last_apk_download,
        t.confirmed_time,
        t.coordinates,
        t.accuracy,
        t.last_coordinate_time
    FROM tid t
    WHERE tid_id = tid;
END