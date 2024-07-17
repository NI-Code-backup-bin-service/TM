--multiline;
CREATE PROCEDURE get_third_party_apk_name (
    IN id INT
)
BEGIN
    SELECT
        `name`
    FROM
        third_party_apks
    WHERE
        apk_id = id;
END