--multiline;
CREATE PROCEDURE get_pending_change_approvals()
BEGIN
  SELECT
    a.approval_id as profile_data_id,
    p.name,
    de.name,
    a.current_value as original_value,
    a.new_value as updated_value,
    a.created_by updated_by,
    a.created_at as updated_at,
    a.approved,
    a.approved_by as reviewd_by,
    a.approved_at as reviewed_at
  FROM approvals a
  RIGHT JOIN site_profiles sp ON sp.profile_id = a.profile_id
  LEFT JOIN data_element de ON de.data_element_id = a.data_element_id
  LEFT JOIN profile p ON p.profile_id = a.profile_id
  LEFT JOIN profile_type pt ON pt.profile_type_id = p.profile_type_id
  WHERE pt.priority = 1 OR pt.priority = 2 AND a.approved = 0
  ORDER BY a.created_at DESC;
END