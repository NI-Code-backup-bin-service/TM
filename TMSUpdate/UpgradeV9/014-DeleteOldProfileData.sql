--multiline;
DELETE pd FROM profile_data pd
LEFT JOIN (SELECT pd2.pdid as pdid,
						pd2.deid as deid,
						pd2.pid as pid,
						pd2.ver as ver,
						pd2.approved as approved
FROM (	SELECT d.profile_data_id as pdid,
				d.data_element_id as deid,
				d.profile_id as pid,
				d.version as ver,
				d.approved as approved
FROM  profile_data d
) as pd2

) as profile_data_2
ON profile_data_2.pid = pd.profile_id
AND profile_data_2.deid = pd.data_element_id
AND profile_data_2.ver < pd.version
AND profile_data_2.approved > 0
WHERE profile_data_2.pdid IS NOT NULL;