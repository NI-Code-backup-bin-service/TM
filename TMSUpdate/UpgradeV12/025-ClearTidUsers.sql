--multiline;
CREATE PROCEDURE `clear_tid_users`(IN tidId INT)
BEGIN
	DELETE FROM tid_user_override WHERE tid_id = tidId;
END