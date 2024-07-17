--multiline
CREATE PROCEDURE `fetch_tid_details` (IN tid INT)
BEGIN
	SELECT 
		t.software_version,
        t.firmware_version,
        t.last_transaction_time
    FROM tid t
    WHERE tid_id = tid;
END
