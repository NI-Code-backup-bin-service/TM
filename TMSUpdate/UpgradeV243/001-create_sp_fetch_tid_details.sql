--multiline
CREATE PROCEDURE `fetch_tid_details`(IN tid int)
BEGIN
      SELECT
        t.software_version,
        t.firmware_version,
        t.last_transaction_time,
        t.last_checked_time,
        t.last_apk_download,
        t.confirmed_time,
        t.eod_auto,
        t.auto_time,
        t.coordinates,
        t.accuracy,
        t.last_coordinate_time,
        t.free_internal_storage,
        t.total_internal_storage,
        t.softui_last_downloaded_file_name,
        t.softui_last_downloaded_file_hash,
        t.softui_last_downloaded_file_list,
        t.softui_last_downloaded_file_date_time
    FROM tid t
    WHERE tid_id = tid;   
END