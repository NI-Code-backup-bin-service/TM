INSERT IGNORE INTO operations_permissiongroup_permission (permissiongroup_id, permission_id) VALUES ((select group_id from operations_permissiongroup pg where pg.name = 'GlobalAdmin'), 12);
INSERT IGNORE INTO operations_permissiongroup_permission (permissiongroup_id, permission_id) VALUES ((select group_id from operations_permissiongroup pg where pg.name = 'GlobalAdmin'), 14);
INSERT IGNORE INTO operations_permissiongroup_permission (permissiongroup_id, permission_id) VALUES ((select group_id from operations_permissiongroup pg where pg.name = 'GlobalAdmin'), 16);