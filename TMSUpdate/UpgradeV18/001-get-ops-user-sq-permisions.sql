--multiline
CREATE PROCEDURE `get_operations_user_groups`(IN userId int)
BEGIN
  select g.group_id
  from operations_permissiongroup g
  left join operations_user_permissiongroup upg on upg.permission_group_id = g.group_id
  where upg.user_id = userId;
END