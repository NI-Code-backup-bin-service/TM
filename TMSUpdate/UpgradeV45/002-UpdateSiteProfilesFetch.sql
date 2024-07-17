--multiline;
CREATE PROCEDURE `site_profiles_fetch`(IN site int)
BEGIN
	SELECT sp.profile_id
    FROM site_profiles as sp
    left join profile as p ON p.profile_id = sp.profile_id
    left join profile_type as pt ON pt.profile_type_id = p. profile_type_id
    WHERE site_id = site
    order by pt.priority
    limit 1;
END