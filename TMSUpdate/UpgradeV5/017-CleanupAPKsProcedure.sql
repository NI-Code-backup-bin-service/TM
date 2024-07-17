--multiline;
create procedure cleanup_apks()
begin
START TRANSACTION;
    DELETE a FROM apk a LEFT OUTER join package_apk ap ON ap.apk_id = a.apk_id WHERE ap.apk_id IS NULL; 
COMMIT;
end;
