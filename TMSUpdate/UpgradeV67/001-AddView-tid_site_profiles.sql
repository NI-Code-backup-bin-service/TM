--multiline
CREATE VIEW tid_site_profiles AS
(
SELECT *
FROM (
         SELECT ts.tid_profile_id profile_id,
                ts.tid_id,
                ts.site_id
         FROM tid_site ts
                  LEFT JOIN site_profiles sp ON
                 ts.site_id = sp.site_id
         UNION
         SELECT sp.profile_id,
                ts.tid_id,
                sp.site_id
         FROM tid_site ts
                  LEFT JOIN site_profiles sp ON
                 ts.site_id = sp.site_id
     ) DATA
    );