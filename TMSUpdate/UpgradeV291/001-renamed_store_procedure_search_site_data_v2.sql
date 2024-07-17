--multiline;
CREATE PROCEDURE `search_report_data_batch_with_site_ids`(IN siteIds LONGTEXT, IN updatesFrom BIGINT)
BEGIN
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
INNER JOIN profile as p ON p.profile_id = e.profile_id
INNER JOIN profile_type pt ON pt.profile_type_id = p.profile_type_id AND pt.profile_type_id = (SELECT profile_type_id FROM profile_type WHERE name = 'site')
LEFT JOIN site_data v2 ON v2.site_id = e.site_id AND v2.data_element_id = e.data_element_id AND v2.priority = 2
WHERE FIND_IN_SET (e.site_id, siteIds);
END;