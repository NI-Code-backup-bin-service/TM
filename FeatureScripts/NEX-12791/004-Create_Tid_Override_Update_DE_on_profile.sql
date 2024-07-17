--multiline
CREATE PROCEDURE `Create_Tid_Override_Update_DE_on_profile`(in TID text, MID text, cashDeskId text, city text)
 BEGIN
    DECLARE ErrorMsg MEDIUMTEXT;
    SET @CashDeskId = "cashDeskId";
    SET @dataGroupName = "IPP";

    SET @profile_id = (Select profile_id from profile_data pd where pd.data_element_id = (Select data_element_id from data_element where name = "merchantNo") AND
						pd.datavalue = MID);
    
    SET @site_id = (Select site_id from site_profiles where profile_id = @profile_id);

    SET @tid = (Select tid_id FROM tid_site WHERE tid_id = TID);

    IF @tid IS NULL THEN
        -- Set error message for tid_id
        SIGNAL SQLSTATE '45000'
        SET MESSAGE_TEXT = 'TID is not found';
    END IF;
	IF @site_id IS NULL THEN
        -- Set error message for site_id
        SIGNAL SQLSTATE '45000'
        SET MESSAGE_TEXT = 'SiteID is not found';
	END IF;
    SET @profileExists := (SELECT IFNULL( (Select tid_profile_id FROM tid_site WHERE tid_id = TID AND site_id = @site_id) ,0));
    
    # Create the override if not exist...
    IF @profileExists = 0 THEN
            SET @profileTypeID := (SELECT profile_type_id FROM profile_type WHERE `name` = "tid");
            INSERT into profile (profile_type_id, name, version, updated_at, updated_by, created_at, created_by) values (@profileTypeID, TID, 1, NOW(), 'system', NOW(), 'system');
            SET @tidProfileID = LAST_INSERT_ID();
            UPDATE tid_site SET tid_profile_id = @tidProfileID where tid_id = TID and site_id = @site_id;

            INSERT ignore INTO profile_data_group (profile_id, data_group_id, version, updated_at, updated_by, created_at, created_by)
            SELECT all_groups.tid_profile_id, all_groups.data_group_id, 1, CURRENT_TIMESTAMP, NULL, CURRENT_TIMESTAMP, NULL
            FROM (SELECT dg.data_group_id, ts.tid_profile_id from data_group dg cross join tid_site ts where ts.tid_profile_id is not null) as all_groups
                left join (select pdg.data_group_id, p.profile_id from profile_data_group pdg
                            inner join profile p on p.profile_id = pdg.profile_id
                            inner join profile_type pt on pt.profile_type_id = p.profile_type_id where pt.name = 'tid') as already_allocated on all_groups.data_group_id = already_allocated.data_group_id and all_groups.tid_profile_id = already_allocated.profile_id where already_allocated.data_group_id is null;

            INSERT INTO profile_data (profile_id, data_element_id, datavalue, version, updated_by, created_by, approved, overriden, is_encrypted)
                (WITH profileData AS (SELECT pd.data_element_id, pd.datavalue, pd.is_encrypted, ROW_NUMBER() OVER(PARTITION BY pd.data_element_id ORDER BY pt.priority) AS RowNum
                                      FROM profile_data pd
                                     JOIN data_element de
                                          ON de.data_element_id=pd.data_element_id
                                     JOIN site_profiles sp
                                          ON pd.profile_id = sp.profile_id
                                     JOIN profile p
                                          ON p.profile_id = sp.profile_id
                                     JOIN profile_type pt
                                          ON p.profile_type_id = pt.profile_type_id
                                    WHERE sp.site_id = @site_id
                                    AND de.tid_overridable=1
                                    AND de.data_group_id IN
                                        (SELECT DISTINCT(de.data_group_id)
                                        FROM profile_data pd
                                        JOIN site_profiles sp
                                             ON pd.profile_id= sp.profile_id
                                        JOIN profile p
                                             ON p.profile_id=sp.profile_id
                                        JOIN profile_type pt
                                             ON pt.profile_type_id=p.profile_type_id
                                        JOIN data_element de
                                             ON de.data_element_id=pd.data_element_id
                                        WHERE sp.site_id=@site_id)
                )
                 SELECT @tidProfileID, data_element_id, datavalue, 0, 'system', 'system', 1, 0, is_encrypted
                 FROM profileData
                 WHERE RowNum = 1);

            SET @acquirer = (SELECT DISTINCT p4.name FROM profile p
                                                              JOIN tid_site ts ON ts.tid_profile_id = profile_id
                                                              JOIN site t ON t.site_id = ts.site_id
                                                              JOIN (site_profiles tp4
                JOIN profile p4 ON p4.profile_id = tp4.profile_id
                JOIN profile_type pt4 ON pt4.profile_type_id = p4.profile_type_id AND pt4.priority = 4)
                                                                   ON tp4.site_id = t.site_id
                             WHERE p.profile_type_id = (SELECT profile_type_id FROM profile_type WHERE profile_type.name = "tid")
                             AND p.profile_id = @tidProfileID);

            SET @createChangeType := (SELECT approval_type_id FROM approval_type WHERE approval_type_name = "Create" LIMIT 1);
            SET @change_type := 5;
            SET @dataValue := "Override Created";
            SET @approved := 1;

            INSERT INTO approvals(profile_id, data_element_id, change_type, current_value, new_value, created_at, approved, created_by, approved_by, approved_at, tid_id, acquirer)
            VALUES (@tidProfileID, 1, @change_type, '', @dataValue, NOW(), @approved, 'system', 'system', NOW(), TID, @acquirer);

            INSERT into profile_data (profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden, is_encrypted)
            values (@tidProfileID, (SELECT data_element_id FROM data_element WHERE name = @CashDeskId AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = @dataGroupName)), cashDeskId, 1, NOW(), 'system', NOW(), 'system', 1, 1, 0)
            ON DUPLICATE KEY
                UPDATE datavalue=VALUES(datavalue), updated_at = VALUES(updated_at), updated_by = VALUES(updated_by);

            INSERT into profile_data (profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden, is_encrypted)
            values (@tidProfileID, (SELECT data_element_id FROM data_element WHERE name = "city" AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = "store")), city, 1, NOW(), 'system', NOW(), 'system', 1, 1, 0)
            ON DUPLICATE KEY
                UPDATE datavalue=VALUES(datavalue), updated_at = VALUES(updated_at), updated_by = VALUES(updated_by);

        INSERT INTO profile_data (
            profile_id,
            data_element_id,
            datavalue,
            version,
            updated_at,
            updated_by,
            created_at,
            created_by,
            approved,
            overriden,
            is_encrypted
        )
        VALUES (
            @tidProfileID,
            (
                SELECT data_element_id
                FROM data_element
                WHERE name = 'active'
                    AND data_group_id = (
                        SELECT data_group_id
                        FROM data_group
                        WHERE name = 'modules'
                    )
            ),
            JSON_ARRAY('IPP'),
            1,
            NOW(),
            'system',
            NOW(),
            'system',
            1,
            1,
            0
        )
        ON DUPLICATE KEY UPDATE
        datavalue = JSON_ARRAY_APPEND(IFNULL(datavalue, '[]'), '$', 'IPP'),
        updated_at = VALUES(updated_at),
        updated_by = VALUES(updated_by);
    ELSE
        INSERT ignore INTO profile_data_group (
            profile_id, data_group_id, version,
            updated_at, updated_by, created_at,
            created_by
        )
        SELECT
            all_groups.tid_profile_id,
            all_groups.data_group_id,
            1,
            CURRENT_TIMESTAMP,
            NULL,
            CURRENT_TIMESTAMP,
            NULL
        FROM
            (
                SELECT
                    dg.data_group_id,
                    ts.tid_profile_id
                from
                    data_group dg cross
                                      join tid_site ts
                where
                    ts.tid_profile_id is not null
            ) as all_groups
                left join (
                select
                    pdg.data_group_id,
                    p.profile_id
                from
                    profile_data_group pdg
                        inner join profile p on p.profile_id = pdg.profile_id
                        inner join profile_type pt on pt.profile_type_id = p.profile_type_id
                where
                    pt.name = 'tid'
            ) as already_allocated on all_groups.data_group_id = already_allocated.data_group_id
                and all_groups.tid_profile_id = already_allocated.profile_id
        where
            already_allocated.data_group_id is null;

        SET @tidProfileID :=(Select tid_profile_id FROM tid_site WHERE tid_id = TID);

        INSERT into profile_data (profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden, is_encrypted)
        values (@tidProfileID, (SELECT data_element_id FROM data_element WHERE name = @CashDeskId AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = @dataGroupName)), cashDeskId, 1, NOW(), 'system', NOW(), 'system', 1, 1, 0)
        ON DUPLICATE KEY
            UPDATE datavalue=VALUES(datavalue), updated_at = VALUES(updated_at), updated_by = VALUES(updated_by);

        INSERT into profile_data (profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden, is_encrypted)
        values (@tidProfileID, (SELECT data_element_id FROM data_element WHERE name = "city" AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = "store")), city, 1, NOW(), 'system', NOW(), 'system', 1, 1, 0)
        ON DUPLICATE KEY
            UPDATE datavalue=VALUES(datavalue), updated_at = VALUES(updated_at), updated_by = VALUES(updated_by);

       # Set active Module IPP value only if it not exist...
       SET @activeIppModuleExists := ( SELECT COUNT(*) FROM profile_data WHERE profile_id = @tidProfileID AND data_element_id = (SELECT data_element_id FROM data_element
                    WHERE name = 'active' AND data_group_id = ( SELECT data_group_id FROM data_group WHERE name = 'modules')) AND JSON_SEARCH(datavalue, 'one', 'IPP') IS NOT NULL);

        IF @activeIppModuleExists = 0 THEN
            INSERT INTO profile_data (
                profile_id,
                data_element_id,
                datavalue,
                version,
                updated_at,
                updated_by,
                created_at,
                created_by,
                approved,
                overriden,
                is_encrypted
            )
            VALUES (
                @tidProfileID,
                (
                    SELECT data_element_id
                    FROM data_element
                    WHERE name = 'active'
                        AND data_group_id = (
                            SELECT data_group_id
                            FROM data_group
                            WHERE name = 'modules'
                        )
                ),
                JSON_ARRAY('IPP'),
                1,
                NOW(),
                'system',
                NOW(),
                'system',
                1,
                1,
                0
            )
            ON DUPLICATE KEY UPDATE
            datavalue = JSON_ARRAY_APPEND(IFNULL(datavalue, '[]'), '$', 'IPP'),
            updated_at = VALUES(updated_at),
            updated_by = VALUES(updated_by);
        END IF;
    END IF;
END