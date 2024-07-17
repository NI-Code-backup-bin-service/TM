--multiline;
CREATE PROCEDURE `search_report_data_batch`()
BEGIN
SELECT
	e.site_id,
	e.profile_id,
	e.data_group_id,
	dg.name AS 'data_group',
	e.data_element_id,
	e.name,
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
FROM profile_data as v2
INNER JOIN site_data_elements as e ON v2.profile_id = e.profile_id AND v2.data_element_id = e.data_element_id
INNER JOIN profile as p ON p.profile_id = e.profile_id
INNER JOIN profile_type pt ON pt.profile_type_id = p.profile_type_id AND pt.profile_type_id = (SELECT profile_type_id FROM profile_type WHERE name = 'site')
INNER JOIN data_group as dg ON dg.data_group_id = e.data_group_id;
END