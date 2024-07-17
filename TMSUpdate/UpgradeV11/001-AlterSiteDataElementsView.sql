--multiline;
alter VIEW site_data_elements AS
SELECT
  sp.site_id,
  sp.profile_id,
  de.data_element_id,
  de.data_group_id,
  de.name,
  de.datatype,
  de.is_allow_empty,
  de.version,
  de.updated_at,
  de.updated_by,
  de.created_at,
  de.created_by,
  de.max_length,
  de.validation_expression,
  de.validation_message,
  de.front_end_validate,
  de.options
FROM
  data_element de
  JOIN profile_data_group pg ON pg.data_group_id = de.data_group_id
  JOIN site_profiles sp ON sp.profile_id = pg.profile_id