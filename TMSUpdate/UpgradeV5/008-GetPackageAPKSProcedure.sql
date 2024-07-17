--multiline;
CREATE PROCEDURE get_package_apks(in package_id int)
BEGIN
	select a.name
    from apk a
    left join package_apk pa on pa.apk_id = a.apk_id
    where pa.package_id = package_id;
END