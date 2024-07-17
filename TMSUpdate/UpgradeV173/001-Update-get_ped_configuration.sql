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
    FROM
    (
		SELECT
			JSONd.DataGroupId,
			JSONd.DataGroupName,
			JSON_ARRAYAGG(
				JSON_OBJECT(
					'Name', JSONd.Name,
					'Type', JSONd.Type,
					'DataValue', JSONd.DataValue,
					'Encrypted', JSONd.Encrypted
				)
			) 'DataElements'
		FROM
		(
			SELECT
				*
			FROM
			(
				SELECT
					dg.data_group_id 'DataGroupId',
					dg.name 'DataGroupName',
					de.name 'Name',
					de.datatype 'Type',
					if(pd.is_encrypted = 0 OR pd.is_encrypted = '' OR pd.is_encrypted IS NULL, CAST(false AS JSON), CAST(true AS JSON)) 'Encrypted',
                    FIRST_VALUE(pd.datavalue) OVER(PARTITION BY dg.data_group_id, de.name ORDER BY pt.priority) 'DataValue'
				FROM
					tid_site_profiles tsp
						INNER JOIN
					profile p ON tsp.profile_id = p.profile_id
						INNER JOIN
					profile_type pt ON p.profile_type_id = pt.profile_type_id
						LEFT JOIN
					profile_data pd ON p.profile_id = pd.profile_id
						INNER JOIN
					data_element de ON pd.data_element_id = de.data_element_id
						INNER JOIN
					data_group dg ON de.data_group_id = dg.data_group_id
				WHERE
					tsp.tid_id = P_TID
					AND
					(
						pd.updated_at IS NULL
						OR pd.updated_at >= FROM_UNIXTIME(0)
					)
				ORDER BY
					dg.data_group_id
			) dataSet
			GROUP BY
				dataSet.DataGroupId,
				dataSet.Name
		) JSONd
		GROUP BY
			JSONd.DataGroupId
	) dataGroups;
END;
