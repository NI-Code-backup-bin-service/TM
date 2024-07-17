--multiline
CREATE FUNCTION `get_parent_datavalue`(p_data_element_id int, p_profile_id int) RETURNS text CHARSET latin1
BEGIN
    DECLARE v_datavalue TEXT;

    SELECT pd.datavalue INTO v_datavalue
    FROM   site_profiles AS sp1
           INNER JOIN site_profiles AS sp2
                   ON sp1.site_id = sp2.site_id
           INNER JOIN profile_data AS pd
                   ON pd.profile_id = sp2.profile_id
           INNER JOIN profile AS p
                   ON p.profile_id = sp2.profile_id
    WHERE  sp1.profile_id = p_profile_id
           AND pd.data_element_id = p_data_element_id
    ORDER  BY p.profile_type_id DESC
    LIMIT  1;

    RETURN v_datavalue;
END;