--multiline;
CREATE PROCEDURE get_data_groups_by_profile_id(IN profile_id int)
BEGIN
    SELECT dg.data_group_id, dg.name FROM data_group dg
    LEFT JOIN profile_data_group pdg ON pdg.data_group_id = dg.data_group_id
    WHERE pdg.profile_id = profile_id;
END