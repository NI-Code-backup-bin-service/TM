--multiline;
CREATE PROCEDURE add_apk_to_package (IN version varchar(45), IN apk_name varchar(45))
BEGIN
set @package_id = (select package_id from package p 
                    where p.version = version);

if (@package_id is not null) then
 insert into apk(name) values (apk_name);
 insert into package_apk values (@package_id, last_insert_id());
end if;
END