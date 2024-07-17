--multiline;
CREATE PROCEDURE `add_site_user`(IN userId INT, IN siteId INT, IN Username VARCHAR(255), IN Pin VARCHAR(5), IN Modules VARCHAR(255) )
BEGIN
	IF userId <= 0 THEN
    BEGIN
		INSERT INTO site_level_users(site_id, Username,PIN,Modules) VALUES (siteId,Username,Pin,Modules);
	END;
    ELSE
    BEGIN
		UPDATE site_level_users SET site_id = site_id, Username = Username, PIN = Pin, Modules = Modules WHERE user_id = userId;
	END;
    END IF;
END