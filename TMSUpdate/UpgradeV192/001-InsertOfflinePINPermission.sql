INSERT INTO `permission` (`permission_id`, `name`) VALUES ('21', 'Offline PIN');
INSERT INTO `permissiongroup_permission` (`permissiongroup_id`, `permission_id`) VALUES ((SELECT `group_id` FROM `permissiongroup` WHERE `name` = 'GlobalAdmin'), (SELECT `permission_id` FROM `permission` WHERE `name` = 'Offline PIN'));