--multiline
CREATE PROCEDURE get_user_groups(IN userId int)
BEGIN
  select g.group_id
  from permissiongroup g
  left join user_permissiongroup upg on upg.permission_group_id = g.group_id
  where upg.user_id = userId;
END