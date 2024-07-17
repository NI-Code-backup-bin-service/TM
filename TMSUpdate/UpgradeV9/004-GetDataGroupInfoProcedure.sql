--multiline
CREATE PROCEDURE get_data_group_info(IN data_group_id INT)
BEGIN
  SELECT
    dg.data_group_id,
    dg.name,
    de.name
  FROM data_group dg
  JOIN data_element de ON de.data_group_id = dg.data_group_id
  WHERE dg.data_group_id = data_group_id
  ORDER BY de.data_element_id;
END;