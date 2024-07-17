--multiline
CREATE PROCEDURE `save_tid`(IN tid int, IN serial_number varchar(50), IN site INT, IN is_encrypted BOOLEAN)
BEGIN
    INSERT INTO tid (`tid_id`, `serial`,`is_encrypted`) VALUES (tid, serial_number,is_encrypted);
    INSERT INTO tid_site VALUES (tid, site, NULL, NOW());
END