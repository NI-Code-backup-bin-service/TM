--multiline;
CREATE PROCEDURE `get_site_profiles`(siteId INT)
BEGIN
	SELECT 
    p.profile_id,
    pt.profile_type_id,
    pt.`name`
    FROM site_profiles sp
    LEFT JOIN `profile` p ON p.profile_id = sp.profile_id
    LEFT JOIN profile_type pt ON p.profile_type_id = pt.profile_type_id;
END