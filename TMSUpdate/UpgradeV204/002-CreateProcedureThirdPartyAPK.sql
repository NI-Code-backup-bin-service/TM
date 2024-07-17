--multiline;
CREATE PROCEDURE get_third_party_apks(IN dataValue varchar(45))
BEGIN
    select apk_id, `name` from third_party_apks where name  LIKE CONCAT('%', dataValue , '%');
END