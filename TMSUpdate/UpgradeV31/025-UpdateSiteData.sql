--multiline
CREATE
    VIEW `site_data` AS
SELECT
    `sp`.`site_id` AS `site_id`,
    `pt`.`priority` AS `priority`,
    `pt`.`name` AS `level`,
    `pd`.`profile_data_id` AS `profile_data_id`,
    `pd`.`profile_id` AS `profile_id`,
    `pd`.`data_element_id` AS `data_element_id`,
    `pd`.`datavalue` AS `datavalue`,
    `pd`.`version` AS `version`,
    `pd`.`updated_at` AS `updated_at`,
    `pd`.`updated_by` AS `updated_by`,
    `pd`.`created_at` AS `created_at`,
    `pd`.`created_by` AS `created_by`,
    `pd`.`overriden` AS `overriden`,
    `pd`.`is_encrypted` AS `is_encrypted`
FROM
    (((`profile_data` `pd`
        JOIN `profile` `p` ON ((`p`.`profile_id` = `pd`.`profile_id`)))
        JOIN `profile_type` `pt` ON ((`pt`.`profile_type_id` = `p`.`profile_type_id`)))
        JOIN `site_profiles` `sp` ON ((`sp`.`profile_id` = `pd`.`profile_id`)))
WHERE
    (`pd`.`approved` = 1)