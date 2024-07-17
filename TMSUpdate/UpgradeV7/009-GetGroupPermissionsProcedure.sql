--multiline
CREATE PROCEDURE get_group_permissions(IN groupId int)
BEGIN
  select p.permission_id
  from permission p
  left join permissiongroup_permission pgp on pgp.permission_id = p.permission_id
  where pgp.permissiongroup_id = groupId;
END