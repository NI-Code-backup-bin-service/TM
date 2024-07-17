--multiline
CREATE PROCEDURE `user_fetch`(
  IN $token varchar(255)
)
BEGIN

  select user_id, username, token, roleId
  from user
  where token = $token;

END