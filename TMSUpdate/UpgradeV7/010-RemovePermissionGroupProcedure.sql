--multiline;
create procedure remove_permission_group(IN name varchar(45))
begin
START TRANSACTION;
    set @group_id = (select group_id from permissiongroup p
                        where p.name = name);

    delete from user_permissiongroup where permission_group_id = @group_id and @group_id is not null;
    delete from permissiongroup_permission where permissiongroup_id = @group_id and @group_id is not null;
    delete from permissiongroup where group_id = @group_id and @group_id is not null;

commit;
end;