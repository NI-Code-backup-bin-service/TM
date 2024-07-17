--multiline;
CREATE PROCEDURE `delete_site_user`(IN userId int)
BEGIN
    -- First delete all TID users for the site with the same username
    DELETE tu
    FROM site_level_users su
    INNER JOIN tid_site ts
        ON su.site_id = ts.site_id
    INNER JOIN tid_user_override tu
        ON ts.tid_id = tu.tid_id
    WHERE su.user_id = userId AND tu.username = su.username;

    DELETE FROM site_level_users WHERE user_id = userId;
END