--multiline;
CREATE PROCEDURE `update_eod_auto`(IN profileId INT, IN auto BOOLEAN)
BEGIN
    SET @profileTypeId = (SELECT profile_type_id from profile where profile_id=profileId);

    IF (@profileTypeId=4) THEN

        UPDATE tid t
        JOIN tid_site ts
            ON t.tid_id = ts.tid_id
        JOIN site_profiles sp
            ON sp.site_id = ts.site_id
        SET t.eod_auto = auto
        WHERE sp.profile_id = profileId;

    ELSEIF (@profileTypeId=3) THEN

        UPDATE tid t
        JOIN tid_site ts
            ON t.tid_id = ts.tid_id
        JOIN site_profiles sp
            ON sp.site_id = ts.site_id
        SET t.eod_auto = auto
        WHERE sp.profile_id = profileId
        AND ts.tid_id NOT IN (
            SELECT ts.tid_id
            FROM tid_site ts
            JOIN site_profiles sp
                ON sp.site_id = ts.site_id
            JOIN profile p
                ON p.profile_id=sp.profile_id
            JOIN profile_type pt
                ON pt.profile_type_id=p.profile_type_id
            JOIN profile_data pd
                ON pd.profile_id=p.profile_id
            WHERE sp.site_id IN (SELECT site_id
                FROM site_profiles
                WHERE profile_id = profileId)
            AND pt.profile_type_id = 4
            AND pd.data_element_id=(SELECT de.data_element_id
                FROM data_element de
                JOIN data_group dg
                    ON dg.data_group_id=de.data_group_id
                WHERE de.name = "auto"
                AND dg.name = "endOfDay"));

    ELSEIF (@profileTypeId=2) THEN

        UPDATE tid t
        JOIN tid_site ts
            ON t.tid_id = ts.tid_id
        JOIN site_profiles sp
            ON sp.site_id = ts.site_id
        SET t.eod_auto = auto
        WHERE sp.profile_id = profileId
        AND ts.tid_id NOT IN (
            SELECT ts.tid_id
            FROM tid_site ts
            JOIN site_profiles sp
                ON sp.site_id = ts.site_id
            JOIN profile p
                ON p.profile_id=sp.profile_id
            JOIN profile_type pt
                ON pt.profile_type_id=p.profile_type_id
            JOIN profile_data pd
                ON pd.profile_id=p.profile_id
            WHERE sp.site_id IN (SELECT site_id
                FROM site_profiles
                WHERE profile_id = profileId)
            AND pt.profile_type_id IN(3,4)
            AND pd.data_element_id=(SELECT de.data_element_id
                FROM data_element de
                JOIN data_group dg
                    ON dg.data_group_id=de.data_group_id
                WHERE de.name = "auto"
                AND dg.name = "endOfDay"));

    END IF;
END