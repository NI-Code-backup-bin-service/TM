--multiline
create table if not exists permissiongroup (
group_id int(11) NOT NULL AUTO_INCREMENT,
name varchar(255) not null unique,
default_group int(11) DEFAULT 0,
primary key(group_id)
);