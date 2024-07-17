--multiline
CREATE PROCEDURE `delete_tid_user`(IN p_tidUserId int)
BEGIN
    DELETE FROM tid_user_override WHERE tid_user_id = p_tidUserId;
END