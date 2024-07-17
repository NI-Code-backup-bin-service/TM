--multiline
CREATE PROCEDURE get_tid_count_filtered(IN p_search_term varchar(255))
BEGIN
    SELECT     COUNT(*)
    FROM tid AS t
             LEFT JOIN tid_site AS ts ON ts.tid_id = t.tid_id
             LEFT JOIN site AS s ON s.site_id = ts.site_id
             LEFT JOIN (
        site_profiles AS tp2
            JOIN profile AS p2
            ON p2.profile_id = tp2.profile_id
            JOIN profile_type AS pt2
            ON pt2.profile_type_id = p2.profile_type_id AND pt2.priority = 2
        ) ON tp2.site_id = s.site_id
    WHERE (upper(t.tid_id) LIKE CONCAT('%', p_search_term, '%')
        OR upper(t.serial) LIKE CONCAT('%', p_search_term, '%')
        OR upper(p2.name) LIKE CONCAT('%', p_search_term, '%'));
END