--multiline;
CREATE PROCEDURE `search_report_data`(IN profileID INT, IN updatesFrom BIGINT)
BEGIN
  DECLARE id INT; 
  DECLARE profileType TEXT;
  DECLARE updatesFromDate DATETIME ;
  SET @updatesFromDate = FROM_UNIXTIME(updatesFrom);
  SET SESSION group_concat_max_len = 1000000;

  SELECT
    e.site_id,
    e.data_group_id,
    dg.name AS 'data_group',
    e.data_element_id,
    e.name,
    COALESCE(GROUP_CONCAT(CASE WHEN v1.priority = l1.priority THEN COALESCE(v1.level, '') END ORDER BY v1.priority), '') AS `source`,
    COALESCE(GROUP_CONCAT(CASE WHEN v1.priority = l1.priority THEN COALESCE(v1.datavalue, '') END ORDER BY v1.priority), '') AS `datavalue`,
    COALESCE(GROUP_CONCAT(CASE WHEN v1.priority = l1.priority THEN COALESCE(v1.overriden, '') END ORDER BY v1.priority), '') AS `overriden`,
    e.datatype,
    e.is_allow_empty,
    e.max_length,
    e.validation_expression,
    e.validation_message,
    e.front_end_validate,
    e.options,
    e.display_name
  FROM
    site_data_elements e
    JOIN data_group dg ON dg.data_group_id = e.data_group_id
    LEFT JOIN site_data v1 ON v1.site_id = e.site_id AND v1.data_element_id = e.data_element_id
    LEFT JOIN profile_type l1 ON l1.name = profileType
  WHERE
    e.site_id = id
    AND COALESCE(v1.updated_at, '') >= @updatesFromDate
  GROUP BY
    e.site_id,
    e.data_group_id,
    e.data_element_id,
    e.name;
END;