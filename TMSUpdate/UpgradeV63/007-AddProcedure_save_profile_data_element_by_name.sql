--multiline
create procedure save_profile_data_element_by_name(IN p_profileId int,
                                                     IN p_elementName varchar(255),
                                                     IN p_elementValue text)
BEGIN
    DECLARE v_existingElementCount INT;
    DECLARE v_parentElementCount INT;

    SELECT
        COUNT(*) INTO v_existingElementCount
    FROM profile_data pd
             INNER JOIN data_element de ON
            pd.data_element_id = de.data_element_id
    WHERE
            pd.profile_id = p_profileId
      AND
            de.name = p_elementName;

    SELECT
        COUNT(*) INTO v_parentElementCount
    FROM profile_data pd
             INNER JOIN data_element de ON
            pd.data_element_id = de.data_element_id
    WHERE
            pd.profile_id = GetProfileParentId(p_profileId)
      AND
            de.name = p_elementName;

    IF v_existingElementCount > 0 THEN
        UPDATE profile_data
        SET
            datavalue = p_elementValue,
            updated_at = NOW(),
            updated_by = 'system'
        WHERE
                data_element_id = (SELECT data_element_id FROM data_element WHERE name = p_elementName)
          AND
                profile_id = p_profileId;
    ELSE
        INSERT INTO profile_data
        (
            profile_id,
            data_element_id,
            datavalue,
            version,
            updated_at,
            updated_by,
            created_at,
            created_by,
            approved,
            overriden
        )
        VALUES
        (
            p_profileId,
            (SELECT data_element_id FROM data_element WHERE name = p_elementName),
            p_elementValue,
            1,
            NOW(),
            'system',
            NOW(),
            'system',
            1,
            0
        );
    END IF;

    /*are we overriding?*/
    IF p_elementValue != get_parent_datavalue((SELECT data_element_id FROM data_element WHERE name = p_elementName), p_profileId) OR (get_parent_datavalue((SELECT data_element_id FROM data_element WHERE name = p_elementName), p_profileId) IS NULL)  THEN
        IF v_parentElementCount = 0 THEN
            INSERT INTO profile_data
            (
                profile_id,
                data_element_id,
                datavalue,
                version,
                updated_at,
                updated_by,
                created_at,
                created_by,
                approved,
                overriden
            )
            VALUES
            (
                GetProfileParentId(p_profileId),
                (SELECT data_element_id FROM data_element WHERE name = p_elementName),
                '',
                1,
                NOW(),
                'system',
                NOW(),
                'system',
                1,
                1
            );
        ELSE
            UPDATE profile_data
            SET overriden = 1
            WHERE
                    data_element_id = (SELECT data_element_id FROM data_element WHERE name = p_elementName)
              AND
                    profile_id = GetProfileParentId(p_profileId);
        END IF;
    END IF;
END;