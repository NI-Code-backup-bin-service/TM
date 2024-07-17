--multiline;
CREATE PROCEDURE `acquirer_list_fetch`(
  IN search_term varchar(255),
  IN acquirers TEXT
)
BEGIN
  set @search = upper(concat('%', ifnull(search_term,''), '%'));
  select
    p.profile_id as 'country_profile_id',
    p.name
  from profile p
  where upper(p.name) like @search
    and FIND_IN_SET(p.name, acquirers)
    and p.profile_type_id = 2;
END