--multiline;
INSERT
	IGNORE INTO profile_data (
		profile_id,
		datavalue,
		data_element_id,
		version,
		updated_at,
		updated_by,
		approved
	)
SELECT
	pd.profile_id,
	'false',
	de.data_element_id,
	1,
	NOW(),
	'system',
	1
FROM
	(
		SELECT
			DISTINCT profile_id
		FROM
			profile_data
	) pd
	CROSS JOIN data_element de
WHERE
	de.datatype = 'BOOLEAN' AND de.tid_overridable = 1
	AND de.data_group_id IN (
		SELECT
			data_group_id
		FROM
			profile_data_group
		where
			profile_id = pd.profile_id
	);