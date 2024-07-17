--multiline
CREATE VIEW `approved_profile_data` AS
SELECT
    `p`.`profile_id` AS `profile_id`,
    `pd`.`data_element_id` AS `data_element_id`,
    `pt`.`name` AS `source`,
    `pd`.`datavalue` AS `datavalue`,
    `pd`.`overriden` AS `overriden`,
    `pd`.`is_encrypted` AS `is_encrypted`
FROM
    ((`profile_data` `pd`
        JOIN `profile` `p` ON ((`p`.`profile_id` = `pd`.`profile_id`)))
        JOIN `profile_type` `pt` ON ((`p`.`profile_type_id` = `pt`.`profile_type_id`)))
WHERE
    (`pd`.`approved` = 1)