--multiline;
CREATE PROCEDURE delete_tid(IN tid INT, IN site int)
  BEGIN
    DELETE FROM tid_site WHERE tid_id = tid AND site_id = site;
    DELETE FROM tid_updates WHERE tid_id = tid;

    IF NOT EXISTS(SELECT * FROM tid_site WHERE tid_id = tid)
    THEN
      DELETE FROM tid WHERE tid_id = tid;
    END IF;
  END