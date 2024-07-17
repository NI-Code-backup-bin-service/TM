--multiline;
CREATE PROCEDURE `add_tid_user`(IN tidId INT, IN Username VARCHAR(255), IN Pin VARCHAR(5), IN Modules VARCHAR(255), IN Encrypted TINYINT )
BEGIN
	INSERT INTO tid_user_override(tid_id, Username, PIN, Modules, is_encrypted) VALUES (tidId, Username, Pin, Modules, Encrypted);
END