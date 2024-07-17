--multiline;
CREATE PROCEDURE discard_change(IN approvalID int, IN approval_user varchar(256))
BEGIN
    SELECT profile_id,change_type,current_value INTO @profileId,@changeType,@currentValue from approvals a WHERE a.approval_id = approvalID;
    SELECT ps.name,ps.group_id,psg.name INTO @paymentServiceName,@groupId,@paymentServiceGroupName from payment_service ps JOIN payment_service_group psg ON psg.group_id = ps.group_id where ps.service_id = @currentValue;
    SET @terminalId=(select name from profile where profile_id=@profileId and profile_type_id = (SELECT profile_type_id FROM profile_type WHERE name='tid'));
    IF @terminalId IS NOT NULL THEN
        UPDATE approvals a SET a.tid_id=@terminalId,a.approved = -1, approved_by = approval_user, approved_at = NOW() WHERE a.approval_id = approvalID;
    ELSEIF @paymentServiceName IS NOT NULL AND @changeType=11 THEN
        UPDATE approvals a SET a.current_value = @paymentServiceName,a.new_value = @paymentServiceGroupName,a.approved = -1, approved_by = approval_user, approved_at = NOW() WHERE a.approval_id = approvalID;
    ELSE
        UPDATE approvals a SET a.approved = -1, approved_by = approval_user, approved_at = NOW() WHERE a.approval_id = approvalID;
    END IF;
END;