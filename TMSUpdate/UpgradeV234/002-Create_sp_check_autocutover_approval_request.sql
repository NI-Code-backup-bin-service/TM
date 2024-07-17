--multiline;
CREATE PROCEDURE `get_if_autocutover_change_approval_requtest`(IN siteIdVal int, IN profileIDVal int, IN overriddenVal int)

BEGIN

DECLARE autoCutOverElementId,timeElementId,countChainAcqNoSite,countChainAcqSite,acqAutoCutOver,ifSiteId,checkSiteTimeRange,countSitesTimeRangeNoAuto,chainTimeRangeNoAuto,ifAcqNoChain,AcqCount INT;
DECLARE checkAutoFromChain,checkAutoFromAcq,ifAcqAutoEnable varchar(10);

SET @autoCutOverElementId = (SELECT data_element_id FROM data_element WHERE data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'endOfDay') AND name = 'autoCutOver');
SET @timeElementId = (SELECT data_element_id FROM data_element WHERE name = 'time' AND data_group_id IN (SELECT data_group_id FROM data_group WHERE name = 'endOfDay'));

IF (siteIdVal = -1) THEN
    SET @ifAcqNoChain = (SELECT COUNT(DISTINCT(chain_profile_id)) FROM chain_profiles WHERE acquirer_id = profileIDVal);
    IF (@ifAcqNoChain = 0) THEN
        SELECT 'false' as elementVal;
    ELSE
        SET @ifAcqNoChain = (SELECT COUNT(DISTINCT(chain_profile_id)) FROM chain_profiles WHERE acquirer_id = profileIDVal);
        SET @acqAutoCutOver = (SELECT COUNT(datavalue) FROM profile_data WHERE profile_id IN
                              (SELECT chain_profile_id FROM chain_profiles WHERE acquirer_id = profileIDVal) AND data_element_id = @autoCutOverElementId AND datavalue = 'true');
        IF (@ifAcqNoChain != @acqAutoCutOver) THEN
            SET @chainTimeRangeNoAuto = (SELECT COUNT(DISTINCT(profile_id)) FROM approvals WHERE data_element_id IN
                                        (@timeElementId)  AND LENGTH(new_value) = 13 AND approved = 0 AND profile_id
                                        IN (SELECT DISTINCT(chain_profile_id) from chain_profiles where acquirer_id = profileIDVal) AND
                                        profile_id NOT IN (select DISTINCT(profile_id) from profile_data where profile_id IN (
                                        SELECT DISTINCT(chain_profile_id) from chain_profiles where acquirer_id = profileIDVal) AND  data_element_id IN
                                        (@autoCutOverElementId)  AND datavalue = 'true'));
            IF (@chainTimeRangeNoAuto != 0) THEN
                SELECT 'true' as elementVal;
            ELSE
                SET @countChainAcqNoSite = (SELECT COUNT(site_id) FROM site_profiles WHERE profile_id IN (SELECT DISTINCT(chain_profile_id) FROM chain_profiles WHERE acquirer_id = profileIDVal));
                IF (@countChainAcqNoSite = 0 ) THEN
                    SELECT 'false' as elementVal;
                ELSEIF (@countChainAcqNoSite != 0 ) THEN
                    SET @countChainAcqSiteTime = (SELECT COUNT(DISTINCT(profile_id)) FROM approvals WHERE approved = 0 AND profile_id IN (
                                              SELECT DISTINCT(pd.profile_id) from profile_data pd JOIN site_profiles sp ON ((sp.profile_id = pd.profile_id)) WHERE sp.site_id IN (
                                              SELECT DISTINCT(site_id) FROM site_profiles WHERE profile_id IN (
                                              SELECT DISTINCT(chain_profile_id) FROM chain_profiles WHERE acquirer_id = profileIDVal AND chain_profile_id NOT IN (
                                              SELECT DISTINCT(profile_Id) from profile_data where data_element_id = (SELECT data_element_id FROM data_element WHERE data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'endOfDay') AND name = 'autoCutOver') AND datavalue ='true')))) AND data_element_id = (SELECT data_element_id FROM data_element WHERE name = 'time' AND data_group_id IN (SELECT data_group_id FROM data_group WHERE name = 'endOfDay')) AND LENGTH(new_value) = 13 AND profile_id NOT IN (
                                              SELECT DISTINCT(pd.profile_id) from profile_data pd JOIN site_profiles sp ON ((sp.profile_id = pd.profile_id)) WHERE sp.site_id IN (
                                              SELECT DISTINCT(site_id) FROM site_profiles WHERE profile_id IN (
                                              SELECT DISTINCT(chain_profile_id) FROM chain_profiles WHERE acquirer_id = profileIDVal AND chain_profile_id NOT IN (
                                              SELECT DISTINCT(profile_Id) from profile_data where data_element_id = (SELECT data_element_id FROM data_element WHERE data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'endOfDay') AND name = 'autoCutOver') AND datavalue ='true'))) AND data_element_id = (SELECT data_element_id FROM data_element WHERE data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'endOfDay') AND name = 'autoCutOver') AND overriden=1));
                    IF (@countChainAcqSiteTime = 0) THEN
                        SELECT 'false' as elementVal;
                    ELSE
                        SELECT 'true' as elementVal;
                    END IF;
                ELSE
                    SELECT 'false' as elementVal;
                END IF;
            END IF;
        ELSE
            SELECT 'false' as elementVal;
        END IF;
    END IF;
ELSEIF (siteIdVal = -2) THEN
    SET @ifSiteId = (SELECT COUNT(DISTINCT(site_id)) FROM site_profiles WHERE profile_id = profileIDVal);
    IF (@ifSiteId = 0) THEN
        IF (overriddenVal = 1) THEN
            SET @ifAcqAutoEnable = (SELECT new_value FROM approvals WHERE approved = 0 AND data_element_id = @autoCutOverElementId AND profile_id = (SELECT acquirer_id FROM chain_profiles WHERE chain_profile_id = profileIDVal));
            IF (@ifAcqAutoEnable = 'false') THEN
                SELECT 'true' as elementVal;
            ELSE
                SELECT 'false' as elementVal;
            END IF;
        ELSE
            SELECT 'false' as elementVal;
        END IF;
    ELSE
        SET @countSitesTimeRangeNoAuto = (SELECT COUNT(DISTINCT(profile_id)) FROM approvals WHERE approved = 0 AND profile_id IN (
                                          SELECT DISTINCT(pd.profile_id) from profile_data pd JOIN site_profiles sp ON (sp.profile_id = pd.profile_id)
                                          WHERE sp.site_id IN (
                                          SELECT DISTINCT(site_id) FROM site_profiles WHERE profile_id IN (profileIDVal))) AND profile_id NOT IN (profileIDVal) AND data_element_id = @timeElementId
                                          AND LENGTH(new_value) = 13 AND
                                          profile_id NOT IN (
                                          SELECT DISTINCT(pd.profile_id) from profile_data pd JOIN site_profiles sp ON (sp.profile_id = pd.profile_id)
                                          WHERE sp.site_id IN (
                                          SELECT DISTINCT(site_id) FROM site_profiles WHERE profile_id IN (profileIDVal)) AND data_element_id =@autoCutOverElementId AND overriden=1));
        SET @AcqCount = (SELECT count(profile_id) FROM approvals WHERE approved = 0 AND new_value = 'false' AND data_element_id = @autoCutOverElementId AND profile_id = (SELECT acquirer_id FROM chain_profiles WHERE chain_profile_id = profileIDVal));
        IF (overriddenVal = 1 AND @countSitesTimeRangeNoAuto != 0) THEN
            IF (@AcqCount != 0) THEN
                SELECT 'true' as elementVal;
            ELSE
                SELECT 'true' as elementVal;
            END IF;
        ELSE
            IF (@countSitesTimeRangeNoAuto !=0 ) THEN
                SELECT 'true' as elementVal;
            ELSE
                SELECT 'false' as elementVal;
            END IF;
        END IF;
    END IF;
ELSE
    IF (overriddenVal = 1) THEN
        SET @checkSiteTimeRange = (SELECT COUNT(profile_id) FROM approvals WHERE approved = 0 AND data_element_id = @timeElementId AND LENGTH(new_value) = 13 AND profile_id = profileIDVal);
        IF (@checkSiteTimeRange = 0) THEN
            SELECT 'false' as elementVal;
        ELSE
            SET @checkAutoFromChain = (SELECT COUNT(DISTINCT(profile_id)) FROM approvals WHERE approved = 0 AND data_element_id = @autoCutOverElementId AND new_value = 'false' AND profile_id IN (SELECT DISTINCT(profile_id) FROM site_profiles WHERE site_id = siteIdVal));
            IF (@checkAutoFromChain != 0) THEN
                SELECT 'true' as elementVal;
            ELSE
                SET @checkAutoFromAcq = (SELECT new_value FROM approvals WHERE approved = 0 AND data_element_id = @autoCutOverElementId AND profile_id = (SELECT DISTINCT(profile_id) FROM site_profiles WHERE site_id = siteIdVal AND profile_id IN (SELECT DISTINCT(profile_id) FROM profile_data_elements WHERE profile_type_priority = 4)));
                IF (@checkAutoFromAcq = 'false') THEN
                    SELECT 'true' as elementVal;
                ELSE
                    SELECT 'false' as elementVal;
                END IF;
            END IF;
        END IF;
    ELSE
        SELECT 'false' as elementVal;
    END IF;
END IF;
END;