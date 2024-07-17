--multiline
CREATE PROCEDURE `add_enoc_site_id_to_sites`()
BEGIN
    DECLARE tr_id INT DEFAULT (SELECT data_group_id FROM data_group WHERE name='transactionRetrieval' LIMIT 1);
    DECLARE siteid_id INT DEFAULT (SELECT data_element_id FROM data_element WHERE name="siteID" AND data_group_id=tr_id LIMIT 1);
    DECLARE done INT DEFAULT FALSE;
    DECLARE id INT;

    DECLARE cur CURSOR FOR SELECT DISTINCT p.profile_id
                           FROM profile p
                                    INNER JOIN profile_data pd ON p.profile_id = pd.profile_id
                                    INNER JOIN data_element de ON pd.data_element_id = de.data_element_id
                                    INNER JOIN data_group dg ON de.data_group_id = dg.data_group_id
                                    INNER JOIN profile_type pt ON p.profile_type_id = pt.profile_type_id
                           WHERE pt.name = 'site';
    DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;
    OPEN cur;

    read_loop: LOOP
        FETCH cur INTO id;

        IF done THEN
            LEAVE read_loop;
        END IF;

        INSERT IGNORE INTO profile_data(profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden, is_encrypted, not_overridable)
        VALUES (id, siteid_id, "", 1, NOW(), 'system', NOW(), 'system', 1, 0, 0, 0);

    END LOOP;


END