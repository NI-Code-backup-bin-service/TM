--multiline;
CREATE PROCEDURE `get_site_level_change_history`(
    IN profileID INT,
    IN siteId INT,
    IN offset_value INT,
    IN limit_value INT
)
BEGIN
SELECT
    CONCAT(
            IFNULL(dg.name, ''),
            '/',
            IFNULL(de.displayname_en, de.name)
        ) AS display_name,
    a.change_type AS change_type,
    a.current_value AS original_value,
    a.new_value,
    a.created_by AS updated_by,
    a.created_at AS updated_at,
    a.approved,
    a.tid_id,
    a.is_password,
    a.is_encrypted
FROM approvals a
         LEFT JOIN data_element de ON de.data_element_id = a.data_element_id
         LEFT JOIN data_group dg ON dg.data_group_id = de.data_group_id
WHERE
    (a.profile_id = profileID)
   OR EXISTS (
        SELECT 1
        FROM tid_site ts
                 JOIN site_profiles sp ON ts.site_id = sp.site_id
        WHERE ts.tid_id = a.tid_id AND sp.profile_id = profileID
    )
ORDER BY a.created_at DESC
    LIMIT limit_value
OFFSET offset_value;
END;