--multiline;
CREATE PROCEDURE get_third_party_apk_names (IN apkIDs TEXT)
BEGIN
    SELECT name FROM third_party_apks WHERE apk_id IN (apkIDs);
END