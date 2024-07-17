--multiline;
insert ignore into permissiongroup_acquirer (permissiongroup_id, acquirer_profile_id, acquirer_name)
select
  pg.group_id,
  p2.profile_id,
  p2.name
from permissiongroup pg
left join profile p2 on p2.profile_type_id = 2 and p2.name = "NI";