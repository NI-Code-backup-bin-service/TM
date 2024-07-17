--multiline
CREATE PROCEDURE chain_list_fetch(
 IN search_term varchar(255)
)
BEGIN
set @search = upper(concat('%', ifnull(search_term,''), '%'));
select
  p.profile_id as 'chain_profile_id',
  p.name
from profile p
where upper(p.name) like @search
and p.profile_type_id = 3;
END