--multiline;
CREATE PROCEDURE `get_tid_user_data`(IN tidId int)
BEGIN
	SELECT * FROM tid_user_override tuo WHERE tuo.tid_id = tidId ORDER BY tuo.Username;
END