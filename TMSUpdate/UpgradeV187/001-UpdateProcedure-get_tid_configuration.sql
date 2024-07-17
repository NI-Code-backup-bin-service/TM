--multiline
CREATE PROCEDURE `get_tid_configuration`(IN P_TID int, IN P_lastChecked bigint)
BEGIN
    IF ((SELECT updated_at FROM tid_site WHERE tid_id = P_TID) > FROM_UNIXTIME(P_lastChecked))
        OR ((SELECT s.updated_at FROM tid_site ts INNER JOIN site s ON ts.site_id = s.site_id WHERE ts.tid_id = P_TID) > FROM_UNIXTIME(P_lastChecked))
    THEN
        SET @v_updatesFromDate = FROM_UNIXTIME(0);
ELSE
        SET @v_updatesFromDate = FROM_UNIXTIME(P_lastChecked);
END IF;

SELECT
    data_group_name,
    data_element_name,
    datatype,
    datavalue,
    encrypted
FROM (
         SELECT
             dg.name data_group_name,
             de.name data_element_name,
             de.datatype,
             # NOTE:
             # The below priority field is in fact in the inverse order of the priority that in pt.priority,
             # where 1 is the highest priority.
             # This priority field is dynamic and not tied to a particular profile type, so when chain is
             # the highest level that a field has been configured at, the priority will be 1; and when TID
             # is the highest level that a field has been configured at, the priority will be 1 too.
             # This has been done simply so we can filter by priority = 1 and always get the highest
             # priority profile_data datavalue.
             ROW_NUMBER() OVER (PARTITION BY de.data_element_id ORDER BY pt.priority asc) priority,
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
                      INNER JOIN profile_data AS pd ON p.profile_id = pd.profile_id
                      INNER JOIN data_element AS de ON pd.data_element_id = de.data_element_id
                      INNER JOIN data_group AS dg ON de.data_group_id = dg.data_group_id
			 WHERE
				 tsp.tid_id = P_TID
			 ORDER BY dg.data_group_id) AS data
             where data.priority = 1 and (data.updated_at IS NULL OR data.updated_at >= @v_updatesFromDate);
END