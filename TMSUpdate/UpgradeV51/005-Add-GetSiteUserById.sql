--multiline
CREATE PROCEDURE `get_site_user_by_id`(IN p_userId int)
BEGIN
    SELECT * FROM site_level_users WHERE user_id = p_userId;
END