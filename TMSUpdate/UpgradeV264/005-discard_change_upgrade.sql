--multiline;
CREATE PROCEDURE discard_change(IN approvalID int, IN approval_user varchar(256))
BEGIN
    SET @profileId = (SELECT profile_id from approvals a WHERE a.approval_id = approvalID);
    SET @paymentServiceName = (SELECT name from payment_service where service_id = (SELECT current_value from approvals a WHERE a.approval_id = approvalID));
    SET @paymentServiceGroupName = (SELECT name from payment_service_group where group_id = (SELECT group_id from payment_service where service_id = (SELECT current_value from approvals a WHERE a.approval_id = approvalID)));
    SET @terminalId=(select name from profile where profile_id=@profileId and profile_type_id=5);
    UPDATE approvals a SET a.approved = -1, approved_by = approval_user, approved_at = NOW() WHERE a.approval_id = approvalID;
    IF @terminalId IS NOT NULL THEN
        UPDATE approvals a SET a.tid_id=@terminalId WHERE a.approval_id = approvalID;
    END IF;
    IF @paymentServiceName IS NOT NULL THEN
        UPDATE approvals a SET a.current_value = @paymentServiceName,a.new_value = @paymentServiceGroupName WHERE a.approval_id = approvalID;
    END IF;
END;