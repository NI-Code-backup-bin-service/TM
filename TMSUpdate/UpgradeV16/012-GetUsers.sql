--multiline;
CREATE PROCEDURE `get_users`()
BEGIN
    select
        u.user_id,
        u.username
    from user u
    where u.roleId != 3
    order by username asc;
END