--multiline;
CREATE PROCEDURE `get_if_autocutover_disable`(IN siteIdVal int,IN dataElementIdVal int, IN profileIDVal int)

BEGIN
DECLARE siteCount,chainCount INT;

IF (siteIdVal = -1) THEN
    SET @chainCount = (SELECT COUNT(profile_id) FROM profile_data WHERE profile_id IN (SELECT chain_profile_id FROM chain_profiles WHERE acquirer_id = profileIDVal) AND data_element_id = dataElementIdVal AND LENGTH(datavalue) = 13);
    IF (@chainCount != 0) THEN
        SELECT 'true' AS elementval;
    ELSE
        SET @siteCount = (SELECT COUNT(profile_id) FROM site_data WHERE site_id IN (SELECT site_id FROM site_profiles WHERE profile_id IN (SELECT profile_id FROM profile_data WHERE profile_id IN (SELECT chain_profile_id FROM chain_profiles WHERE acquirer_id = profileIDVal) AND data_element_id = dataElementIdVal)) AND level = 'site' AND data_element_id = dataElementIdVal AND LENGTH(datavalue) = 13);
        IF (@siteCount != 0) THEN
            SELECT 'true' AS elementval;
        ELSE
            SELECT 'false' AS elementval;
        END IF;
    END IF;
ELSEIF (siteIdVal = -2) THEN
    SET @siteCount = (SELECT COUNT(profile_id) FROM site_data WHERE site_id IN (SELECT site_id FROM site_profiles WHERE profile_id = profileIDVal) AND LEVEL = 'site' AND data_element_id = dataElementIdVal AND LENGTH(datavalue) = 13);
    IF (@siteCount != 0) THEN
        SELECT 'true' AS elementval;
    ELSE
        SELECT 'false' AS elementval;
    END IF;
ELSE
    SELECT 'false' AS elementval;
END IF;
END;