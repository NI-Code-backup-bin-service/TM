--multiline
CREATE PROCEDURE `get_operations_group_permissions`(IN groupId int)
BEGIN
  select p.permission_id
  from operations_permission p
  left join operations_permissiongroup_permission pgp on pgp.permission_id = p.permission_id
  where pgp.permissiongroup_id = groupId;
END