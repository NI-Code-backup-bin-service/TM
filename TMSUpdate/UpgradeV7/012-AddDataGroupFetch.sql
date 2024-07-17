--multiline
CREATE PROCEDURE `fetch_data_groups`()
BEGIN
 SELECT dg.data_group_id, dg.name FROM data_group dg;
END