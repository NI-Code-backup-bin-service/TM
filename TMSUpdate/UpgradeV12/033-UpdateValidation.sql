update data_element set validation_expression = '^[\\S](.{0,21})(\\S)$' where name = "name";
update data_element set validation_message = "Must be up to 23 characters long, not have leading or trailing spaces and not blank" where name = "name";
update data_element set validation_expression = '^[\\S](.{0,23})(\\S)$' where name = "addressLine1";
update data_element set validation_message = "Must be up to 25 characters long, not have leading or trailing spaces and not blank" where name = "addressLine1";
update data_element set validation_expression = '^[\\S](.{0,23})(\\S)$' where name = "addressLine2";
update data_element set validation_message = "Must be up to 25 characters long, not have leading or trailing spaces and not blank" where name = "addressLine2";