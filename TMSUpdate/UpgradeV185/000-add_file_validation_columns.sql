-- Create new columns
ALTER TABLE data_element ADD COLUMN file_max_size INT DEFAULT NULL;
ALTER TABLE data_element ADD COLUMN file_min_ratio FLOAT DEFAULT NULL;
ALTER TABLE data_element ADD COLUMN file_max_ratio FLOAT DEFAULT NULL;
-- init as null
UPDATE data_element SET file_max_size = NULL;
UPDATE data_element SET file_min_ratio = NULL;
UPDATE data_element SET file_max_ratio = NULL;
-- set existing files max lengths:
UPDATE data_element SET file_max_size = 100000 WHERE `name`='acquirerLogo' OR `name`='applicationHeaderImg'  OR `name`='applicationFooterImg';
-- set exiting files min/max ratios
UPDATE data_element SET file_min_ratio = 1.0  WHERE `name`='acquirerLogo' OR `name`='applicationHeaderImg'  OR `name`='applicationFooterImg';
UPDATE data_element SET file_max_ratio = 6.0 WHERE `name`='acquirerLogo' OR `name`='applicationHeaderImg'  OR `name`='applicationFooterImg';