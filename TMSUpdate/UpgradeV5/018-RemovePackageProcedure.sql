--multiline;
create procedure remove_package(IN version varchar(45))
begin
START TRANSACTION;
    set @package_id = (select package_id from package p
                        where p.version = version);

    delete from package_apk where package_id = @package_id;
    delete from package where package_id = @package_id and @package_id is not null;   

    call cleanup_apks();
commit;
end;