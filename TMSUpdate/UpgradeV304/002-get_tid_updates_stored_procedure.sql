--multiline;
CREATE PROCEDURE get_tid_updates(IN tidUpdateId int, IN tidId int)
BEGIN
    declare tid_updateId int;
    SET tid_updateId=(SELECT tid_update_id FROM tid_updates where tid_update_id=tidUpdateId and tid_id=tidId);
    IF tid_updateId > 0 THEN
        SELECT tid_update_id, target_package_id, DATE_FORMAT(update_date, '%Y-%m-%d %H:%i'), third_party_apk FROM tid_updates where tid_update_id=tidUpdateId and tid_id = tidId;
    ELSE
        SELECT tid_update_id, target_package_id, DATE_FORMAT(update_date, '%Y-%m-%d %H:%i'), third_party_apk FROM tid_updates where tid_id = tidId order by update_date desc limit 1;
    END IF;
END