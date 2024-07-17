--multiline;
CREATE PROCEDURE `get_apk_details`(IN tidID INT)
BEGIN
    SELECT tu.target_package_id, replace(tu.third_party_apk,'[]',NULL) as third_party_apk, p.version
    FROM tid_updates tu
    JOIN package p ON  tu.target_package_id = p.package_id
    WHERE tu.tid_id = tidID AND tu.update_date <= NOW() ORDER BY tu.update_date DESC limit 1;
END;