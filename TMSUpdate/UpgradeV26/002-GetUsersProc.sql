--multiline
CREATE PROCEDURE `get_users`(IN acquirers TEXT, IN user_id int)
BEGIN
    set @user_role_id = (select u.roleId from user u where u.user_id = user_id);
    if @user_role_id != 3 then
        select distinct
            u.user_id,
            u.username
        from user u
                 left join user_permissiongroup upg on upg.user_id = u.user_id
                 left join permissiongroup pg on pg.group_id = upg.permission_group_id
                 left join permissiongroup_acquirer pga on pga.permissiongroup_id = pg.group_id
            and FIND_IN_SET(pga.acquirer_name, acquirers)
        where u.roleId != 3 and pga.permissiongroup_id is not null
        order by username asc;
    else
        select
            u.user_id,
            u.username
        from user u
        where u.roleId != 3
        order by username asc;
    end if;
END