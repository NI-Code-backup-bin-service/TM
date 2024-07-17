--multiline
CREATE PROCEDURE GET_DATA_GROUPS_FOR_SITE_BY_PROFILE(IN p_profile_id int)
BEGIN
    SELECT dg.data_group_id,
           dg.name,
           dg.displayname_en,
           IF(COUNT(sp.profile_id) > 0, true, false) selected,
           IF(COUNT(sp.profile_id) > 1, true, false) preselected
    FROM data_group dg
             LEFT OUTER JOIN profile_data_group pdg on dg.data_group_id = pdg.data_group_id
             LEFT OUTER JOIN (
        select sp2.profile_id
        from site_profiles sp
                 INNER JOIN site_profiles sp2 ON sp.site_id = sp2.site_id
             #We don't care what data groups global has selected
        where sp.profile_id = p_profile_id
          AND sp2.profile_id != 1
    ) sp ON sp.profile_id = pdg.profile_id
    GROUP BY dg.data_group_id, dg.name, dg.displayname_en;
END;