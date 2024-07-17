--multiline;
CREATE PROCEDURE discard_change(IN approvalID int, IN approval_user varchar(256))
BEGIN
    SET @profileId = (SELECT profile_id from approvals a WHERE a.approval_id = approvalID);
    SET @terminalId=(select name from profile where profile_id=@profileId and profile_type_id=5);
    UPDATE approvals a SET a.approved = -1, approved_by = approval_user, approved_at = NOW() WHERE a.approval_id = approvalID;
    IF @terminalId IS NOT NULL THEN
        UPDATE approvals a SET a.tid_id=@terminalId WHERE a.approval_id = approvalID;
    END IF;
END;