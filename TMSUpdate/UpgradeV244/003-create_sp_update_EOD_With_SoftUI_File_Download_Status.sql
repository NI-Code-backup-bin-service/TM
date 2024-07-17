--multiline
CREATE PROCEDURE `update_EOD_With_SoftUI_File_Download_Status`(IN tid int, IN update_coordinates varchar(50), IN update_accuracy varchar(50), IN update_coordinate_time DATETIME,IN update_free_internal_storage varchar(50), IN update_total_internal_storage varchar(50), IN last_download_datetime DATETIME, IN last_download_mainmenu_fileName varchar(30),last_download_mainmenu_filehash TEXT, last_downloaded_filelist TEXT)
BEGIN
UPDATE tid t SET t.coordinates = update_coordinates,
                 t.accuracy = update_accuracy,
                 t.last_coordinate_time = update_coordinate_time,
                 t.free_internal_storage = update_free_internal_storage,
                 t.total_internal_storage = update_total_internal_storage,
                 t.softui_last_downloaded_file_date_time = last_download_datetime,
                 t.softui_last_downloaded_file_name = last_download_mainmenu_fileName,
                 t.softui_last_downloaded_file_hash = last_download_mainmenu_filehash,
                 t.softui_last_downloaded_file_list = last_downloaded_filelist
WHERE t.tid_id = tid;
END