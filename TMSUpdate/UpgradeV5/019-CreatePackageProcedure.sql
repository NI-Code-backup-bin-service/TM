--multiline;
create procedure create_package(IN version varchar(45))
begin
    insert ignore into package(version) values (version);
end;
