--multiline;
CREATE PROCEDURE get_tid_updates(IN tidUpdateId int, IN tidId int)
BEGIN
    SELECT tid_update_id, target_package_id, DATE_FORMAT(update_date, '%Y-%m-%d %H:%i'), third_party_apk INTO @tidUpdateID, @targetPackageID, @updateDate, @thirdPartyAPK FROM tid_updates where tid_update_id=tidUpdateId and tid_id = tidId;
    IF @tidUpdateID > 0 THEN
        SELECT @tidUpdateID, @targetPackageID, @updateDate, @thirdPartyAPK;
    ELSE
        SELECT tid_update_id, target_package_id, DATE_FORMAT(update_date, '%Y-%m-%d %H:%i'), third_party_apk FROM tid_updates where tid_id = tidId order by update_date desc limit 1;
    END IF;
END