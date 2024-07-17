--multiline
CREATE PROCEDURE country_list_fetch(
 IN search_term varchar(255)
)
BEGIN
set @search = upper(concat('%', ifnull(search_term,''), '%'));
select
  p.profile_id as 'country_profile_id',
  p.name
from profile p
where upper(p.name) like @search
and p.profile_type_id = 2;
END