ALTER TABLE data_group ADD COLUMN `displayname_en` TEXT NULL AFTER `created_by`;
UPDATE data_group SET displayname_en = name WHERE displayname_en IS NULL;
ALTER TABLE data_group CHANGE COLUMN `displayname_en` `displayname_en` TEXT NOT NULL;
ALTER TABLE data_element ADD COLUMN `sort_order_in_group` INT(11) NULL AFTER `is_password`;