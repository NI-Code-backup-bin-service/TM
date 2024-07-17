--multiline;
CREATE PROCEDURE `update_ped_time_and_get_tid_details`(IN tidID INT, IN lastCheckedTime BIGINT)
BEGIN
    DECLARE previousCheckedTime, currentTime BIGINT;
    DECLARE flagStatus, eodAuto BOOL;
    DECLARE flaggedDate, autoTime TEXT;
    SELECT last_checked_time, UNIX_TIMESTAMP()*1000, flag_status, flagged_date, eod_auto, auto_time FROM tid WHERE tid_id = tidID INTO previousCheckedTime, currentTime, flagStatus, flaggedDate, eodAuto, autoTime;
    IF previousCheckedTime = lastCheckedTime THEN
        UPDATE tid SET confirmed_time = lastCheckedTime, last_checked_time = currentTime WHERE tid_id = tidID;
    ELSE
        UPDATE tid SET last_checked_time = currentTime WHERE tid_id = tidID;
    END IF;
    SELECT currentTime, flagStatus, flaggedDate, eodAuto, autoTime;
END;