--multiline
CREATE PROCEDURE `add_datagroup_to_profile`(profileId INT, datagroupId INT, user VARCHAR(256))
BEGIN
	INSERT INTO profile_data_group (profile_id, data_group_id, version, updated_at, updated_by, created_at, created_by) values (profileId, datagroupId, 1, NOW(), user, NOW(), user);
END