--multiline;
CREATE PROCEDURE `get_tid_overridable_active_data_group_ids`(IN siteId INT)
BEGIN
SELECT distinct dg.data_group_id
	FROM site_profiles sp
	INNER JOIN site_profiles sp2 ON sp.site_id = sp2.site_id
	INNER JOIN profile_data_group pdg ON pdg.profile_id=sp2.profile_id
	INNER JOIN data_group dg ON pdg.data_group_id=dg.data_group_id
    INNER JOIN data_element de ON dg.data_group_id=de.data_group_id
	WHERE sp.profile_id = (SELECT sp.profile_Id FROM site_profiles sp 
		LEFT JOIN profile p ON p.profile_id = sp.profile_id 
		LEFT JOIN profile_type pt ON pt.profile_type_id = p.profile_type_id
        WHERE sp.site_id = siteId  AND pt.priority = 2)
	AND sp2.profile_id != 1 AND de.tid_overridable=1; 
END