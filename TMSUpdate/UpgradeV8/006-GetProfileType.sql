--multiline;
CREATE PROCEDURE `get_profile_type`(IN profileId int)
BEGIN
	SELECT
		pt.name
	from profile_type pt
    left join profile p on p.profile_type_id = pt.profile_type_id
    WHERE p.profile_id = profileId;
END