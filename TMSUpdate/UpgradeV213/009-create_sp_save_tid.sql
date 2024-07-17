--multiline
CREATE PROCEDURE `save_tid`(IN tid int, IN serialNumber varchar(50), IN site INT,  IN eodAuto BOOLEAN, IN autoTime varchar(45))
BEGIN
    INSERT INTO tid (tid_id, serial, flag_status,flagged_date, eod_auto, auto_time) VALUES (tid, serialNumber, 1, now(), eodAuto, autoTime);
    INSERT INTO tid_site VALUES (tid, site, NULL, NOW());
END