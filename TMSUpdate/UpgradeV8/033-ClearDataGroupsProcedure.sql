--multiline
CREATE PROCEDURE `clear_profile_datagroups`(profileId INT)
BEGIN
	DELETE FROM profile_data_group WHERE profile_id = profileId;
END