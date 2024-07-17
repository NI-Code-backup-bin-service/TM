--multiline
CREATE PROCEDURE `get_configured_data_elements_for_tid`(IN P_TID INT, IN P_lastChecked BIGINT)
BEGIN
	IF ((SELECT updated_at FROM tid_site WHERE tid_id = P_TID) > FROM_UNIXTIME(P_lastChecked))
		OR ((SELECT s.updated_at FROM tid_site ts INNER JOIN site s ON ts.site_id = s.site_id WHERE ts.tid_id = P_TID) > FROM_UNIXTIME(P_lastChecked))
	THEN
		SET @v_updatesFromDate = FROM_UNIXTIME(0);
	ELSE
		SET @v_updatesFromDate = FROM_UNIXTIME(P_lastChecked);
	END IF;
	SELECT # A second select is used to exclude pd.updated_at (needed for HAVING)
		data_group_name,
		data_element_name,
		datatype,
		priority,
		datavalue,
		encrypted
	FROM (
		SELECT
			dg.name data_group_name,
			de.name data_element_name,
			de.datatype,
			pt.priority,
			pd.datavalue,
			if(pd.is_encrypted = 0 OR pd.is_encrypted = '' OR pd.is_encrypted IS NULL, 0, 1) encrypted,
			pd.updated_at
		FROM (
		    SELECT
		        `tspid`.`profile_id` AS `profile_id`,
		        `tspid`.`tid_id` AS `tid_id`,
		        `tspid`.`site_id` AS `site_id`
		    FROM (
		        SELECT
		            `ts`.`tid_profile_id` AS `profile_id`,
		            `ts`.`tid_id` AS `tid_id`,
		            `ts`.`site_id` AS `site_id`
		        FROM `tid_site` AS `ts`
		            LEFT JOIN `site_profiles` AS `sp` ON `ts`.`site_id` = `sp`.`site_id`
		        WHERE `ts`.`tid_id` = P_TID
		        UNION
		        SELECT
		            `sp`.`profile_id` AS `profile_id`,
		            `ts`.`tid_id` AS `tid_id`,
		            `sp`.`site_id` AS `site_id`
		        FROM `tid_site` AS `ts`
		            LEFT JOIN `site_profiles` AS `sp` ON `ts`.`site_id` = `sp`.`site_id`
		        WHERE `ts`.`tid_id` = P_TID
            ) AS `tspid`
        ) AS tsp
			INNER JOIN profile AS p ON tsp.profile_id = p.profile_id
			INNER JOIN profile_type AS pt ON p.profile_type_id = pt.profile_type_id
			LEFT JOIN profile_data AS pd ON p.profile_id = pd.profile_id
			INNER JOIN data_element AS de ON pd.data_element_id = de.data_element_id
			INNER JOIN data_group AS dg ON de.data_group_id = dg.data_group_id
		WHERE
			tsp.tid_id = P_TID
		GROUP BY dg.data_group_id, de.name
		HAVING pd.updated_at IS NULL OR pd.updated_at >= @v_updatesFromDate
		ORDER BY dg.data_group_id) AS data;
END