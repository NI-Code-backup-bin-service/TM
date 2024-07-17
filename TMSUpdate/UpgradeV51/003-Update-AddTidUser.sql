--multiline
CREATE PROCEDURE `add_tid_user`(IN p_tidUserId INT, IN p_tidId INT, IN p_Username VARCHAR(255), IN p_Pin VARCHAR(66), IN p_Modules VARCHAR(255), IN p_Encrypted bool )
BEGIN
    IF p_tidUserId <= 0 THEN
		UPDATE tid_user_override SET Pin = p_Pin AND Username = p_Username WHERE tid_user_id = p_tidUserId;
    ELSE
		INSERT INTO tid_user_override(tid_id, Username, PIN, Modules, is_encrypted) VALUES (p_tidId, p_Username, p_Pin, p_Modules, p_Encrypted);
    END IF;

    /*Return the user ID of the just inserted user*/
    SELECT tid_user_id FROM tid_user_override WHERE username = Username AND tid_id = p_tidId;
END