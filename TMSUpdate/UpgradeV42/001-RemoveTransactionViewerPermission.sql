DELETE FROM permission WHERE `name` = 'Transaction Viewer';
DELETE FROM permissiongroup_permission WHERE permission_id not in (select permission_id from permission);