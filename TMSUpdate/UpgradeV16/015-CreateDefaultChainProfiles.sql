--multiline;
insert ignore into chain_profiles (chain_profile_id, acquirer_id)
select
  p.profile_id,
  p2.profile_id
from profile p
left join profile p2 on p2.profile_type_id = 2 and p2.name = "NI"
where p.profile_type_id = 3;