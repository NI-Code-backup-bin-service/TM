--multiline;
CREATE PROCEDURE get_site_id_from_profile_id(
   IN profileId INT
)
BEGIN
    SELECT site_id FROM site_profiles WHERE profile_id = profileId;
END;