--multiline;
CREATE PROCEDURE `set_last_checked_time`(IN tid INT, IN checkTime BIGINT)
BEGIN
    UPDATE tid AS t SET last_checked_time = checkTime WHERE t.tid_id = tid;
END