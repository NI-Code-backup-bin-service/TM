--multiline;
CREATE PROCEDURE `update_tid` (IN tid_id INT(11), IN presence INT(11))
BEGIN
	update tid set tid.Presence = presence
    where tid.tid_id = tid_id;
END;
