ALTER TABLE profile_data ADD COLUMN `is_encrypted` BOOLEAN;
ALTER TABLE tid ADD COLUMN `is_encrypted` BOOLEAN NOT NULL DEFAULT 0;
ALTER TABLE tid_user_override ADD COLUMN `is_encrypted` BOOLEAN NOT NULL DEFAULT 0;
ALTER TABLE data_element ADD COLUMN `is_encrypted` BOOLEAN NOT NULL DEFAULT 0;
ALTER TABLE data_element ADD COLUMN `is_password` BOOLEAN NOT NULL DEFAULT 0;
ALTER TABLE approvals ADD COLUMN `is_encrypted` BOOLEAN NOT NULL DEFAULT 0;
ALTER TABLE approvals ADD COLUMN `is_password` BOOLEAN NOT NULL DEFAULT 0;
ALTER TABLE site_level_users ADD COLUMN `is_encrypted` BOOLEAN NOT NULL DEFAULT 0;
UPDATE `approvals` SET `is_password` = '1' WHERE data_element_id IN (SELECT GROUP_CONCAT(data_element_id) FROM data_element WHERE is_password = 1);