--multiline
create procedure get_available_schemes_for_profile_id(IN profileId int, IN includeQrSchemes tinyint(1))
BEGIN
    DECLARE v_CardDefinitionsJson TEXT;
    DECLARE v_CardDefinitionsCount INT;
    DECLARE v_CardName VARCHAR(255);
    DECLARE v_CardDefinitionsIterator INT DEFAULT 0;
    DECLARE v_isTid BOOL DEFAULT FALSE;

    /*Create the temporary table we will use to generate our result set*/
    CREATE TEMPORARY TABLE tt_AvailableSchemes(
                                                  scheme_id INT PRIMARY KEY,
                                                  scheme_name VARCHAR(255)
    );
    SELECT if(count(*) > 0, TRUE, FALSE) INTO v_isTid FROM tid_site where tid_profile_id = profileId;
    /*Get the correct card definition JSON*/
    SELECT
        pd.datavalue INTO v_CardDefinitionsJson
    FROM profile_data pd
             INNER JOIN profile p ON
            p.profile_id = pd.profile_id
             INNER JOIN profile_type pt ON
            p.profile_type_id = pt.profile_type_id
             INNER JOIN data_element de ON
            pd.data_element_id = de.data_element_id
    WHERE
            upper(de.name) = 'CARDDEFINITIONS'
      AND
            pd.datavalue != ''
      AND
        pd.datavalue IS NOT NULL
      AND
        (
                    p.profile_id = profileId
                OR
                    UPPER(pt.name) != 'SITE'
            )
    ORDER BY pt.priority ASC
    LIMIT 1;

    /*How many card definitions are there?*/
    SELECT JSON_LENGTH(v_CardDefinitionsJson) INTO v_CardDefinitionsCount;
    /*Loop over all of the cards*/
    CardDefinitionsLoop: LOOP
        /*Loop exit condition*/
        IF v_CardDefinitionsIterator = v_CardDefinitionsCount THEN
            LEAVE CardDefinitionsLoop;
        end if;
        SELECT NULL INTO v_CardName;
        SELECT JSON_UNQUOTE(JSON_EXTRACT(v_CardDefinitionsJson,CONCAT('$[', v_CardDefinitionsIterator, '].cardName'))) INTO v_CardName;
        /*Perform the lookup against that cardName and add a row to our temp table*/
        IF v_CardName IS NOT NULL THEN
            INSERT INTO tt_AvailableSchemes (
                scheme_id,
                scheme_name
            )
            SELECT
                s.scheme_id,
                s.scheme_name
            FROM schemes s
            WHERE
                    UPPER(s.scheme_name) = v_CardName;
        END IF;

        /*Increment the iterator*/
        SET v_CardDefinitionsIterator = v_CardDefinitionsIterator + 1;
    end loop;

    /*
    Now we need to handle QR schemes.
    QR Schemes are made available by profile data groups and thus we need to do this in a separate step.
    */
    IF includeQrSchemes THEN
        INSERT INTO tt_AvailableSchemes (
            scheme_id,
            scheme_name
        )
        SELECT
            scheme_id,
            scheme_name
        FROM schemes
        WHERE scheme_name IN (
            SELECT UPPER(dg.name)
            FROM profile_data_group
                     INNER JOIN data_group dg ON
                    profile_data_group.data_group_id = dg.data_group_id
            WHERE profile_id = if(v_isTid, GetProfileParentId(profileId), profileId)
        );
    END IF;

    /*select everything from our temporary table to return*/
    SELECT * FROM tt_AvailableSchemes;
    DROP TEMPORARY TABLE tt_AvailableSchemes;
END;

