ALTER TABLE tid DROP COLUMN Presence;
DELETE FROM permissiongroup_permission WHERE permission_id = (SELECT permission.permission_id FROM permission WHERE name = 'TID Logs');
DELETE FROM permission WHERE name = 'TID Logs';