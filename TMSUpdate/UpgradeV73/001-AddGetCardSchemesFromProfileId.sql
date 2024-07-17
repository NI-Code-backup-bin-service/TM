--multiline
create procedure get_card_schemes_from_profile_id(IN profileId int)
BEGIN
    SELECT
        pd.datavalue
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
END;