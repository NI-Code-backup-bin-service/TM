--multiline;
CREATE PROCEDURE `add_site_data_groups_to_tid_profiles`(IN profileID INT)
BEGIN
    INSERT ignore INTO profile_data_group (profile_id, data_group_id, version, updated_at, updated_by, created_at, created_by)
    SELECT
    P.tid_profile_id, D.data_group_id, 1, CURRENT_TIMESTAMP, NULL, CURRENT_TIMESTAMP, NULL
    FROM
        (SELECT ts.tid_profile_id
        FROM tid_site ts
        LEFT JOIN site_profiles sp
            ON sp.site_id = ts.site_id
        WHERE sp.profile_id = profileID
        AND ts.tid_profile_id IS NOT NULL) P
    CROSS JOIN
        (SELECT data_group_id
        FROM profile_data_group
        WHERE profile_id = profileID) D;
END;