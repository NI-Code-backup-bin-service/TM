--multiline;
CREATE PROCEDURE `save_bulk_tid_flagging`(IN fileName TEXT, IN createdBy varchar(45))
BEGIN
    INSERT INTO bulk_tid_flagging (filename, created_by, created_at) VALUES (fileName, createdBy, NOW());
END