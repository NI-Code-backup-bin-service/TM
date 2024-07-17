--multiline;
CREATE PROCEDURE `tid_list_fetch`(IN search_term varchar(255))
BEGIN
  set @search = upper(concat('%', ifnull(search_term,''), '%'));
  select
    t.tid_id,
    t.serial,
    t.PIN,
    t.ExpiryDate,
    t.ActivationDate,
    t.Presence,
    s.site_id,
    pdn.datavalue as 'name',
    pd.datavalue as 'merchant_id'
  from tid t
         left join tid_site ts on ts.tid_id = t.tid_id
         left join site s on s.site_id = ts.site_id

         LEFT JOIN (site_profiles tp2
    join profile p2 on p2.profile_id = tp2.profile_id
    join profile_type pt2 on pt2.profile_type_id = p2.profile_type_id and pt2.priority = 2) on tp2.site_id = s.site_id

         LEFT JOIN profile_data pd ON pd.profile_id = p2.profile_id AND pd.data_element_id = 1
    AND pd.version = (SELECT MAX(d.version) FROM profile_data d WHERE d.data_element_id = 1 AND d.profile_id = p2.profile_id AND d.approved = 1)

         LEFT JOIN profile_data pdn ON pdn.profile_id = p2.profile_id AND pdn.data_element_id = 3
    AND pdn.version = (SELECT MAX(d.version) FROM profile_data d WHERE d.data_element_id = 3 AND d.profile_id = p2.profile_id AND d.approved = 1)

  where ts.tid_id = t.tid_id
    and (upper(t.tid_id) like @search
      or upper(t.serial) like @search)
     or upper(pdn.datavalue) like @search;
END