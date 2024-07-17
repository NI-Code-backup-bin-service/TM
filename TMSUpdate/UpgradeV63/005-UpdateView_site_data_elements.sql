--multiline
create view site_data_elements as
select `sp`.`site_id`               AS `site_id`,
       `sp`.`profile_id`            AS `profile_id`,
       `de`.`data_element_id`       AS `data_element_id`,
       `de`.`data_group_id`         AS `data_group_id`,
       `de`.`name`                  AS `name`,
       `de`.`datatype`              AS `datatype`,
       `de`.`is_allow_empty`        AS `is_allow_empty`,
       `de`.`version`               AS `version`,
       `de`.`updated_at`            AS `updated_at`,
       `de`.`updated_by`            AS `updated_by`,
       `de`.`created_at`            AS `created_at`,
       `de`.`created_by`            AS `created_by`,
       `de`.`max_length`            AS `max_length`,
       `de`.`validation_expression` AS `validation_expression`,
       `de`.`validation_message`    AS `validation_message`,
       `de`.`front_end_validate`    AS `front_end_validate`,
       `de`.`options`               AS `options`,
       `de`.`displayname_en`        AS `display_name`,
       `de`.`is_password`           AS `is_password`,
       `de`.`is_encrypted`          AS `is_encrypted`,
       `de`.`sort_order_in_group`   AS `sort_order_in_group`,
       `del`.`location_name`        AS `location_name`
FROM data_element de
         INNER JOIN profile_data_group pg ON
        pg.data_group_id = de.data_group_id
         INNER JOIN site_profiles sp ON
        sp.profile_id = pg.profile_id
         LEFT OUTER JOIN data_element_locations_data_element delde ON
        delde.data_element_id = de.data_element_id
         INNER JOIN data_element_locations del ON
        delde.location_id = del.location_id;