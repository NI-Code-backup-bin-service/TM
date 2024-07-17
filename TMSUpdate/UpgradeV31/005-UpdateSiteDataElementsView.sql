--multiline
CREATE
    VIEW `site_data_elements` AS
SELECT
    `sp`.`site_id` AS `site_id`,
    `sp`.`profile_id` AS `profile_id`,
    `de`.`data_element_id` AS `data_element_id`,
    `de`.`data_group_id` AS `data_group_id`,
    `de`.`name` AS `name`,
    `de`.`datatype` AS `datatype`,
    `de`.`is_allow_empty` AS `is_allow_empty`,
    `de`.`version` AS `version`,
    `de`.`updated_at` AS `updated_at`,
    `de`.`updated_by` AS `updated_by`,
    `de`.`created_at` AS `created_at`,
    `de`.`created_by` AS `created_by`,
    `de`.`max_length` AS `max_length`,
    `de`.`validation_expression` AS `validation_expression`,
    `de`.`validation_message` AS `validation_message`,
    `de`.`front_end_validate` AS `front_end_validate`,
    `de`.`options` AS `options`,
    `de`.`displayname_en` AS `display_name`,
    `de`.`is_password` AS `is_password`,
    `de`.`is_encrypted` AS `is_encrypted`
FROM
    ((`data_element` `de`
        JOIN `profile_data_group` `pg` ON ((`pg`.`data_group_id` = `de`.`data_group_id`)))
        JOIN `site_profiles` `sp` ON ((`sp`.`profile_id` = `pg`.`profile_id`)));
