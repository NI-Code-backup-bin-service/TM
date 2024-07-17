--multiline;
CREATE PROCEDURE `save_pending_profile_change`(
  IN profileId int,
  IN change_type INT,
  IN dataValue TEXT,
  IN updated_by varchar(255),
  IN tidId TEXT,
  IN approvalState int)
BEGIN
  IF (SELECT count(a.approval_id) FROM approvals a 
					WHERE a.profile_id = profileId 
                    AND a.data_element_id = 1 
                    AND a.tid_id = tidId 
                    AND a.change_type = change_type 
                    AND a.approved = 0)  > 0
  THEN 
	UPDATE approvals ap
		SET
		  new_value = dataValue,
		  created_at = NOW(),
		  created_by = updated_by
		WHERE ap.profile_id = profileId
		  AND ap.data_element_id = 1
		  AND ap.tid_id = tidId
		  AND a.change_type = change_type 
		  AND ap.approved = 0;
  ELSE
    INSERT INTO approvals(profile_id, data_element_id, change_type, current_value, new_value, created_at, created_by, approved, tid_id, approved_at, approved_by)
    VALUES (profileId, 1, change_type, current_value, dataValue, NOW(), updated_by, approvalState, tidId, CASE WHEN approvalState = 1 THEN NOW() ELSE NULL END, CASE WHEN approvalState = 1 THEN updated_by ELSE NULL END );
  END IF;
END