--multiline
create procedure `bulkUpdate`(
	IN versionId int,
	IN versionDataElement int
)
begin
	REPLACE INTO profile_data 
	(profile_id,data_element_id,datavalue,version,updated_at,updated_by,created_at,created_by,approved,overriden)
	(SELECT 
	p.profile_id,
	versionDataElement,
	307,
	1,
	NOW(),
	'BulkUpdate',
	NOW(),
	'BulkUpdate',
	1,
	1
	FROM site s
	LEFT JOIN site_profiles sp on s.site_id = sp.site_id
	LEFT JOIN `profile` p on  p.profile_id = sp.profile_id
	LEFT JOIN profile_type pt on pt.profile_type_id = p.profile_type_id
	WHERE pt.name = 'Site'
	AND IFNULL((SELECT pd.datavalue FROM profile_data pd WHERE pd.profile_id = p.profile_id AND pd.data_element_id = versionDataElement),0) < 307)
	LIMIT 100;

	INSERT INTO tid_updates (tid_update_id, tid_id, target_package_id, update_date)
	SELECT 
	(SELECT max(tid_update_id) from tid_updates tu WHERE tu.tid_id = t.tid_id)+1,
	t.tid_id, versionId, CURDATE()
	FROM tid t
	WHERE (SELECT p.version FROM tid_updates u 
			LEFT JOIN package p ON p.package_id = u.target_package_id 
			WHERE u.tid_id = t.tid_id 
			ORDER BY u.update_date DESC LIMIT 1) < 307
			
	AND ( SELECT 
	pd.datavalue 
	FROM tid_site ts
	LEFT JOIN site_profiles sp on ts.site_id = sp.site_id
	LEFT JOIN `profile` p on  p.profile_id = sp.profile_id
	LEFT JOIN profile_type pt on pt.profile_type_id = p.profile_type_id
	LEFT JOIN profile_data pd ON pd.profile_id = p.profile_id
	WHERE pd.data_element_id = versionDataElement AND ts.tid_id = t.tid_id AND pt.name = 'Site') = 307;
end