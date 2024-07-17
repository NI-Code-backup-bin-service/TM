--multiline
CREATE PROCEDURE `get_tid_acquirer`(
    in tid text
)
BEGIN
    SELECT
        p.name
    from profile p
         left join tid_site ts on ts.tid_id = tid
         left join site_profiles sp on sp.site_id = ts.site_id
         left join profile_type pt on pt.profile_type_id = p.profile_type_id
    where
        pt.name = "acquirer" and p.profile_id = sp.profile_id limit 1;
END