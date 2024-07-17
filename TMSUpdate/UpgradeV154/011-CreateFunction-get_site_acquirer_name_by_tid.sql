--multiline
create function get_site_acquirer_name_by_tid(p_tid int) RETURNS TEXT
BEGIN
    DECLARE v_acquirer_name TEXT;
    SELECT p.name INTO v_acquirer_name
    FROM site_profiles sp
             INNER JOIN tid_site ts on sp.site_id = ts.site_id
             INNER JOIN profile p on sp.profile_id = p.profile_id
    WHERE ts.tid_id = p_tid and p.profile_type_id = 2
    LIMIT 1;
    RETURN v_acquirer_name;
END;