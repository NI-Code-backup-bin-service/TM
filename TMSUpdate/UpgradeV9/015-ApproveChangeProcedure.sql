--multiline;
CREATE PROCEDURE approve_change(IN approval_id int, IN approval_user Varchar(256))
BEGIN
	DECLARE approval_type int;
    DECLARE newVal VARCHAR(256);
    DECLARE profileId INT;
    DECLARE elementId INT;
    
	SET profileId = (SELECT profile_id from approvals a WHERE a.approval_id = approval_id);
    SET elementId = (SELECT data_element_id from approvals a WHERE a.approval_id = approval_id);
    SET approval_type = (SELECT change_type FROM approvals a WHERE a.approval_id = approval_id);
        
    IF approval_type = 1 OR  approval_type = 2 THEN
		SET newVal = (SELECT new_value from approvals a WHERE a.approval_id = approval_id);
        IF EXISTS (SELECT profile_data_id FROM profile_data pd WHERE pd.profile_id = profileId AND pd.data_element_id = elementId) THEN
			UPDATE profile_data pd SET pd.datavalue = newVal, pd.updated_at=NOW(), pd.updated_by = approval_user 
            WHERE pd.profile_id = profileId AND pd.data_element_id = elementId;
		else
            insert into profile_data( profile_id, data_element_id, datavalue, version,  updated_at,updated_by, created_at, created_by, approved, overriden) 
            values (profileId, 
            elementId, 
            newVal, 
            1, 
            current_timestamp, 
            approval_user, 
            current_timestamp, 
            approval_user, 
            1, 
            CASE approval_type WHEN 2 THEN 1 ELSE 0 END); -- override if change type is overriden
		end if;
	ELSEIF approval_type=4 THEN
		DELETE pd FROM profile_data pd WHERE pd.profile_id = profileId AND pd.data_element_id = elementId;
	end if;
    
    UPDATE approvals a SET approved = 1, approved_by = approval_user, approved_at = NOW() WHERE a.approval_id = approval_id;
END