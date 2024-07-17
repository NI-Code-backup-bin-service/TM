--multiline
create
    procedure get_site_page(IN p_search_term varchar(255), IN p_acquirers text)
BEGIN
    SELECT t.site_id     AS 'site_id',
           p2.profile_id AS 'site_profile_id',
           pdn.datavalue AS 'site_name',
           p3.profile_id AS 'chain_profile_id',
           p3.name          'chain_name',
           p4.profile_id AS 'acquirer_profile_id',
           p4.name       AS 'acquirer_name',
           p5.profile_id AS 'global_profile_id',
           p5.name       AS 'global_name',
           pd.datavalue  AS 'merchant_id'

    FROM site AS t
             LEFT JOIN (
        site_profiles AS tp1
            JOIN profile AS p1
            ON p1.profile_id = tp1.profile_id
            JOIN profile_type AS pt1
            ON pt1.profile_type_id = p1.profile_type_id
                AND pt1.priority = 1

        ) ON tp1.site_id = t.site_id
             LEFT JOIN (
        site_profiles AS tp2
            JOIN profile AS p2
            ON p2.profile_id = tp2.profile_id
            JOIN profile_type AS pt2
            ON pt2.profile_type_id = p2.profile_type_id
                AND pt2.priority = 2

        ) ON tp2.site_id = t.site_id

        -- Get the merchant number
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
        -- Get the merchant name
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
        site_profiles AS tp3
            JOIN profile AS p3
            ON p3.profile_id = tp3.profile_id
            JOIN profile_type AS pt3
            ON pt3.profile_type_id = p3.profile_type_id
                AND pt3.priority = 3

        ) ON tp3.site_id = t.site_id
             LEFT JOIN (
        site_profiles AS tp4
            JOIN profile AS p4
            ON p4.profile_id = tp4.profile_id
            JOIN profile_type AS pt4
            ON pt4.profile_type_id = p4.profile_type_id
                AND pt4.priority = 4

        ) ON tp4.site_id = t.site_id
             LEFT JOIN (
        site_profiles AS tp5
            JOIN profile AS p5
            ON p5.profile_id = tp5.profile_id
            JOIN profile_type AS pt5
            ON pt5.profile_type_id = p5.profile_type_id
                AND pt5.priority = 5

        ) ON tp5.site_id = t.site_id
    WHERE (upper(pdn.datavalue) LIKE CONCAT('%', p_search_term, '%')
        OR upper(p3.name) LIKE CONCAT('%', p_search_term, '%')
        OR upper(p4.name) LIKE CONCAT('%', p_search_term, '%')
        OR upper(p5.name) LIKE CONCAT('%', p_search_term, '%')
        OR pd.datavalue LIKE CONCAT('%', p_search_term, '%')
        )
      AND FIND_IN_SET(p4.name, p_acquirers);
END;
