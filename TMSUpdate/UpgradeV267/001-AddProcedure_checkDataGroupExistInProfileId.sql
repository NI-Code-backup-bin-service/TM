--multiline;
CREATE PROCEDURE checkDataGroupExistInProfileId(IN profileId INT, IN groupName VARCHAR(255))
BEGIN
    SELECT data_group_id FROM profile_data_group WHERE profile_id = profileId AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = groupName);
END;