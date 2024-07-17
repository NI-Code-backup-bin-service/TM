Alter TABLE data_element ADD COLUMN max_length int(11) DEFAULT NULL;
Alter TABLE data_element ADD COLUMN validation_expression int(11) DEFAULT NULL;
Alter TABLE data_element ADD COLUMN validation_message int(11) DEFAULT NULL;
Alter TABLE data_element ADD COLUMN front_end_validate int(11) DEFAULT 0;