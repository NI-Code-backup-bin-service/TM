--multiline
create view site_data as
select `sp`.`site_id`         AS `site_id`,
       `pt`.`priority`        AS `priority`,
       `pt`.`name`            AS `level`,
       `pd`.`profile_data_id` AS `profile_data_id`,
       `pd`.`profile_id`      AS `profile_id`,
       `pd`.`data_element_id` AS `data_element_id`,
       `pd`.`datavalue`       AS `datavalue`,
       `pd`.`version`         AS `version`,
       `pd`.`updated_at`      AS `updated_at`,
       `pd`.`updated_by`      AS `updated_by`,
       `pd`.`created_at`      AS `created_at`,
       `pd`.`created_by`      AS `created_by`,
       `pd`.`overriden`       AS `overriden`,
       `pd`.`is_encrypted`    AS `is_encrypted`,
       `pd`.`not_overridable`    AS `not_overridable`
from (((`NextGen_TMS`.`profile_data` `pd` join `NextGen_TMS`.`profile` `p` on ((`p`.`profile_id` = `pd`.`profile_id`))) join `NextGen_TMS`.`profile_type` `pt` on ((`pt`.`profile_type_id` = `p`.`profile_type_id`)))
         join `NextGen_TMS`.`site_profiles` `sp` on ((`sp`.`profile_id` = `pd`.`profile_id`)))
where (`pd`.`approved` = 1);
