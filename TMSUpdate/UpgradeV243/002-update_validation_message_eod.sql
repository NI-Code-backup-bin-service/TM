UPDATE data_element SET validation_message = 'Must be non empty and an integer' WHERE name in ('xReadMaxPrints','zReadMaxPrints') AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'endOfDay');