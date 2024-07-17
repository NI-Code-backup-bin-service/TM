--multiline
CREATE PROCEDURE `get_operations_users`()
BEGIN
    select
        u.user_id,
        u.username
    from operations_user u
    where u.roleId != 3
    order by username asc;
END