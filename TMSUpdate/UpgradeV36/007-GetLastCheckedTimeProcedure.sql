--multiline;
CREATE PROCEDURE `get_last_checked_time`(IN tid INT)
BEGIN
    SELECT t.last_checked_time
    FROM tid AS t
    WHERE t.tid_id = tid;
END