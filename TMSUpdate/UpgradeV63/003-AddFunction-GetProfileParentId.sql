--multiline
CREATE FUNCTION GetProfileParentId(profileId INT) RETURNS INT NOT DETERMINISTIC
BEGIN
    DECLARE parentProfileId INT;
    SELECT
        CASE
            WHEN pt.name = 'global' THEN
                (SELECT profile_id FROM profile WHERE profile.profile_type_id = pt.profile_type_id)
            WHEN pt.name = 'site' THEN
                (
                    SELECT p.profile_id
                    FROM site_profiles sp
                             LEFT JOIN tid_site ts ON
                            ts.site_id = sp.site_id
                             INNER JOIN profile p ON
                            sp.profile_id = p.profile_id
                    WHERE
                            ts.tid_profile_id = profileId
                      AND
                            p.profile_type_id = pt.profile_type_id
                )
            ELSE
                (
                    SELECT sp.profile_id
                    FROM site_profiles sp
                             INNER JOIN profile p ON
                            sp.profile_id = p.profile_id
                    WHERE
                            sp.site_id = (SELECT site_id FROM site_profiles WHERE profile_id = profileId LIMIT 1)
                      AND
                            p.profile_type_id = pt.profile_type_id
                )
            END INTO parentProfileId
    FROM profile_type pt
    WHERE pt.priority > (
        SELECT pt.priority
        FROM profile_type pt
                 INNER JOIN profile p ON
                pt.profile_type_id = p.profile_type_id
        WHERE p.profile_id = profileId
    )
    ORDER BY priority ASC
    LIMIT 1;

    RETURN parentProfileId;
END
