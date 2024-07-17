--multiline;
CREATE PROCEDURE `add_to_bulk_approvals` (IN fileName TEXT, IN fileType varchar(45) ,IN createdBy varchar(45),IN change_type INT, IN approved INT)
BEGIN
IF approved=-1 THEN
    INSERT INTO bulk_approvals (filename, filetype, created_by, created_at, change_type, approved, approved_by ,approved_at) VALUES (fileName, fileType, createdBy, NOW(), change_type, approved, createdBy, NOW());
ELSE
	INSERT INTO bulk_approvals (filename, filetype, created_by, created_at, change_type) VALUES (fileName, fileType, createdBy, NOW(),change_type);
END IF; 
END