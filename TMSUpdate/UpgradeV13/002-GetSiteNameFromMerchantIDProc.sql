--multiline;
CREATE PROCEDURE `Get_site_name_from_merchantID`(IN merchantId varchar(255))
BEGIN
  select
    pdn.datavalue as 'site_name'
  from site t
         LEFT JOIN (site_profiles tp
    join profile p on p.profile_id = tp.profile_id
    join profile_type pt on pt.profile_type_id = p.profile_type_id and pt.priority = 2)
                   on tp.site_id = t.site_id
    -- Get the merchant number
         LEFT JOIN profile_data pd ON pd.profile_id = p.profile_id AND pd.data_element_id = 1
    AND pd.version = (SELECT MAX(d.version) FROM profile_data d WHERE d.data_element_id = 1 AND d.profile_id = p.profile_id AND d.approved = 1)
    -- Get the merchant name
         LEFT JOIN profile_data pdn on pdn.profile_id = p.profile_id and pdn.data_element_id = 3
    AND pdn.version = (SELECT MAX(d.version) FROM profile_data d WHERE d.data_element_id = 3 AND d.profile_id = p.profile_id AND d.approved = 1 )
  where pd.datavalue = merchantId;
END