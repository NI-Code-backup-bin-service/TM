--multiline
CREATE PROCEDURE `get_configured_data_elements_for_tid`(IN P_TID int, IN P_lastChecked bigint)
BEGIN
    IF ((SELECT updated_at FROM tid_site WHERE tid_id = P_TID) > FROM_UNIXTIME(P_lastChecked))
        OR
       ((select s.updated_at from tid_site ts INNER JOIN site s ON ts.site_id = s.site_id WHERE ts.tid_id = P_TID) > FROM_UNIXTIME(P_lastChecked))
    THEN
        SET @v_updatesFromDate = FROM_UNIXTIME(0);
    ELSE
        SET @v_updatesFromDate = FROM_UNIXTIME(P_lastChecked);
    END IF;

    SELECT
        dg.name data_group_name,
        de.name data_element_name,
        de.datatype,
        pt.priority,
        pd.datavalue,
        if(pd.is_encrypted = 0 OR pd.is_encrypted = '' OR pd.is_encrypted IS NULL, 0, 1)
    FROM tid_site_profiles tsp
             INNER JOIN profile p ON
            tsp.profile_id = p.profile_id
             INNER JOIN profile_type pt ON
            p.profile_type_id = pt.profile_type_id
             LEFT JOIN profile_data pd ON
            p.profile_id = pd.profile_id
             INNER JOIN data_element de ON
            pd.data_element_id = de.data_element_id
             INNER JOIN data_group dg ON
            de.data_group_id = dg.data_group_id
    WHERE
            tsp.tid_id = P_TID
      AND
        (
                pd.updated_at IS NULL
                OR
                pd.updated_at >= @v_updatesFromDate
            )
    ORDER BY dg.data_group_id;
END