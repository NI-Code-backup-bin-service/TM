ALTER TABLE tid_site ADD COLUMN updated_at DATETIME NULL AFTER tid_profile_id;
UPDATE tid_site SET updated_at = '2019-01-01 12:00:00';