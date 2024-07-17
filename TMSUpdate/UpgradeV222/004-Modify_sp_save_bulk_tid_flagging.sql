--multiline
CREATE PROCEDURE `save_bulk_tid_flagging`(IN fileName TEXT, IN fileType varchar(45) ,IN createdBy varchar(45))
BEGIN
INSERT INTO bulk_approvals (filename, filetype, created_by, created_at) VALUES (fileName, fileType, createdBy, NOW());
END