--multiline;
CREATE PROCEDURE `get_acquirer_name`(IN profileID INT)
BEGIN
    SELECT p.`name`
    FROM site_profiles sp
    LEFT JOIN `profile` p
        ON p.profile_id = sp.profile_id
    LEFT JOIN profile_type pt
        ON p.profile_type_id = pt.profile_type_id
    WHERE sp.site_id = (SELECT site_id FROM site_profiles WHERE profile_id = profileID limit 1) AND pt.name = 'acquirer';
END;