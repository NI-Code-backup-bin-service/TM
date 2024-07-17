--multiline
CREATE VIEW `chain_data` AS
SELECT
    `profile_ids`.`chain_profile_id` AS `profile_id`,
    `de`.`data_element_id` AS `data_element_id`,
    COALESCE(`apd_chain`.`source`,
             `apd_acquirer`.`source`,
             `apd_global`.`source`) AS `source`,
    COALESCE(`apd_chain`.`datavalue`,
             `apd_acquirer`.`datavalue`,
             `apd_global`.`datavalue`) AS `datavalue`,
    COALESCE(`apd_chain`.`overriden`,
             `apd_acquirer`.`overriden`,
             `apd_global`.`overriden`) AS `overriden`,
    COALESCE(`apd_chain`.`is_encrypted`,
             `apd_acquirer`.`is_encrypted`,
             `apd_global`.`is_encrypted`) AS `is_encrypted`
FROM
    ((((((SELECT
              `cp`.`chain_profile_id` AS `chain_profile_id`,
              `cp`.`acquirer_id` AS `acquirer_profile_id`,
              1 AS `global_profile_id`
          FROM
              `chain_profiles` `cp`) `profile_ids`
        JOIN `data_element` `de`)
        JOIN `data_group` `dg` ON ((`dg`.`data_group_id` = `de`.`data_group_id`)))
        LEFT JOIN `approved_profile_data` `apd_chain` ON (((`apd_chain`.`profile_id` = `profile_ids`.`chain_profile_id`)
            AND (`apd_chain`.`data_element_id` = `de`.`data_element_id`))))
        LEFT JOIN `approved_profile_data` `apd_acquirer` ON (((`apd_acquirer`.`profile_id` = `profile_ids`.`acquirer_profile_id`)
            AND (`apd_acquirer`.`data_element_id` = `de`.`data_element_id`))))
        LEFT JOIN `approved_profile_data` `apd_global` ON (((`apd_global`.`profile_id` = `profile_ids`.`global_profile_id`)
        AND (`apd_global`.`data_element_id` = `de`.`data_element_id`))))