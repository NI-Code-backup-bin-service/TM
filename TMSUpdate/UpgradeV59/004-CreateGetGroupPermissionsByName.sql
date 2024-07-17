--multiline
CREATE PROCEDURE get_group_permissions_by_name (
    IN groupName VARCHAR(255))
BEGIN
    SELECT p.name
    FROM permission p LEFT JOIN permissiongroup_permission pgp ON p.permission_id = pgp.permission_id
    LEFT JOIN permissiongroup pg ON pgp.permissiongroup_id = pg.group_id
    WHERE pg.name = groupName;
END