--multiline
CREATE PROCEDURE get_distinct_operations_user_acquirers (
    IN userId INT
)
BEGIN
    SELECT DISTINCT acquirer_name
    FROM operations_permissiongroup_acquirer pga
             LEFT JOIN operations_user_permissiongroup upg ON upg.permission_group_id = pga.permissiongroup_id
             LEFT JOIN operations_user u ON u.user_id = upg.user_id
    WHERE u.user_id = userId;
END