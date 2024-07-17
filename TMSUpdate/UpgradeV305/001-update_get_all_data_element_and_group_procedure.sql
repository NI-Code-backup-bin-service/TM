--multiline;
CREATE PROCEDURE `get_all_data_elements_and_group_name`()
BEGIN
    SELECT de.data_element_id, de.name, de.options,dg.name FROM data_element AS de
    INNER JOIN data_group dg ON de.data_group_id = dg.data_group_id;
END;