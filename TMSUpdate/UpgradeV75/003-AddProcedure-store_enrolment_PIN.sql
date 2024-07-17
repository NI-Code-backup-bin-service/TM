--multiline
create
    procedure store_enrolment_PIN(IN tid int, IN PIN varchar(5), IN timeout int)
BEGIN

    UPDATE tid t SET t.PIN = PIN, t.ExpiryDate = DATE_ADD(NOW(), INTERVAL timeout MINUTE)

    WHERE t.tid_id  = tid;

END;