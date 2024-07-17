--multiline;
CREATE PROCEDURE `add_apk_to_package`(IN version varchar(45), IN apk_name varchar(45))
BEGIN
    set @package_id = (select package_id from package p
                       where p.version = version);

    set @apk_id = (select apk_id from apk a
                   where a.name = apk_name);
                   
	if @apkId IS NULL then
		      insert into apk(name) values (apk_name);
			  set @apk_id = (select apk_id from apk a where a.name = apk_name);
	end if;

    if (@package_id is not null and @apk_id is not null) then
      insert ignore into package_apk values (@package_id, @apk_id);
    end if;
  END