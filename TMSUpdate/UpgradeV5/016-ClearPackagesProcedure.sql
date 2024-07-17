--multiline;
create procedure clear_packages()
begin
    SET FOREIGN_KEY_CHECKS = 0;
    truncate table package_apk;
    truncate table package;
    alter table package auto_increment = 1;
    truncate table apk;
    alter table apk auto_increment = 1;
    SET FOREIGN_KEY_CHECKS = 1;
end;