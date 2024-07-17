--multiline;
CREATE PROCEDURE `get_site_user_data`(IN siteId INT)
BEGIN
 SELECT * FROM site_level_users slu WHERE slu.site_id = siteId ORDER BY slu.username;
END