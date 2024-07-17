--multiline
create
    procedure get_tid_page(IN p_search_term varchar(255), IN p_acquirers text)
BEGIN
    SELECT t.tid_id         AS 'tid_id',
           t.serial         AS 'serial',
           t.PIN            AS 'PIN',
           t.ExpiryDate     AS 'expiry_date',
           t.ActivationDate AS 'activation_date',
           t.Presence       AS 'presence',
           s.site_id        AS 'site_id',
           pdn.datavalue    as 'site_name',
           pd.datavalue     as 'merchant_id'

    FROM tid AS t
             LEFT JOIN tid_site AS ts ON ts.tid_id = t.tid_id
             LEFT JOIN site AS s ON s.site_id = ts.site_id

             LEFT JOIN (
        site_profiles AS tp2
            JOIN profile AS p2
            ON p2.profile_id = tp2.profile_id
            JOIN profile_type AS pt2
            ON pt2.profile_type_id = p2.profile_type_id
                AND pt2.priority = 2

        ) ON tp2.site_id = s.site_id
             LEFT JOIN profile_data AS pd
                       ON pd.profile_id = p2.profile_id
                           AND pd.data_element_id =
                               (SELECT de.data_element_id FROM data_element AS de WHERE de.name = 'merchantNo')
                           AND pd.version = (SELECT MAX(d.version)
                                             FROM profile_data AS d
                                             WHERE d.data_element_id = (SELECT de.data_element_id
                                                                        FROM data_element AS de
                                                                        WHERE de.name = 'merchantNo')
                                               AND d.profile_id = p2.profile_id
                                               AND d.approved = 1)
             LEFT JOIN profile_data AS pdn
                       ON pdn.profile_id = p2.profile_id
                           AND pdn.data_element_id =
                               (SELECT de.data_element_id FROM data_element AS de WHERE de.name = 'name')
                           AND pdn.version = (SELECT MAX(d.version)
                                              FROM profile_data AS d
                                              WHERE d.data_element_id = (SELECT de.data_element_id
                                                                         FROM data_element AS de
                                                                         WHERE de.name = 'name')
                                                AND d.profile_id = p2.profile_id
                                                AND d.approved = 1)
             LEFT JOIN (
        site_profiles AS tp4
            JOIN profile AS p4
            ON p4.profile_id = tp4.profile_id
            JOIN profile_type AS pt4
            ON pt4.profile_type_id = p4.profile_type_id
                AND pt4.priority = 4

        ) ON tp4.site_id = s.site_id

    WHERE ts.tid_id = t.tid_id
      AND FIND_IN_SET(p4.name, p_acquirers)
      AND (upper(t.tid_id) LIKE CONCAT('%', p_search_term, '%')
        OR upper(t.serial) LIKE CONCAT('%', p_search_term, '%')
        OR upper(pdn.datavalue) LIKE CONCAT('%', p_search_term, '%'));
END;
