--multiline
CREATE PROCEDURE get_operations_group_permissions_by_name (
    IN groupName VARCHAR(255))
BEGIN
    SELECT p.name
    FROM operations_permission p LEFT JOIN operations_permissiongroup_permission pgp ON p.permission_id = pgp.permission_id
    LEFT JOIN operations_permissiongroup pg ON pgp.permissiongroup_id = pg.group_id
    WHERE pg.name = groupName;
END