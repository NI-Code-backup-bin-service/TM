--multiline;
CREATE PROCEDURE get_datagroup_with_profile(
   IN acquirerId INT,
   IN chainId INT,
   IN profileId INT,
   IN isChain BOOL
)
BEGIN
   IF isChain =false THEN
        SELECT
            distinct(dg.data_group_id) as data_group_id,
            dg.name,
            dg.displayname_en,
            CASE WHEN pdg.data_group_id = dg.data_group_id THEN true ELSE false END as preSelected,
            CASE WHEN spdg.data_group_id = dg.data_group_id THEN true ELSE false END as isSelected
        FROM data_group as dg
        LEFT JOIN profile_data_group pdg ON pdg.data_group_id = dg.data_group_id AND pdg.profile_id IN(acquirerId,chainId)
        LEFT JOIN profile_data_group spdg ON spdg.data_group_id = dg.data_group_id AND spdg.profile_id = profileId;
   ELSE
        SELECT
            distinct(dg.data_group_id) as data_group_id,
            dg.name,
            dg.displayname_en,
            CASE WHEN pdg.data_group_id = dg.data_group_id THEN true ELSE false END as preSelected,
            CASE WHEN spdg.data_group_id = dg.data_group_id THEN true ELSE false END as isSelected
        FROM data_group as dg
        LEFT JOIN profile_data_group pdg ON pdg.data_group_id = dg.data_group_id AND pdg.profile_id IN(acquirerId)
        LEFT JOIN profile_data_group spdg ON spdg.data_group_id = dg.data_group_id AND spdg.profile_id = profileId;
   END IF;
END;