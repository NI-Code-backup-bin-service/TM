INSERT IGNORE INTO permission (name) VALUES ('Payment Services');
INSERT IGNORE INTO permissiongroup_permission (`permissiongroup_id`, `permission_id`) VALUES ((select group_id from permissiongroup where `name` = 'GlobalAdmin' LIMIT 1), (select permission_id from permission where `name` = 'Payment Services' LIMIT 1));