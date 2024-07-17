--multiline;
CREATE PROCEDURE check_data_group_exist_in_profile(IN profileId INT, IN dataGroupName VARCHAR(255))
BEGIN
    SELECT COUNT(dg.data_group_id)
    FROM profile_data_group pdg
    JOIN data_group dg
        ON dg.data_group_id = pdg.data_group_id
    WHERE pdg.profile_id = profileId
    AND dg.name = dataGroupName;
END;