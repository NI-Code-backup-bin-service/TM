INSERT IGNORE INTO permission (permission_id, name) VALUES (16, 'User Management Audit');
INSERT IGNORE INTO permissiongroup_permission (permissiongroup_id, permission_id) VALUES ((select group_id from permissiongroup pg where pg.name = 'GlobalAdmin'), 16);
INSERT IGNORE INTO operations_permission (permission_id, name) VALUES (9, 'User Management Audit');
INSERT IGNORE INTO operations_permissiongroup_permission (permissiongroup_id, permission_id) VALUES ((select group_id from operations_permissiongroup pg where pg.name = 'GlobalAdmin'), 9);