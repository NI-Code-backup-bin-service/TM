--multiline;
CREATE PROCEDURE `get_all_packages`()
BEGIN
    SELECT p.package_id, p.version, GROUP_CONCAT(a.name)
    FROM package p
        LEFT JOIN  package_apk pa
            ON pa.package_id = p.package_id
        LEFT JOIN apk a
            ON pa.apk_id = a.apk_id
    GROUP BY p.package_id
    ORDER BY p.package_id DESC;
END;