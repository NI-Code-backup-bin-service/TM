--multiline;
CREATE PROCEDURE `delete_site_user`(IN userId int)
BEGIN
	DELETE FROM site_level_users WHERE user_id = userId;
END