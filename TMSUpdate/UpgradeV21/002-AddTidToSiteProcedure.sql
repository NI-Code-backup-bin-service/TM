--multiline
CREATE PROCEDURE `add_tid_to_site`(
IN newTid int, 
IN serial_number varchar(50), 
IN site INT
)
BEGIN
	IF (site is not null) THEN
		INSERT INTO tid (tid_id, `serial`) VALUES (newTid, serial_number);
		INSERT INTO tid_site VALUES (newTid, site, NULL, Now());
    END IF;
END