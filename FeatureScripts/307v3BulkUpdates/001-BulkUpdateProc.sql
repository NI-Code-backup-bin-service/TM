--multiline
CREATE PROCEDURE `bulkUpdate`(IN versionNum TEXT)
BEGIN
	SET @versionId = (SELECT package_id FROM package WHERE `version` = versionNum);

	INSERT INTO tid_updates (tid_update_id, tid_id, target_package_id, update_date)
	SELECT (SELECT IFNULL(MAX(tid_update_id), 0) FROM tid_updates tu WHERE tu.tid_id = t.tid_id) + 1, t.tid_id, @versionId, CURDATE()
	FROM tid t
	WHERE IFNULL((SELECT p.version FROM tid_updates u 
			LEFT JOIN package p ON p.package_id = u.target_package_id 
			WHERE u.tid_id = t.tid_id and t.tid_id IN (88881755)
			ORDER BY u.update_date DESC LIMIT 1), '') NOT IN ('71010', '71010BT')
	AND firmware_version = '19091701'
	AND software_version = '70041'
	LIMIT 1;
END