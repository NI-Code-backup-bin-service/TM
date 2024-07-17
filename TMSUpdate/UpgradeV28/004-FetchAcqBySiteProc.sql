--multiline;
CREATE PROCEDURE `fetch_site_acquirer`(IN siteId INT)
BEGIN
    SELECT p.`name`
    FROM site_profiles sp
             LEFT JOIN `profile` p ON p.profile_id = sp.profile_id
             LEFT JOIN profile_type pt ON p.profile_type_id = pt.profile_type_id
    WHERE sp.site_id = siteId AND pt.name = 'acquirer';
END