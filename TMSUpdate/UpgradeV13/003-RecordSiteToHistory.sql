--multiline;
CREATE PROCEDURE `record_site_to_history`(
  IN profile_id int,
  IN change_type INT,
  IN dataValue TEXT,
  IN updated_by varchar(255),
  IN approved int)
BEGIN
  INSERT INTO approvals(profile_id, data_element_id, change_type, current_value, new_value, created_at, approved_at, created_by, approved_by, approved)
  VALUES (profile_id, 1, change_type, current_value, dataValue, NOW(), NOW(), updated_by, updated_by, approved);
END