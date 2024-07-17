--multiline
CREATE PROCEDURE `remove_operations_permission_group`(IN name varchar(45))
begin
START TRANSACTION;
    set @group_id = (select group_id from operations_permissiongroup p
                        where p.name = name);

    delete from operations_user_permissiongroup where permission_group_id = @group_id and @group_id is not null;
    delete from operations_permissiongroup_permission where permissiongroup_id = @group_id and @group_id is not null;
    delete from operations_permissiongroup where group_id = @group_id and @group_id is not null;

commit;
end