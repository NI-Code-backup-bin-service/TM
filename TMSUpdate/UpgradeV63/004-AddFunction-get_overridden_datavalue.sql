--multiline
create function get_parent_datavalue(p_data_element_id int, p_profile_id int) returns text
BEGIN
    DECLARE v_datavlue TEXT;
    SELECT
        pd.datavalue INTO v_datavlue
    FROM profile_data pd
             INNER JOIN profile p ON
            p.profile_id = pd.profile_id
             INNER JOIN profile_type pt ON
            p.profile_type_id = pt.profile_type_id
    WHERE
            pd.data_element_id = p_data_element_id
      AND
            pd.datavalue != ''
      AND
        pd.datavalue IS NOT NULL
      AND
        (
            p.profile_id = p_profile_id
            OR
            (
                pt.profile_type_id != (SELECT profile_type_id FROM profile WHERE profile_id = p_profile_id)
                AND pt.priority > (
                    SELECT profile_type.priority
                    FROM profile
                    INNER JOIN profile_type ON
                        profile.profile_type_id = profile_type.profile_type_id
                    WHERE profile_id = p_profile_id)
                )
            )
    ORDER BY pt.priority ASC
    LIMIT 1;

    RETURN v_datavlue;
END;

