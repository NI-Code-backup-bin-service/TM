--multiline;
ALTER TABLE data_element
CHANGE COLUMN `validation_expression` `validation_expression` VARCHAR(255) NULL DEFAULT NULL ,
CHANGE COLUMN `validation_message` `validation_message` VARCHAR(255) NULL DEFAULT NULL ;
