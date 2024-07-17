--multiline;
CREATE PROCEDURE fetch_active_group_by_profileId(IN profileId INT)
BEGIN
SELECT DISTINCT dg.name
FROM site_profiles sp
         INNER JOIN site_profiles sp2 ON sp.site_id = sp2.site_id
         INNER JOIN profile_data_group pdg ON pdg.profile_id = sp2.profile_id
         INNER JOIN data_group dg ON pdg.data_group_id = dg.data_group_id
WHERE sp.profile_id = profileId AND sp2.profile_id != 1;
END;