--multiline;
CREATE PROCEDURE save_tid(IN tid int, IN serial_number varchar(50), IN site INT)
BEGIN
  INSERT INTO tid (tid_id, `serial`) VALUES (tid, serial_number);
  INSERT INTO tid_site VALUES (tid, site, NULL);
END