--multiline;
CREATE PROCEDURE get_qr_schemes_from_site_id(
   IN siteId INT
)
BEGIN
    SELECT distinct s.scheme_id, s.scheme_name FROM site_profiles sp
    INNER JOIN profile_data_group pdg ON sp.profile_id = pdg.profile_id
    INNER JOIN data_group dg ON pdg.data_group_id = dg.data_group_id
    INNER JOIN schemes s ON UPPER(s.scheme_name) = UPPER(dg.name)
    WHERE sp.site_id = siteId AND s.qr_scheme = 1;
END;