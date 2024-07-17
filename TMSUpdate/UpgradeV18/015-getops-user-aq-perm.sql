--multiline
CREATE PROCEDURE `get_operations_user_acquirer_permissions`(
  IN user_id int
)
begin
  set @user_role_id = (select u.roleId from operations_user u where u.user_id = user_id);
  if @user_role_id != 3 then
    select
      acquirer_name
    from operations_permissiongroup_acquirer pga
           left join operations_user_permissiongroup upg on upg.permission_group_id = pga.permissiongroup_id
           left join operations_user u on u.user_id = upg.user_id
    where u.user_id = user_id;
  else
    select
      p.name
    from profile p where p.profile_type_id = 2;
  end if;
end