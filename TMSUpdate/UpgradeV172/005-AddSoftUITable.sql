--multiline;
CREATE TABLE `softUi` (
  `mmid` INT NOT NULL AUTO_INCREMENT,
  `service_name` VARCHAR(45) NOT NULL,
  `key_name` VARCHAR(45) NOT NULL,
  `key_value` TEXT NOT NULL,
  PRIMARY KEY (`mmid`),
  CONSTRAINT `service_key` UNIQUE (`service_name`, `key_name`));
