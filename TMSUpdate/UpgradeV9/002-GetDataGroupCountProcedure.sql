--multiline
CREATE PROCEDURE get_data_group_count()
BEGIN
  SELECT
    count(*)
  FROM data_group
  ORDER BY data_group_id ASC;
END