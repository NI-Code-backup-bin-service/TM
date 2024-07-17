--multiline;
CREATE PROCEDURE add_and_update_tid_updates_and_flag(IN tidUpdateId int, IN tidId int, IN targetPackageId int, IN updateDate datetime, IN thirdPartyApk text, IN dataValue varchar(45))
BEGIN
    INSERT INTO tid_updates(tid_update_id, tid_id, target_package_id, update_date, third_party_apk) values (tidUpdateId, tidId, targetPackageId, updateDate, thirdPartyApk);
    UPDATE tid set flag_status=true, flagged_date=CURRENT_TIMESTAMP where tid_id = tidId;
    SELECT apk_id, `name` from third_party_apks where name  LIKE CONCAT('%', dataValue , '%');
END