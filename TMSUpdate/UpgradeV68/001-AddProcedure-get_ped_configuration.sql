--multiline
CREATE PROCEDURE get_ped_configuration(IN P_TID int, IN P_lastChecked bigint)
BEGIN
    DECLARE v_updatesFromDate DATETIME;
    IF ((SELECT updated_at FROM tid_site WHERE tid_id = P_TID) > FROM_UNIXTIME(P_lastChecked))
        OR
       ((select s.updated_at from tid_site ts INNER JOIN site s ON ts.site_id = s.site_id WHERE ts.tid_id = P_TID) > FROM_UNIXTIME(P_lastChecked))
    THEN
        SET v_updatesFromDate = FROM_UNIXTIME(0);
    ELSE
        SET v_updatesFromDate = FROM_UNIXTIME(P_lastChecked);
    END IF;

    SELECT
        JSON_ARRAYAGG(JSON_OBJECT(
                'DataGroupId', dataGroups.DataGroupId,
                'DataGroup', dataGroups.DataGroupName,
                'DataElements', dataGroups.DataElements
            ))
        /*
        We are selecting from another select because MySQL does not allow nesting of JSON aggregate functions
        */
    FROM (
             SELECT
                 dg.data_group_id 'DataGroupId',
                 dg.name 'DataGroupName',
                 JSON_ARRAYAGG(JSON_OBJECT(
                         'Name', de.name,
                         'Type', de.datatype,
                         'DataValue', pd.datavalue,
                         'Encrypted', if(de.is_password = 0 OR de.is_password = '' OR de.is_password IS NULL, CAST(false AS JSON), CAST(true AS JSON))
                     )) 'DataElements'
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
                 /*
                     This self join is done to filter out the highest priority profile type.
                     This would normally be done with a Window function but they are not available until MySQL 8.
                 */
                      INNER JOIN (
                 SELECT
                     de2.data_element_id,
                     tsp2.tid_id,
                     MIN(pt2.priority) 'profileTypePriority'
                 FROM tid_site_profiles tsp2
                          INNER JOIN profile p2 ON
                         tsp2.profile_id = p2.profile_id
                          INNER JOIN profile_type pt2 ON
                         p2.profile_type_id = pt2.profile_type_id
                          LEFT JOIN profile_data pd2 ON
                         p2.profile_id = pd2.profile_id
                          INNER JOIN data_element de2 ON
                         pd2.data_element_id = de2.data_element_id
                          INNER JOIN data_group dg2 ON
                         de2.data_group_id = dg2.data_group_id
		 WHERE tsp2.tid_id = P_TID
                 GROUP BY de2.data_element_id, tsp2.tid_id) SELF_1 ON
                         de.data_element_id = SELF_1.data_element_id
                     AND
                         tsp.tid_id = SELF_1.tid_id
                     AND
                         pt.priority = SELF_1.profileTypePriority
             WHERE
                     tsp.tid_id = P_TID
               AND
                 (
                         pd.updated_at IS NULL
                         OR
                         pd.updated_at >= v_updatesFromDate
                     )
             GROUP BY dg.data_group_id, dg.name
             ORDER BY dg.data_group_id
         ) dataGroups;
END;
