--multiline
CREATE PROCEDURE `user_store`(
    IN $username varchar(255),
    IN $roleId int,
    IN $updated_by varchar(255)
)
BEGIN

    insert into user(
        username,
        roleId,
        updated_at,
        updated_by,
        created_at,
        created_by
    ) values (
                 $username,
                 $roleId,
                 current_timestamp,
                 $updated_by,
                 current_timestamp,
                 $updated_by) ON DUPLICATE KEY UPDATE username=$username,roleId=$roleId,updated_at=current_timestamp,updated_by=$updated_by;

END