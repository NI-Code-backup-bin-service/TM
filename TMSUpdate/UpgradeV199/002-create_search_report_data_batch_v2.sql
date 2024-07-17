--multiline;
CREATE PROCEDURE `search_report_data_batch_v2`( IN siteIds LONGTEXT, IN updatesFrom BIGINT )
BEGIN
select concat("With SiteDataElementVersions as (SELECT d2.site_id, d2.data_element_id, Max(d2.version) as version FROM site_data d2 WHERE d2.site_id in (",siteIds,") and d2.priority = 2 group by d2.site_id, d2.data_element_id)
SELECT
	e.site_id,
	e.profile_id,
	e.data_group_id,
	dg.name AS 'data_group',
	e.data_element_id,
	e.name,
	v2.level AS `source`,
	v2.datavalue AS `datavalue`,
	v2.overriden AS `overriden`,
	e.datatype,
	e.is_allow_empty,
	e.max_length,
	e.validation_expression,
	e.validation_message,
	e.front_end_validate,
	e.options,
	e.display_name
FROM site_data_elements as e
INNER JOIN data_group as dg ON dg.data_group_id = e.data_group_id
LEFT JOIN site_data v2 ON v2.site_id = e.site_id AND v2.data_element_id = e.data_element_id AND v2.priority = 2
INNER JOIN SiteDataElementVersions as sdv on sdv.site_id = v2.site_id and sdv.data_element_id = v2.data_element_id and v2.version = sdv.version
WHERE e.site_id IN (",siteIds,")") into @sql;
PREPARE stmt1 FROM @sql;
EXECUTE stmt1 ;
DEALLOCATE PREPARE stmt1;
END