--multiline
create
    procedure store_reset_PIN(IN tid int, IN PIN varchar(5), IN timeout int)
BEGIN

    UPDATE tid t SET t.reset_pin = PIN, t.reset_pin_expiry_date = DATE_ADD(NOW(), INTERVAL timeout MINUTE)

    WHERE t.tid_id  = tid;

END;