--multiline
create table if not exists permission (
permission_id int(11) NOT NULL AUTO_INCREMENT,
name varchar(255) not null unique,
primary key(permission_id)
);