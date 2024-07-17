--multiline
CREATE PROCEDURE get_operations_group_acquirers_by_name (
    IN groupName VARCHAR(255)
)
BEGIN
    SELECT pga.acquirer_name
    FROM operations_permissiongroup_acquirer pga
    LEFT JOIN operations_permissiongroup pg ON pg.group_id = pga.permissiongroup_id
    WHERE pg.name = groupName;
END