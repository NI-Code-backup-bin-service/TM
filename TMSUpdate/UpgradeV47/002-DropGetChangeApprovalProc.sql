--multiline;
CREATE PROCEDURE `get_change_approval_history`(
    IN afterDate VARCHAR(255),
    IN name VARCHAR(255),
    IN user VARCHAR(255),
    IN beforeDate VARCHAR(255),
    IN field VARCHAR(255),
    IN acquirers TEXT
)
BEGIN
    set @name = upper(concat('%', ifnull(name,''), '%'));
    set @user = upper(concat('%', ifnull(user,''), '%'));
    set @field = upper(concat('%', ifnull(field,''), '%'));
    SELECT distinct
        a.approval_id as profile_data_id,
        p.name,
        de.name as 'field',
        a.current_value as original_value,
        a.new_value as updated_value,
        a.created_by as updated_by,
        a.created_at as updated_at,
        a.approved,
        a.approved_by as reviewed_by,
        a.approved_at as reviewed_at,
        a.tid_id,
        a.merchant_id,
        a.is_password,
        a.is_encrypted
    FROM approvals a
             LEFT JOIN data_element de ON de.data_element_id = a.data_element_id
             LEFT JOIN profile p ON p.profile_id = a.profile_id
        -- chains
             left join chain_profiles cp on cp.chain_profile_id = p.profile_id
             left join profile p2 on p2.profile_id = cp.acquirer_id
        -- acquirers
             left join chain_profiles cp2 on cp2.acquirer_id = p.profile_id
             left join profile p3 on p3.profile_id = cp2.acquirer_id
        -- sites
             left join ( profile p4
        JOIN site_profiles sp on sp.profile_id = p4.profile_id
        join site_profiles sp2 on sp2.site_id = sp.site_id
        join profile p5 on p5.profile_id = sp2.profile_id) on p5.profile_id = a.profile_id and p5.profile_type_id = 4
        -- tids
             left join tid td on td.tid_id = p.name
             left join tid_site ts on ts.tid_id = p.name
             left join site t on t.site_id = ts.site_id
             LEFT JOIN (site_profiles tp4
        join profile p6 on p6.profile_id = tp4.profile_id
        join profile_type pt4 on pt4.profile_type_id = p6.profile_type_id and pt4.priority = 4) on tp4.site_id = t.site_id and td.tid_id != 0
    WHERE a.approved != 0
      AND (afterDate IS NULL OR a.created_at >= afterDate)
      AND p.name like @name
      AND a.approved_by like @user
      AND (beforeDate IS NULL OR a.created_at <= beforeDate)
      AND de.name like @field
      and ( FIND_IN_SET(p2.name, acquirers)
        or FIND_IN_SET(p3.name, acquirers)
        or ( FIND_IN_SET(p4.name, acquirers) or FIND_IN_SET(a.acquirer, acquirers))
        or ( FIND_IN_SET(p6.name, acquirers) or FIND_IN_SET(a.acquirer, acquirers)) )
    ORDER BY a.approved_at DESC;
END
