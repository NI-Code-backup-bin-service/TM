--multiline
CREATE PROCEDURE get_operations_user_group_names (
    IN userId INT
)
BEGIN
    SELECT g.name
    FROM operations_permissiongroup g
             LEFT JOIN operations_user_permissiongroup upg ON upg.permission_group_id = g.group_id
    WHERE upg.user_id = userId;
END