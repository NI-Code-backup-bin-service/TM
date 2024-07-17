--multiline;
CREATE PROCEDURE `operations_user_fetch`(
  IN username varchar(255)
)
BEGIN
  select user_id, username, roleId
  from operations_user ou
  where ou.username = username;
END