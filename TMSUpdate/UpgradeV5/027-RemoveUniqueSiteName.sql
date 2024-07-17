ALTER TABLE `site` DROP INDEX `idx_site_name_version` ;
ALTER TABLE `profile` DROP INDEX `name_UNIQUE` ;
ALTER TABLE `data_element` ADD COLUMN `unique` INT NULL DEFAULT 0 AFTER `front_end_validate`;
UPDATE `data_element` SET `unique`=1 WHERE `data_element_id`=1;