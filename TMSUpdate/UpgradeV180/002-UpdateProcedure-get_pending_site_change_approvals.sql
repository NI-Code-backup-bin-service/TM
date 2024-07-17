--multiline
CREATE PROCEDURE `get_pending_site_change_approvals`(IN acquirers TEXT)
BEGIN
    SELECT a.approval_id   AS profile_data_id,
           pdn.datavalue,
           Concat(dg.name, "/", de.name),
           a.change_type   AS change_type,
           a.current_value AS original_value,
           a.new_value     AS updated_value,
           a.created_by       updated_by,
           a.created_at    AS updated_at,
           a.approved,
           a.approved_by   AS reviewd_by,
           a.approved_at   AS reviewed_at,
           a.tid_id,
           a.merchant_id,
           a.is_encrypted,
           a.is_password
    FROM approvals a

             INNER JOIN profile p
                        ON p.profile_id = a.profile_id
             INNER JOIN profile_type pt
                        ON pt.profile_type_id = p.profile_type_id
             INNER JOIN profile_data pdn
                        ON pdn.profile_id = p.profile_id
                            AND pdn.data_element_id = 3
                            AND pdn.version = (SELECT Max(d.version)
                                               FROM profile_data d
                                               WHERE d.data_element_id = 3
                                                 AND d.profile_id = p.profile_id
                                                 AND d.approved = 1)
             INNER JOIN data_element de
                        ON de.data_element_id = a.data_element_id
             INNER JOIN data_group dg
                        ON dg.data_group_id = de.data_group_id
             INNER JOIN (site_profiles tp4
        INNER JOIN profile p4
        ON p4.profile_id = tp4.profile_id
        INNER JOIN profile_type pt4
        ON pt4.profile_type_id = p4.profile_type_id
            AND pt4.priority = 4)
                        ON tp4.site_id = (SELECT DISTINCT t.site_id
                                          FROM site t
                                          WHERE t.site_id = (SELECT DISTINCT sp.site_id
                                                             FROM site_profiles sp
                                                             WHERE sp.profile_id = a.profile_id))

    WHERE a.approved = 0
      AND (pt.priority = 1 OR pt.priority = 2)
      AND Find_in_set(p4.name, acquirers)
    ORDER BY a.created_at DESC;
END