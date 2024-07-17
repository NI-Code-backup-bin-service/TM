--multiline;
CREATE PROCEDURE set_tid_update_date(IN tid int, update_date DATE)
BEGIN
	UPDATE tid t SET t.update_date = update_date
    WHERE t.tid_id = tid;
END;