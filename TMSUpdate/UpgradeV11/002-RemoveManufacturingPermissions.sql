delete pp from permissiongroup_permission pp left join permission p on pp.permission_id = p.permission_id where p.name = 'Manufacturing';
delete from permission where name = 'Manufacturing';
