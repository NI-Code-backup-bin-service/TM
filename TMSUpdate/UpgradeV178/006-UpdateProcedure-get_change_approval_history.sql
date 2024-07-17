--multiline
CREATE PROCEDURE get_change_approval_history(IN afterDate varchar(255), IN name varchar(255), IN user varchar(255), IN beforeDate varchar(255), IN field varchar(255), IN acquirers text)
BEGIN
    SET @name = upper(concat('%', ifnull(name, ''), '%'));
    SET @user = upper(concat('%', ifnull(user, ''), '%'));
    SET @field = upper(concat('%', ifnull(field, ''), '%'));

    SELECT p.name,
           CONCAT(dg.name, '/', de.name) AS 'field',
           a.current_value               AS original_value,
           a.new_value                   AS updated_value,
           a.created_by                  AS updated_by,
           a.created_at                  AS updated_at,
           a.approved,
           a.approved_by                 AS reviewed_by,
           a.approved_at                 AS reviewed_at,
           a.tid_id,
           a.merchant_id,
           a.is_password,
           a.is_encrypted,
           a.change_type
    FROM approvals a
             LEFT JOIN data_element de ON de.data_element_id = a.data_element_id
             LEFT JOIN data_group dg ON dg.data_group_id = de.data_group_id
             LEFT JOIN profile p ON p.profile_id = a.profile_id

    WHERE a.approved != 0
      AND (
            afterDate IS NULL
            OR a.created_at >= afterDate
        )
      AND p.name LIKE @name
      AND a.approved_by LIKE @user
      AND (
            beforeDate IS NULL
            OR a.created_at <= beforeDate
        )
      AND de.name LIKE @field
      AND (
        FIND_IN_SET(a.acquirer, acquirers)
        )

    GROUP BY p.name,
             dg.name,
             de.name,
             original_value,
             updated_value,
             updated_by,
             updated_at,
             a.approved,
             reviewed_by,
             reviewed_at,
             a.tid_id,
             a.merchant_id,
             a.is_password,
             a.is_encrypted,
             a.change_type
    ORDER BY a.approved_at DESC;
END;