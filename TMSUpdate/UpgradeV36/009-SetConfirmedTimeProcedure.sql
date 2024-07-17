--multiline;
CREATE PROCEDURE `set_confirmed_time`(IN tid INT, IN checkTime BIGINT)
BEGIN
    UPDATE tid AS t SET confirmed_time = checkTime WHERE t.tid_id = tid;
END