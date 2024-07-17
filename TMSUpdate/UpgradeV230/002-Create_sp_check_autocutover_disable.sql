--multiline;
CREATE PROCEDURE `get_if_autocutover_disable`(IN siteIdVal int, IN profileIDVal int, IN overriddenVal int)

BEGIN

DECLARE autoCutOverElementId,timeElementId,countChainAcqNoSite,countChainAcqSite,acqAutoCutOver,ifSiteId,checkAutoFromSite,checkSiteTimeRange,countSitesTimeRangeNoAuto,chainTimeRangeNoAuto,ifAcqNoChain INT;
DECLARE checkAutoFromChain,checkAutoFromAcq,ifAcqAutoEnable varchar(10);

SET @autoCutOverElementId = (SELECT data_element_id FROM data_element WHERE data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'endOfDay') AND name = 'autoCutOver');
SET @timeElementId = (SELECT data_element_id FROM data_element WHERE name = 'time' AND data_group_id IN (SELECT data_group_id FROM data_group WHERE name = 'endOfDay'));

IF (siteIdVal = -1) THEN
    SET @ifAcqNoChain = (SELECT COUNT(DISTINCT(chain_profile_id)) FROM chain_profiles WHERE acquirer_id = profileIDVal);
    IF (@ifAcqNoChain = 0) THEN
        SELECT 'false';
    ELSE
        SET @chainTimeRangeNoAuto = (SELECT COUNT(DISTINCT(profile_id)) FROM profile_data WHERE data_element_id = @timeElementId AND LENGTH(datavalue) = 13 AND profile_id NOT IN (SELECT DISTINCT(profile_id) FROM profile_data WHERE data_element_id = @autoCutOverElementId) AND profile_id IN (SELECT DISTINCT(chain_profile_id) FROM chain_profiles WHERE acquirer_id = profileIDVal));
        IF (@chainTimeRangeNoAuto != 0) THEN
            SELECT 'true';
        ELSE
            SET @countChainAcqNoSite = (SELECT COUNT(DISTINCT(profile_id)) FROM chain_data WHERE source = 'chain ' AND data_element_id = @timeElementId and LENGTH(datavalue) = 13 AND profile_id NOT IN (SELECT DISTINCT(profile_id) FROM chain_data WHERE data_element_id = @autoCutOverElementId AND source = 'chain') AND profile_id IN (SELECT chain_profile_id FROM chain_profiles WHERE acquirer_id = profileIDVal AND chain_profile_id NOT IN (SELECT profile_id FROM site_profiles)));
            SET @countChainAcqSite = (SELECT COUNT(DISTINCT(profile_id)) FROM chain_data WHERE profile_id NOT IN (SELECT DISTINCT(profile_id) FROM chain_data WHERE data_element_id = @autoCutOverElementId and datavalue = 'true' and source = 'chain') AND profile_id IN (SELECT DISTINCT(profile_id) FROM chain_data WHERE profile_id IN (SELECT DISTINCT(profile_id) FROM site_profiles WHERE site_id IN (SELECT DISTINCT(site_id) FROM site_data WHERE site_id IN (SELECT DISTINCT(site_id) FROM site_profiles WHERE profile_id = profileIDVal) AND level = 'site' AND data_element_id = @timeElementId and LENGTH(datavalue) = 13 AND profile_id NOT IN (SELECT DISTINCT(profile_id) FROM site_data WHERE data_element_id = @autoCutOverElementId)))));
            IF (@countChainAcqNoSite = 0 AND @countChainAcqSite = 0) THEN
                SELECT 'false' as elementVal;
            ELSE
                SELECT 'true' as elementVal;
            END IF;
        END IF;
    END IF;
ELSEIF (siteIdVal = -2) THEN
    SET @ifSiteId = (SELECT COUNT(DISTINCT(site_id)) FROM site_profiles WHERE profile_id = profileIDVal);
    IF (@ifSiteId = 0) THEN
        IF (overriddenVal = 1) THEN
            SET @ifAcqAutoEnable = (SELECT datavalue FROM profile_data WHERE data_element_id = @autoCutOverElementId AND profile_id = (SELECT acquirer_id FROM chain_profiles WHERE chain_profile_id = profileIDVal));
            IF (@ifAcqAutoEnable = 'false') THEN
                SELECT 'true';
            ELSE
                SELECT 'false';
            END IF;
        ELSE
            SELECT 'false';
        END IF;
    ELSE
        SET @countSitesTimeRangeNoAuto = (SELECT COUNT(DISTINCT(profile_id)) FROM site_data WHERE site_id IN (SELECT DISTINCT(site_id) FROM site_profiles WHERE profile_id = profileIDVal) AND LEVEL = 'site' AND data_element_id = @timeElementId AND LENGTH(datavalue) = 13 AND profile_id NOT IN (SELECT DISTINCT(profile_id) FROM site_data WHERE data_element_id = @autoCutOverElementId));
        SET @ifAcqAutoEnable = (SELECT datavalue FROM profile_data WHERE data_element_id = @autoCutOverElementId AND profile_id = (SELECT acquirer_id FROM chain_profiles WHERE chain_profile_id = profileIDVal));
        IF (overriddenVal = 1) THEN
            IF (@countSitesTimeRangeNoAuto!=0 AND @ifAcqAutoEnable = 'false') THEN
                SELECT 'true';
            ELSE
                SELECT 'false';
            END IF;
        ELSE
            IF (@countSitesTimeRangeNoAuto!=0) THEN
                SELECT 'true';
            ELSE
                SELECT 'false';
            END IF;
        END IF;
    END IF;
ELSE
    IF (overriddenVal = 1) THEN
        SET @checkSiteTimeRange = (SELECT COUNT(profile_id) FROM profile_data WHERE data_element_id = @timeElementId AND LENGTH(datavalue) = 13 AND profile_id = profileIDVal);
        IF (@checkSiteTimeRange = 0) THEN
            SELECT 'false';
        ELSE
            SET @checkAutoFromChain = (SELECT datavalue FROM chain_data WHERE source = 'chain' AND data_element_id = @autoCutOverElementId AND profile_id IN (SELECT DISTINCT(profile_id) FROM site_profiles WHERE site_id = siteIdVal));
            IF (@checkAutoFromChain = 'true') OR (@checkAutoFromChain = 'false') THEN
                IF (@checkAutoFromChain = 'false') THEN
                    SELECT 'true';
                ELSE
                    SELECT 'false';
                END IF;
            ELSE
                SET @checkAutoFromAcq = (SELECT datavalue FROM profile_data WHERE data_element_id = @autoCutOverElementId AND profile_id = (SELECT DISTINCT(profile_id) FROM site_profiles WHERE site_id = siteIdVal AND profile_id IN (SELECT DISTINCT(profile_id) FROM profile_data_elements WHERE profile_type_priority = 4)));
                IF (@checkAutoFromAcq = 'false') THEN
                    SELECT 'true';
                ELSE
                    SELECT 'false';
                END IF;
            END IF;
        END IF;
    ELSE
        SELECT 'false' AS elementval;
    END IF;
END IF;
END;