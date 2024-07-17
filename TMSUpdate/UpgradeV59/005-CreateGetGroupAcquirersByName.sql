--multiline
CREATE PROCEDURE get_group_acquirers_by_name (
    IN groupName VARCHAR(255)
)
BEGIN
   SELECT pga.acquirer_name
       FROM permissiongroup_acquirer pga
    LEFT JOIN permissiongroup pg ON pg.group_id = pga.permissiongroup_id
    WHERE pg.name = groupName;
END