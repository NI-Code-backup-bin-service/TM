--multiline
CREATE PROCEDURE get_tid_count_filtered(IN p_search_term varchar(255), IN p_acquirers text)
BEGIN
    SELECT     COUNT(*)
    FROM tid AS t
             LEFT JOIN tid_site AS ts ON ts.tid_id = t.tid_id
             LEFT JOIN site AS s ON s.site_id = ts.site_id
        -- Site Profiles
             LEFT JOIN (site_profiles AS sp
        JOIN profile AS p ON sp.profile_id = p.profile_id AND p.profile_type_id = 4
        ) ON s.site_id = sp.site_id

             LEFT JOIN profile_data AS pd ON pd.profile_id = p.profile_id AND pd.data_element_id = 1
        -- Acquirer Profiles
             LEFT JOIN (site_profiles AS sp2
        JOIN profile AS p2 ON sp2.profile_id = p2.profile_id AND p2.profile_type_id = 2
        ) ON s.site_id = sp2.site_id

        -- Chain
             LEFT JOIN (site_profiles AS sp3
        JOIN profile AS p3 ON sp3.profile_id = p3.profile_id AND p3.profile_type_id = 3
        ) ON s.site_id = sp3.site_id
    WHERE (upper(t.tid_id) LIKE CONCAT('%', p_search_term, '%')
        OR upper(t.serial) LIKE CONCAT('%', p_search_term, '%')
        OR upper(p2.name) LIKE CONCAT('%', p_search_term, '%'))
        AND FIND_IN_SET(p2.name, p_acquirers);
END