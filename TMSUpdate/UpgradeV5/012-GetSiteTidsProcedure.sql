--multiline;
CREATE PROCEDURE get_site_tids(IN site_id INT)
BEGIN
	SELECT t.tid_id,
		   t.serial,
           t.PIN,
           t.ExpiryDate,
           t.ActivationDate,
           t.target_package_id,
           t.update_date
    FROM tid t
    LEFT JOIN tid_site ts ON ts.tid_id = t.tid_id
    WHERE ts.site_id = site_id;
END;