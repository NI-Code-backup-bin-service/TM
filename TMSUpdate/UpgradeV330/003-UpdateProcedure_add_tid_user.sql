--multiline
create procedure add_tid_user(IN p_tidUserId int, IN p_tidId int, IN p_Username varchar(255), IN p_Pin varchar(66), IN p_Modules LONGTEXT, IN p_Encrypted tinyint(1))
BEGIN
    /*Positive but non existent rows are copies from the site data at time of first override, so add them in*/
    IF p_tidUserId >= 0 AND (SELECT COUNT(*) FROM tid_user_override t_o WHERE p_tidUserId = t_o.tid_user_id) > 0 THEN
        UPDATE tid_user_override SET Pin = p_Pin, Username = p_Username, Modules = p_Modules WHERE tid_user_id = p_tidUserId;
    ELSE
        INSERT INTO tid_user_override(tid_id, Username, PIN, Modules, is_encrypted) VALUES (p_tidId, p_Username, p_Pin, p_Modules, p_Encrypted);
    END IF;
    /*Return the user ID of the just inserted user*/
    SELECT tid_user_id FROM tid_user_override WHERE username = Username AND tid_id = p_tidId;
END;

