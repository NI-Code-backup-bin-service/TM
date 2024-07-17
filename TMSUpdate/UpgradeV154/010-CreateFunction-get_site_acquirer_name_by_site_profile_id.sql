--multiline
create function get_site_acquirer_name_by_site_profile_id(p_profile_id int) RETURNS TEXT
BEGIN
    DECLARE v_acquirer_name TEXT;
    SELECT p.name INTO v_acquirer_name
    FROM site_profiles sp
             INNER JOIN site_profiles sp2 ON sp.site_id = sp2.site_id
             INNER JOIN profile p on sp2.profile_id = p.profile_id
    WHERE sp.profile_id = p_profile_id and p.profile_type_id = 2
    LIMIT 1;
    RETURN v_acquirer_name;
END;