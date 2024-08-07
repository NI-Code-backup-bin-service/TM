--multiline;
CREATE PROCEDURE `get_profile_change_history`(IN profileId INT)
BEGIN
  SELECT
    de.name,
    a.current_value as original_value,
    a.new_value,
    a.created_by as updated_by,
    a.created_at as updated_at,
    a.approved
  FROM approvals a
  LEFT JOIN data_element de ON de.data_element_id = a.data_element_id
  WHERE a.profile_id = profileId
  ORDER BY a.created_at DESC;
END