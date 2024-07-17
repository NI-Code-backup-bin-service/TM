-- --multiline;
CREATE PROCEDURE `get_site_tids`(IN site_id int)
BEGIN
SELECT t.tid_id,
       t.serial,
       t.PIN,
       t.ExpiryDate,
       t.reset_pin,
       t.reset_pin_expiry_date,
       t.ActivationDate,
       pd.datavalue as 'merchant_id',
        (SELECT COUNT(*) FROM tid_user_override tuo WHERE tuo.tid_id = t.tid_id) as 'userOverrides',
        (SELECT COUNT(*) FROM velocity_limits vl WHERE vl.tid_id = t.tid_id) as 'fraudOverride',
        IFNULL(ts.tid_profile_id,0) as 'tidProfileID'
FROM tid t
         LEFT JOIN tid_site ts ON ts.tid_id = t.tid_id
         LEFT JOIN site s on s.site_id = ts.site_id
         LEFT JOIN (site_profiles tp2
    join profile p2 on p2.profile_id = tp2.profile_id
    join profile_type pt2 on pt2.profile_type_id = p2.profile_type_id and pt2.priority = 2) on tp2.site_id = s.site_id
         LEFT JOIN profile_data pd ON pd.profile_id = p2.profile_id AND pd.data_element_id = 1
    AND pd.version = (SELECT MAX(d.version) FROM profile_data d WHERE d.data_element_id = 1 AND d.profile_id = p2.profile_id AND d.approved = 1)
WHERE ts.site_id = site_id;
END;

