--multiline
create procedure log_tid_change(IN P_tid int, IN P_data_element_id int, IN P_change_type int, IN P_current_value text, IN P_updated_value text, IN P_updated_by varchar(255), IN P_approved int)
BEGIN
    SELECT tsp.site_id, tsp.profile_id INTO @site_id, @profile_id
    FROM tid_site_profiles tsp
             INNER JOIN profile p on
            tsp.profile_id = p.profile_id
             INNER JOIN profile_type pt ON
            p.profile_type_id = pt.profile_type_id
    WHERE tid_id = P_tid
    ORDER BY pt.priority asc
    LIMIT 1;

    SELECT DISTINCT p4.name
    INTO @acquires
    FROM profile p
             LEFT JOIN (site_profiles tp4
        JOIN profile p4 on p4.profile_id = tp4.profile_id
        JOIN profile_type pt4 on pt4.profile_type_id = p4.profile_type_id and pt4.priority = 4) on tp4.site_id = @site_id;

    INSERT INTO approvals(profile_id,
                          data_element_id,
                          change_type,
                          current_value,
                          new_value,
                          created_at,
                          approved_at,
                          created_by,
                          approved_by,
                          approved,
                          tid_id,
                          acquirer)
    VALUES (@profile_id,
            P_data_element_id,
            P_change_type,
            P_current_value,
            P_updated_value,
            NOW(),
            NOW(),
            P_updated_by,
            P_updated_by,
            P_approved,
            P_tid,
            @acquires);
END;

