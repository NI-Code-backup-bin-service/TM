--multiline;
CREATE PROCEDURE `get_permission_groups`(
    IN acquirers TEXT,
    IN user_id int
)
BEGIN
    set @user_role_id = (select u.roleId from user u where u.user_id = user_id);
    if @user_role_id != 3 then
        select
            pg.group_id,
            pg.name,
            pg.default_group,
            case when pga.permissiongroup_id is null then FALSE else TRUE end as 'userAccess'
        from permissiongroup pg
                 left join permissiongroup_acquirer pga on pga.permissiongroup_id = pg.group_id
            and FIND_IN_SET(pga.acquirer_name, acquirers)
        where name != "GlobalAdmin"
        order by group_id asc;
    else
        select *, TRUE from permissiongroup
        where name != "GlobalAdmin"
        order by group_id asc;
    end if;
END