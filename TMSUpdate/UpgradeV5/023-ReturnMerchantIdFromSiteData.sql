--multiline;
CREATE PROCEDURE `site_list_fetch`(
 IN search_term varchar(255)
)
BEGIN
set @search = upper(concat('%', ifnull(search_term,''), '%'));
select *
from (
	select 
    t.site_id as 'site_profile_id',
	p2.name as 'site_name',
	p3.profile_id as 'chain_profile_id',
	p3.name 'chain_name',
	p4.profile_id as 'country_profile_id',
	p4.name as 'country_name',
	p5.profile_id as 'global_profile_id',
	p5.name as 'global_name',
	pd.datavalue as 'merchant_id'
	from site t
	LEFT JOIN 
	(site_profiles tp1 
	join profile p1 on p1.profile_id = tp1.profile_id 
	join profile_type pt1 on pt1.profile_type_id = p1.profile_type_id and pt1.priority = 1) 
	on tp1.site_id = t.site_id 
    LEFT JOIN
	(site_profiles tp2 
	join profile p2 on p2.profile_id = tp2.profile_id 
	join profile_type pt2 on pt2.profile_type_id = p2.profile_type_id and pt2.priority = 2) 
	on tp2.site_id = t.site_id 
    -- Get the merchant number
    LEFT JOIN profile_data pd ON pd.profile_id = p2.profile_id AND pd.data_element_id = 1
		AND pd.version = (SELECT MAX(d.version) FROM profile_data d WHERE d.data_element_id = 1 AND d.profile_id = p2.profile_id)
	LEFT JOIN 
	(site_profiles tp3 
	join profile p3 on p3.profile_id = tp3.profile_id 
	join profile_type pt3 on pt3.profile_type_id = p3.profile_type_id and pt3.priority = 3) 
	on tp3.site_id = t.site_id 
	LEFT JOIN
	(site_profiles tp4
	join profile p4 on p4.profile_id = tp4.profile_id 
	join profile_type pt4 on pt4.profile_type_id = p4.profile_type_id and pt4.priority = 4) 
	on tp4.site_id = t.site_id 
	LEFT JOIN
	(site_profiles tp5
	join profile p5 on p5.profile_id = tp5.profile_id 
	join profile_type pt5 on pt5.profile_type_id = p5.profile_type_id and pt5.priority = 5) 
	on tp5.site_id = t.site_id 
) l
where l.site_name like @search
OR upper(l.site_name) like @search
OR upper(l.chain_name) like @search
OR upper(l.country_name) like @search
OR upper(l.global_name) like @search
OR l.merchant_id like @search;
END