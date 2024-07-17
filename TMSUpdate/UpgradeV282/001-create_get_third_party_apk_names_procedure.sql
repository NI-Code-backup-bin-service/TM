--multiline;
CREATE PROCEDURE get_third_party_apk_names (IN apkIDs TEXT)
BEGIN
    SELECT name FROM third_party_apks WHERE FIND_IN_SET(apk_id, apkIDs);
END