--multiline
CREATE PROCEDURE get_chain_page(IN p_search_term varchar(255), IN p_acquirers text)
BEGIN
    SELECT DISTINCT p.profile_id AS 'chain_profile_id',
            p.name       AS 'chain_name',
            p2.name      AS 'acquirer_name',
            (SELECT COUNT(tid_id) FROM tid_site ts, site_profiles sp WHERE ts.site_id=sp.site_id AND sp.profile_id = p.profile_id AND sp.profile_id = cp.chain_profile_id) AS 'chain_tid_count'
    FROM profile AS p
        LEFT JOIN chain_profiles AS cp ON cp.chain_profile_id = p.profile_id
        LEFT JOIN profile AS p2 ON p2.profile_id = cp.acquirer_id
    WHERE FIND_IN_SET(p2.name, p_acquirers) AND (upper(p.name) LIKE CONCAT('%', p_search_term, '%') OR p.profile_id like CONCAT('%', p_search_term, '%') OR p2.name like CONCAT('%', p_search_term, '%'));
END;
