INSERT IGNORE INTO `permission` (`permission_id`, `name`) VALUES ('22', 'Db Backup');
INSERT IGNORE INTO `permissiongroup_permission` (`permissiongroup_id`, `permission_id`) VALUES ((SELECT `group_id` FROM `permissiongroup` WHERE `name` = 'GlobalAdmin'), (SELECT `permission_id` FROM `permission` WHERE `name` = 'Db Backup'));