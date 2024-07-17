--multiline;
CREATE PROCEDURE `get_pending_profile_change_approvals`(IN typeId INT)
BEGIN
	SELECT
		a.approval_id as profile_data_id,
		p.name as profileName,
		de.name,
		a.current_value as original_value,
		a.new_value as updated_value,
		a.created_by as updated_by,
		a.created_at as updated_at,
		a.approved,
		a.approved_by as reviewd_by,
		a.approved_at as reviewed_at
	FROM approvals a
    LEFT JOIN data_element de ON de.data_element_id = a.data_element_id
    LEFT JOIN profile p ON p.profile_id = a.profile_id
    WHERE p.profile_type_id = typeId AND a.approved = 0
    ORDER BY a.created_at DESC;
END