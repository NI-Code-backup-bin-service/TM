ALTER TABLE tid ADD COLUMN softui_last_downloaded_file_name VARCHAR(100) DEFAULT NULL;
ALTER TABLE tid ADD COLUMN softui_last_downloaded_file_hash VARCHAR(100) DEFAULT NULL;
ALTER TABLE tid ADD COLUMN softui_last_downloaded_file_list TEXT DEFAULT NULL;
ALTER TABLE tid ADD COLUMN softui_last_downloaded_file_date_time datetime DEFAULT NULL;