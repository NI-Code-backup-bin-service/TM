--multiline
CREATE PROCEDURE `get_pending_tid_change_approvals`(IN acquirers TEXT)
BEGIN
    SELECT
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
        a.tid_id
    FROM approvals a
             LEFT JOIN data_element de ON de.data_element_id = a.data_element_id
             LEFT JOIN profile p ON p.profile_id = a.profile_id
             left join tid_site ts on ts.tid_id = p.name
             left join site t on t.site_id = ts.site_id
             LEFT JOIN (site_profiles tp4
        join profile p4 on p4.profile_id = tp4.profile_id
        join profile_type pt4 on pt4.profile_type_id = p4.profile_type_id and pt4.priority = 4) on tp4.site_id = t.site_id
    WHERE p.profile_type_id = (select profile_type_id from profile_type where profile_type.name = "tid")
      AND a.approved = 0
      and FIND_IN_SET(p4.name, acquirers)
    ORDER BY a.created_at DESC;
END