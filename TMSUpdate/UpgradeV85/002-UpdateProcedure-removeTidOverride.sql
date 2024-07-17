--multiline
CREATE
    PROCEDURE removeTidOverride(IN overrideId int)
BEGIN
    DECLARE v_nonSiteConfigProfileDataCount INT DEFAULT 0;

    SELECT COUNT(*) INTO v_nonSiteConfigProfileDataCount
    FROM
        profile_data pd
            INNER JOIN profile p ON
                pd.profile_id = p.profile_id
            INNER JOIN data_element_locations_data_element delde ON
                pd.data_element_id = delde.data_element_id
            INNER join data_element_locations del ON
                    delde.location_id = del.location_id
                AND
                    del.profile_type_id = p.profile_type_id
    WHERE
            pd.profile_id = 15
      AND
            del.location_name != 'site_configuration';

    #Only delete the profile id if there's no non-site configuration set
    IF v_nonSiteConfigProfileDataCount = 0 THEN
        UPDATE tid_site SET tid_profile_id = NULL, updated_at = NOW() WHERE tid_profile_id = overrideId;
    END IF;

    #Only delete the profile_data_group entries where they are just site_configuration
    DELETE FROM profile_data_group
    WHERE
            profile_id = overrideId
      AND
            data_group_id NOT IN (
            SELECT pdg.data_group_id
            FROM (SELECT * FROM profile_data_group) pdg
                     INNER JOIN profile p ON
                    pdg.profile_id = p.profile_id
                     INNER JOIN data_element de ON
                    de.data_group_id = pdg.data_group_id
                     INNER JOIN data_element_locations_data_element delde ON
                    de.data_element_id = delde.data_element_id
                     INNER join data_element_locations del ON
                        delde.location_id = del.location_id
                    AND
                        del.profile_type_id = p.profile_type_id
            where
                    pdg.profile_id = 15
              AND
                    del.location_name != 'site_configuration'
        )
      # Do not delete the cleanse time as part of TID override delete
      AND
            del.location_name != 'fraud';


    DELETE FROM profile_data
    WHERE
            profile_id = overrideId
      AND
            profile_data_id NOT IN (
            SELECT pd.profile_data_id
            FROM
                (SELECT * FROM profile_data) pd
                    INNER JOIN profile p ON
                        pd.profile_id = p.profile_id
                    INNER JOIN data_element_locations_data_element delde ON
                        pd.data_element_id = delde.data_element_id
                    INNER join data_element_locations del ON
                            delde.location_id = del.location_id
                        AND
                            del.profile_type_id = p.profile_type_id
            WHERE
                    pd.profile_id = 15
              AND
                    (del.location_name != 'site_configuration' OR del.location_name != 'fraud')
        )
      # Do not delete the cleanse time as part of TID override delete
      AND
            del.location_name != 'fraud';


    DELETE FROM approvals WHERE profile_id = overrideId AND approved = 0;
END

