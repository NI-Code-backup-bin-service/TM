--multiline
CREATE PROCEDURE `add_site_user`(IN userId INT, IN siteId INT, IN Username VARCHAR(255), IN Pin VARCHAR(80), IN Modules VARCHAR(255), IN is_encrypted BOOL )
BEGIN
    IF userId <= 0 THEN
        BEGIN
            INSERT INTO site_level_users(site_id, Username,PIN,Modules,is_encrypted) VALUES (siteId,Username,Pin,Modules,is_encrypted);
        END;
    ELSE
        BEGIN
            UPDATE site_level_users SET site_id = site_id, Username = Username, PIN = Pin, Modules = Modules, is_encrypted = is_encrypted WHERE user_id = userId;
        END;
    END IF;
END