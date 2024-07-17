--multiline;
CREATE PROCEDURE insert_tid_update(IN tidUpdateId INT, IN tidId INT, IN targetPackageId INT, IN updateDate DATETIME, IN thirdPartyApk INT)
BEGIN
if (NOT EXISTS (select tid_update_id from tid_updates t where tid_update_id = tidUpdateId and tid_id = tidId)) then
    insert into tid_updates(tid_update_id, tid_id, target_package_id, update_date, third_party_apk)
    values (tidUpdateId, tidId, targetPackageId, updateDate, thirdPartyApk);
else
	update tid_updates set target_package_id = targetPackageId, update_date = updateDate, third_party_apk = thirdPartyApk
    where tid_update_id = tidUpdateId and tid_id = tidId;
end if;
END