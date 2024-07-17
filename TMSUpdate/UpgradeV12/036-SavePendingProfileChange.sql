--multiline;
CREATE PROCEDURE `save_pending_profile_change`(
  IN profile_id int,
  IN change_type INT,
  IN dataValue TEXT,
  IN updated_by varchar(255),
  IN tidId TEXT,
  IN approved int)
BEGIN
  INSERT INTO approvals(profile_id, data_element_id, change_type, current_value, new_value, created_at, created_by, approved, tid_id)
  VALUES (profile_id, 1, change_type, current_value, dataValue, NOW(), updated_by, approved, tidId);
END