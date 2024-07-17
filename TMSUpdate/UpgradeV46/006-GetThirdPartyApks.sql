--multiline;
CREATE PROCEDURE get_third_party_apks()
BEGIN
    SELECT
        apk_id,
        `name`
    FROM
        third_party_apks;
END