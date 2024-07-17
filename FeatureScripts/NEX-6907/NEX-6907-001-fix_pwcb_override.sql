--multiline
/*
 * Fix records that were incorrectly migrated by the '002-MigrateCashbackData_and_cleanup.sql' script.
 * Records in 'profile_data' relating to cashback definitions do not have the 'overridden' value set
 * correctly as this column was ignored by said script. This procedure determines if the overridden value
 * should be true for sites with cashback data elements and updates the 'override' column value to true (1) if so.
 */
CREATE PROCEDURE `fix_pwcb_override`()
BEGIN
    /* Get the 'cashback' data group id and the 'definitions' data element id */
    DECLARE cashback_dg_id INT DEFAULT (SELECT data_group_id FROM data_group WHERE name = "cashback" LIMIT 1);
    DECLARE cashback_de_id INT DEFAULT (SELECT data_element_id FROM data_element WHERE name = "definitions" AND data_group_id = cashback_dg_id LIMIT 1);

    DECLARE done INT DEFAULT FALSE;
    DECLARE id INT;

    /* Find all the sites with the cashback data group set. */
    DECLARE cur CURSOR FOR SELECT DISTINCT p.profile_id
                           FROM profile p
                                    INNER JOIN profile_data pd ON p.profile_id = pd.profile_id
                                    INNER JOIN data_element de ON pd.data_element_id = de.data_element_id
                                    INNER JOIN data_group dg ON de.data_group_id = dg.data_group_id
                                    INNER JOIN profile_type pt ON p.profile_type_id = pt.profile_type_id
                           WHERE pt.name = 'site' AND dg.data_group_id = cashback_dg_id;
    DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;
    OPEN cur;

    /*
     * For each site, determine whether it's associated chain or acquirer has the 'definitions' data
     * element set. If either of them do, this means the data element has been overridden at the
     * site level. If so, update the 'overridden' column in profile_data to reflect this.
    */
    read_loop: LOOP
        FETCH cur INTO id;

        IF done THEN
            LEAVE read_loop;
        END IF;

        SET @site_id = (SELECT sp.site_id FROM site_profiles sp WHERE sp.profile_id = id LIMIT 1);

        SET @acquirer_id = (SELECT p.profile_id FROM site_profiles sp
             LEFT JOIN `profile` p ON p.profile_id = sp.profile_id
             LEFT JOIN profile_type pt ON p.profile_type_id = pt.profile_type_id
             WHERE sp.site_id = @site_id AND pt.name = 'acquirer');

        SET @chain_id = (SELECT p.profile_id FROM site_profiles sp
             LEFT JOIN `profile` p ON p.profile_id = sp.profile_id
             LEFT JOIN profile_type pt ON p.profile_type_id = pt.profile_type_id
             WHERE sp.site_id = @site_id AND pt.name = 'chain');

        /* Check if the 'definitions' data element is set at either the acquirer or chain level. */
        SET @override_count = (SELECT COUNT(*) FROM profile_data pd
            WHERE pd.data_element_id = cashback_de_id AND (pd.profile_id = @acquirer_id OR pd.profile_id = @chain_id)
        );

        /* Update the overridden column to be true for the site if the data element is present at the acquirer or chain level */
        IF @override_count > 0 THEN
            UPDATE profile_data
            SET overriden = 1
            WHERE profile_id = id AND data_element_id = cashback_de_id;
        END IF;
    END LOOP;

    CLOSE cur;
END;