insert into permissiongroup (group_id, name, default_group) values (0, 'GlobalAdmin', 1);
insert into permissiongroup_permission (permissiongroup_id, permission_id) values (1, 11);
update permissiongroup set name = "NI Admin" where name = "Admin";
update permissiongroup set name = "NI User" where name = "User";