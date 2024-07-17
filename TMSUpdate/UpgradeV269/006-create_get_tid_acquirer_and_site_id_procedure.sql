--multiline;
CREATE PROCEDURE `get_tid_acquirer_and_site_id`(IN tidID INT)
BEGIN
    SELECT p.name,  ts.site_id
    FROM profile p
    LEFT JOIN tid_site ts
        ON ts.tid_id = tidID
    LEFT JOIN site_profiles sp
        ON sp.site_id = ts.site_id
    LEFT JOIN profile_type pt
        ON pt.profile_type_id = p.profile_type_id
    WHERE pt.name = "acquirer" AND p.profile_id = sp.profile_id limit 1;
END;