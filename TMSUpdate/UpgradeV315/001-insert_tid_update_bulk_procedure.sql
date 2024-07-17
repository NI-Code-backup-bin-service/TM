-- --multiline;
CREATE PROCEDURE `insert_tid_update_bulk`(
    IN tidUpdateIdList VARCHAR(255),
    IN tidIdList VARCHAR(255),
    IN targetPackageIdList VARCHAR(255),
    IN updateDateList DATETIME,
    IN thirdPartyApkList TEXT
)
BEGIN
    DECLARE i INT DEFAULT 1;
    DECLARE numUpdates INT;

    SET numUpdates = LENGTH(tidUpdateIdList) - LENGTH(REPLACE(tidUpdateIdList, ',', '')) + 1;

    WHILE i <= numUpdates DO
        SET @currentTidUpdateId = SUBSTRING_INDEX(SUBSTRING_INDEX(tidUpdateIdList, ',', i), ',', -1);
        SET @currentTidId = SUBSTRING_INDEX(SUBSTRING_INDEX(tidIdList, ',', i), ',', -1);
        SET @currentTargetPackageId = SUBSTRING_INDEX(SUBSTRING_INDEX(targetPackageIdList, ',', i), ',', -1);
        SET @currentUpdateDate = SUBSTRING_INDEX(SUBSTRING_INDEX(updateDateList, ',', i), ',', -1);
        SET @currentThirdPartyApk = SUBSTRING_INDEX(SUBSTRING_INDEX(thirdPartyApkList, ',', i), ',', -1);

        IF NOT EXISTS (SELECT tid_update_id FROM tid_updates t WHERE tid_update_id = @currentTidUpdateId AND tid_id = @currentTidId) THEN
            INSERT INTO tid_updates (tid_update_id, tid_id, target_package_id, update_date, third_party_apk)
            VALUES (@currentTidUpdateId, @currentTidId, @currentTargetPackageId, @currentUpdateDate, @currentThirdPartyApk);
        ELSE
            UPDATE tid_updates
            SET target_package_id = @currentTargetPackageId,
                update_date = @currentUpdateDate,
                third_party_apk = @currentThirdPartyApk
            WHERE tid_update_id = @currentTidUpdateId AND tid_id = @currentTidId;
        END IF;

        SET i = i + 1;
    END WHILE;
END;

