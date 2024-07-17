--multiline;
CREATE PROCEDURE `get_site_profile_id_from_tid` (tid INT)
BEGIN
	SELECT p.profile_id
    FROM site_profiles AS sp
    LEFT JOIN profile AS p ON sp.profile_id = p.profile_id
    WHERE site_id = (SELECT site_id FROM tid_site WHERE tid_id = tid)
    AND p.profile_type_id = (SELECT profile_type_id FROM profile_type WHERE name = "site");
END
