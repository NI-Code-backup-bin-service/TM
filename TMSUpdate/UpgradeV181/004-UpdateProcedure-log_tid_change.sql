--multiline;
CREATE PROCEDURE `log_tid_change`(IN P_tid int, IN P_data_element_id int, IN P_change_type int, IN P_current_value text, IN P_updated_value text, IN P_updated_by varchar(255), IN P_approved int)
BEGIN
    SELECT tsp.site_id, tsp.profile_id
    INTO @site_id, @profile_id
    FROM
        (SELECT

             `data`.`profile_id` AS `profile_id`,

             `data`.`tid_id` AS `tid_id`,

             `data`.`site_id` AS `site_id`

         FROM

             (SELECT

                  `ts`.`tid_profile_id` AS `profile_id`,

                  `ts`.`tid_id` AS `tid_id`,

                  `ts`.`site_id` AS `site_id`

              FROM

                  (`tid_site` `ts`

                      LEFT JOIN `site_profiles` `sp` ON ((`ts`.`site_id` = `sp`.`site_id`)))

              WHERE

                  (`ts`.`tid_id` = P_tid) UNION SELECT

                                                    `sp`.`profile_id` AS `profile_id`,

                                                    `ts`.`tid_id` AS `tid_id`,

                                                    `sp`.`site_id` AS `site_id`

              FROM

                  (`tid_site` `ts`

                      LEFT JOIN `site_profiles` `sp` ON ((`ts`.`site_id` = `sp`.`site_id`)))

              WHERE

                  (`ts`.`tid_id` = P_tid)) `data`) tsp
            INNER JOIN profile p on
                tsp.profile_id = p.profile_id
            INNER JOIN profile_type pt ON
                p.profile_type_id = pt.profile_type_id
    WHERE tid_id = P_tid
    ORDER BY pt.priority asc
    LIMIT 1;

    SELECT p.name
    INTO @acquires
    FROM profile p
             LEFT JOIN site_profiles tp4 on tp4.profile_id = p.profile_id
    WHERE tp4.site_id = @site_id AND p.profile_type_id = 2;

    INSERT INTO approvals(profile_id,
                          data_element_id,
                          change_type,
                          current_value,
                          new_value,
                          created_at,
                          approved_at,
                          created_by,
                          approved_by,
                          approved,
                          tid_id,
                          acquirer)
    VALUES (@profile_id,
            P_data_element_id,
            P_change_type,
            P_current_value,
            P_updated_value,
            NOW(),
            NOW(),
            P_updated_by,
            P_updated_by,
            P_approved,
            P_tid,
            @acquires);
END