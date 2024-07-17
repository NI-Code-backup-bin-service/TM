--multiline
CREATE PROCEDURE `removeTidOverride` (IN overrideId INT)
BEGIN
	DELETE FROM profile_data_group WHERE profile_id = overrideId;
    DELETE FROM profile_data WHERE profile_id = overrideId;
    UPDATE tid_site SET tid_profile_id = NULL, updated_at = NOW() WHERE tid_profile_id = overrideId;
    DELETE FROM approvals WHERE profile_id = overrideId AND approved = 0;
END
