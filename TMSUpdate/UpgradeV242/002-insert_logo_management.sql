INSERT IGNORE INTO `permission` (`permission_id`,`name`) VALUES ('29','Logo Management');
INSERT IGNORE INTO `permissiongroup_permission` (`permissiongroup_id`, `permission_id`) VALUES ((SELECT `group_id` FROM `permissiongroup` WHERE `name` = 'GlobalAdmin'), (SELECT `permission_id` FROM `permission` WHERE `name` = 'Logo Management'));