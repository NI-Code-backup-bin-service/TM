--multiline
CREATE PROCEDURE `tid_list_fetch`(IN search_term varchar(255))
BEGIN
  set @search = upper(concat('%', ifnull(search_term,''), '%'));
  select
    t.tid_id,
    t.serial,
    t.PIN,
    t.ExpiryDate,
    t.ActivationDate,
    s.site_id,
    s.name
  from tid t
  left join tid_site ts on ts.tid_id = t.tid_id
  left join site s on s.site_id = ts.site_id
  where ts.tid_id = t.tid_id
  and (upper(t.tid_id) like @search
  or upper(t.serial) like @search)
  or upper(s.name) like @search;
END