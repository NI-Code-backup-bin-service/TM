--multiline;
CREATE PROCEDURE insert_tid_update_ignore_TPAPK(IN tidId INT, IN targetPackageId INT, IN updateDate DATETIME, IN thirdPartyApk text)
BEGIN
    SET @tidUpdateId = ((SELECT IFNULL(MAX(tid_update_id), 0) FROM tid_updates WHERE tid_id = tidId)+1);

    IF (thirdPartyApk = '') THEN
        INSERT INTO tid_updates (tid_update_id, tid_id, target_package_id, update_date)
            VALUES (@tidUpdateId, tidId, targetPackageId, updateDate);
    ELSE
            INSERT INTO tid_updates(tid_update_id, tid_id, target_package_id, update_date, third_party_apk)
            VALUES (@tidUpdateId, tidId, targetPackageId, updateDate, thirdPartyApk);
    END IF;
END