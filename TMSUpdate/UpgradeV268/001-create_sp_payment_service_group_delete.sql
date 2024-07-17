--multiline;
CREATE PROCEDURE payment_service_group_deletion(IN profileId int,IN changeType int,IN groupId varchar(256),IN groupName varchar(256),IN userName varchar(256),acquirerName varchar(256))
BEGIN
    DECLARE approvalCount INT;
    SELECT count(*) into approvalCount from approvals where current_value = groupId and  approved=0 and change_type= changeType;
    IF approvalCount > 0 THEN
        UPDATE approvals set created_by= userName,created_at=NOW() where current_value = groupId and change_type= changeType AND new_value = groupName;
    ELSE
        INSERT INTO approvals (profile_id,data_element_id, change_type, current_value, new_value, created_at, approved, created_by, acquirer,approved_by,approved_at)
		VALUE
		(profileId, 1, changeType,groupId, groupName, NOW(), 0, userName, acquirerName,userName,NOW());
    END IF;
END;

