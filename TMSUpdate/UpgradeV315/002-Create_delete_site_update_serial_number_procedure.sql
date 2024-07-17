--multiline;
CREATE PROCEDURE `delete_site_update_serial_number`(IN profileID INT)
BEGIN
    DECLARE siteId INT;
    DECLARE profileExists INT;
    DECLARE tid_ids_list VARCHAR(255); -- Assuming it's a comma-separated string, adjust the size accordingly
    DECLARE v_remaining_elements_count INT;
    -- Fetch site_id based on the provided profileID
SELECT site_id INTO siteId
FROM site_profiles
WHERE profile_id = profileID;
-- Fetch tid_ids_list based on the siteId
SELECT GROUP_CONCAT(t.tid_id) INTO tid_ids_list
FROM tid t
         LEFT JOIN tid_site ts ON ts.tid_id = t.tid_id
WHERE ts.site_id = siteId;
-- Check if profile exists
SET profileExists = (
        SELECT tid_profile_id
        FROM tid_site
        WHERE tid_id IN (SELECT tid_id FROM tid_ids_list)
        LIMIT 1
    );
START TRANSACTION;
SET @profileID = (SELECT sp.profile_id
                     FROM site_profiles sp
                     LEFT JOIN profile p ON p.profile_id = sp.profile_id
                     LEFT JOIN profile_type pt ON p.profile_type_id = pt.profile_type_id
                     WHERE sp.site_id = siteId
                     ORDER BY pt.priority
                     LIMIT 1);
DELETE FROM profile_data WHERE profile_id = @profileID AND @profileID IS NOT NULL;
DELETE FROM profile_data_group WHERE profile_id = @profileID AND @profileID IS NOT NULL;
DELETE FROM site_profiles WHERE site_id = siteId AND siteId IS NOT NULL;
DELETE FROM tid_site WHERE site_id = siteId AND siteId IS NOT NULL;
DELETE FROM site_level_users WHERE site_id = siteId AND siteId IS NOT NULL;
DELETE FROM site WHERE site_id = siteId AND siteId IS NOT NULL;
DELETE u FROM tid_user_override u LEFT OUTER JOIN tid_site ts ON ts.tid_id = u.tid_id WHERE ts.tid_id IS NULL;
    DELETE t FROM tid t LEFT OUTER JOIN tid_site ts ON ts.tid_id = t.tid_id WHERE ts.tid_id IS NULL;
DELETE FROM velocity_limits_txn
WHERE velocity_limit_id IN (
    SELECT velocity_limit_id
    FROM velocity_limits
    WHERE site_id = siteId AND limit_level IN (1, 3) AND tid_id = -1
);
DELETE FROM velocity_limits WHERE site_id = siteId AND limit_level IN (1, 3) AND tid_id = -1;
DELETE FROM tid_site WHERE tid_id IN (SELECT CAST(tid_id AS SIGNED) FROM tid_ids_list WHERE FIND_IN_SET(tid_id, tid_ids_list));
DELETE FROM tid WHERE tid_id IN (SELECT CAST(tid_id AS SIGNED) FROM tid_ids_list WHERE FIND_IN_SET(tid_id, tid_ids_list));
DELETE FROM tid_updates WHERE tid_id IN (SELECT CAST(tid_id AS SIGNED) FROM tid_ids_list WHERE FIND_IN_SET(tid_id, tid_ids_list));
IF profileExists IS NOT NULL THEN
        DELETE pd.*
        FROM profile_data pd
        INNER JOIN data_element de ON pd.data_element_id = de.data_element_id
        INNER JOIN data_element_locations_data_element delde ON de.data_element_id = delde.data_element_id
        INNER JOIN data_element_locations del ON delde.location_id = del.location_id
        WHERE pd.profile_id = @profileID AND del.location_name = 'fraud';
DELETE
FROM profile_data_group
WHERE profile_id = @profileID
  AND data_group_id NOT IN (
    SELECT de.data_group_id
    FROM profile_data pd
             INNER JOIN data_element de ON pd.data_element_id = de.data_element_id
    WHERE pd.profile_id = @profileID
);
-- Remaining rows for the given profile in profile_data
SELECT COUNT(*) INTO v_remaining_elements_count
FROM profile_data pd
WHERE pd.profile_id = @profileID;
-- If there's no remaining rows in profile data, then delete approvals and update tid_site
IF v_remaining_elements_count = 0 THEN
DELETE FROM approvals WHERE profile_id = @profileID AND approved = 0;
UPDATE tid_site SET tid_profile_id = NULL, updated_at = NOW() WHERE tid_profile_id = @profileID;
END IF;
END IF;
COMMIT;
END;