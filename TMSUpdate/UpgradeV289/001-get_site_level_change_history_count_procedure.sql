--multiline;
CREATE PROCEDURE `get_site_level_change_history_count`(
    IN profileID INT,
    IN siteId INT
)
BEGIN
SELECT COUNT(*)
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
    );

END;
