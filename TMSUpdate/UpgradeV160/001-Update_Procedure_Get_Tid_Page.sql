--multiline
CREATE PROCEDURE get_tid_page(IN p_search_term varchar(255), IN p_offset int, IN p_pageSize int)
BEGIN
SELECT     t.tid_id         		AS 'tid_id',
        t.serial         		AS 'serial',
        t.PIN            		AS 'PIN',
        t.ExpiryDate     		AS 'expiry_date',
        t.reset_pin      		AS 'reset_pin',
        t.reset_pin_expiry_date	AS 'reset_pin_expiry_date',
        t.ActivationDate 		AS 'activation_date',
        t.Presence       		AS 'presence',
        s.site_id        		AS 'site_id',
        p2.name         			AS  'site_name',
        pd.datavalue    			AS 'merchant_id',
        prof.name       			AS 'acquirer'
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
         LEFT JOIN (
    profile_data as pd
        INNER JOIN data_element as de
        ON de.data_element_id = pd.data_element_id
            AND de.name = 'merchantNo'
    ) ON pd.profile_id = tp2.profile_id
         LEFT JOIN (
    site_profiles AS sp
        JOIN profile AS prof
        ON prof.profile_id = sp.profile_id
        JOIN profile_type AS pt
        ON pt.profile_type_id = prof.profile_type_id
            AND pt.priority = 4
    ) ON sp.site_id = s.site_id
WHERE (upper(t.tid_id) LIKE CONCAT('%', p_search_term, '%')
    OR upper(t.serial) LIKE CONCAT('%', p_search_term, '%')
    OR upper(p2.name) LIKE CONCAT('%', p_search_term, '%'))
ORDER BY t.tid_id ASC
    LIMIT p_offset, p_pageSize;
END

