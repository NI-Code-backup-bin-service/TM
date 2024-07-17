--multiline
CREATE PROCEDURE `get_pending_chain_change_approvals`(IN acquirers TEXT)
BEGIN
    SELECT DISTINCT
        a.approval_id as profile_data_id,
        p.name as profileName,
        de.name,
        a.change_type as change_type,
        a.current_value as original_value,
        a.new_value as updated_value,
        a.created_by as updated_by,
        a.created_at as updated_at,
        a.approved,
        a.approved_by as reviewd_by,
        a.approved_at as reviewed_at,
        a.tid_id,
        a.is_encrypted,
        a.is_password
    FROM approvals a
             LEFT JOIN data_element de ON de.data_element_id = a.data_element_id
             LEFT JOIN profile p ON p.profile_id = a.profile_id
             left join chain_profiles cp on cp.chain_profile_id = p.profile_id
             left join profile p2 on p2.profile_id = cp.acquirer_id
    where a.approved = 0
      and FIND_IN_SET(p2.name, acquirers)
    ORDER BY a.created_at DESC;
END