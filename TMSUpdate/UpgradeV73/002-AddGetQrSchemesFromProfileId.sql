--multiline
CREATE PROCEDURE `get_qr_schemes_from_profile_id`(IN profileId int)
BEGIN
    SELECT
        scheme_name
    FROM schemes
    WHERE scheme_name IN (
        SELECT UPPER(dg.name)
        FROM profile_data_group
                 INNER JOIN data_group dg ON
                profile_data_group.data_group_id = dg.data_group_id
        WHERE profile_id = profileId
    );
END