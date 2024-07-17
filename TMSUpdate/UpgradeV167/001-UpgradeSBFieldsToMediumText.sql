ALTER TABLE profile_data CHANGE COLUMN `datavalue` `datavalue` MEDIUMTEXT;
ALTER TABLE approvals CHANGE COLUMN `current_value` `current_value` MEDIUMTEXT  , CHANGE COLUMN `new_value` `new_value` MEDIUMTEXT;
