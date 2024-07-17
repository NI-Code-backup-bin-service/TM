INSERT INTO permission (`permission_id`, `name`) VALUES (13, 'Edit Passwords');
INSERT INTO permissiongroup_permission (permissiongroup_id, permission_id) VALUES (10, 1);
INSERT INTO permissiongroup_permission (permissiongroup_id, permission_id) VALUES (10, 13);