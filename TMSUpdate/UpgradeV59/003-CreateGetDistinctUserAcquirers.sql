--multiline
CREATE PROCEDURE get_distinct_user_acquirers (
    IN userId INT
)
BEGIN
    SELECT DISTINCT acquirer_name
    FROM permissiongroup_acquirer pga
    LEFT JOIN user_permissiongroup upg ON upg.permission_group_id = pga.permissiongroup_id
    LEFT JOIN `user` u ON u.user_id = upg.user_id
    WHERE u.user_id = userId;
END