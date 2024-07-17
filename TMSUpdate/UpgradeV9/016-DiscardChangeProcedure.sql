--multiline;
CREATE PROCEDURE discard_change(IN approvalID int, IN approval_user VARCHAR(256))
BEGIN
	UPDATE approvals a SET a.approved = -1, approved_by = approval_user, approved_at = NOW() WHERE a.approval_id = approvalID;
END