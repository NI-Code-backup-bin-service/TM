--multiline
create procedure add_site_user(IN userId int, IN siteId int, IN Username varchar(255), IN Pin varchar(66), IN Modules LONGTEXT, IN Encrypted tinyint(1))
BEGIN
    IF userId <= 0 THEN
        INSERT INTO site_level_users(site_id, Username, PIN, Modules, is_encrypted) VALUES (siteId, Username, Pin, Modules, Encrypted);
    ELSE
        UPDATE site_level_users SET site_id = site_id, Username = Username, PIN = Pin, Modules = Modules, is_encrypted = Encrypted WHERE user_id = userId;
    END IF;

    -- Update the PINs of any TID overrides belonging to the site to match
    UPDATE site_level_users su
        INNER JOIN tid_site ts
        ON su.site_id = ts.site_id
        INNER JOIN tid_user_override tu
        ON ts.tid_id = tu.tid_id
    SET tu.PIN = su.PIN
    WHERE su.user_id = userId AND tu.username = su.username;

    /*Return the user ID of the just inserted user*/
    SELECT user_id FROM site_level_users WHERE username = Username AND site_id = siteId;
END;
