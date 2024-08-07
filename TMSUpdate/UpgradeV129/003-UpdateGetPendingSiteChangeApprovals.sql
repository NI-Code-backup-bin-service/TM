--multiline
CREATE PROCEDURE `get_pending_site_change_approvals`(IN acquirers TEXT)
BEGIN
    SELECT
        a.approval_id as profile_data_id,
        pdn.datavalue,
        CONCAT(dg.name, "/", de.name),
        a.change_type as change_type,
        a.current_value as original_value,
        a.new_value as updated_value,
        a.created_by updated_by,
        a.created_at as updated_at,
        a.approved,
        a.approved_by as reviewd_by,
        a.approved_at as reviewed_at,
        a.tid_id,
        a.merchant_id,
        a.is_encrypted,
        a.is_password
    FROM approvals a
             RIGHT JOIN site_profiles sp ON sp.profile_id = a.profile_id
             LEFT JOIN data_element de ON de.data_element_id = a.data_element_id
             LEFT JOIN data_group dg ON dg.data_group_id = de.data_group_id
             LEFT JOIN profile p ON p.profile_id = a.profile_id
             LEFT JOIN profile_type pt ON pt.profile_type_id = p.profile_type_id
             LEFT JOIN profile_data pdn on pdn.profile_id = p.profile_id and pdn.data_element_id = 3
        AND pdn.version = (SELECT MAX(d.version) FROM profile_data d WHERE d.data_element_id = 3 AND d.profile_id = p.profile_id AND d.approved = 1)
             left join site t on t.site_id = sp.site_id
             LEFT JOIN (site_profiles tp4
        join profile p4 on p4.profile_id = tp4.profile_id
        join profile_type pt4 on pt4.profile_type_id = p4.profile_type_id and pt4.priority = 4) on tp4.site_id = t.site_id
    WHERE pt.priority = 1 OR pt.priority = 2 AND a.approved = 0
        and FIND_IN_SET(p4.name, acquirers)
    ORDER BY a.created_at DESC;
END