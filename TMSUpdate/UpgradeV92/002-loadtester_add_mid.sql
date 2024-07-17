--multiline
CREATE PROCEDURE `loadtest_add_mid`(
    IN IdIn varchar(255),
    IN acquirerName varchar(255)
)
BEGIN
    -- Check if the acquirer already exists
    SET @midcount = (SELECT COUNT(profile_id) FROM profile WHERE `name` = IdIn);
    IF @midcount < 1 THEN
        -- get the typeID for acquirer
        SET @typeID = (SELECT profile_type_id FROM profile_type WHERE `name` = 'site');
        -- save the new acquirer profile
        CALL profile_store(-1, @typeID, IdIn, 1, 'loadtester');
        -- get the profileID of the new acquirer
        SET @profileID = (SELECT profile_id FROM profile WHERE name = IdIn);
        -- save the site
        CALL site_store(-1, 1, 'loadtester');
        -- retrieve the siteID created by the site_store proc
        SET @siteID = (SELECT MAX(site_id) FROM site);
        -- store site profile
        CALL site_profiles_store(-1, @siteID, @profileID, 1, 'loadtester');
        -- get the ID for the site's acquirer
        SET @acquirerID = (SELECT profile_id FROM profile WHERE name = acquirerName);
        -- link site to acquirer profile
        CALL site_profiles_store(-1, @siteID, @acquirerID, 1, 'loadtester');
        -- get the ID for the loadtester chain
        SET @chainID = (SELECT profile_id FROM profile WHERE name = 'loadtesterchain');
        -- link site to chain profile
        CALL site_profiles_store(-1, @siteID, @chainID, 1, 'loadtester');
        -- link site to global profile
        CALL site_profiles_store(-1, @siteID, 1, 1, 'loadtester');
        -- set a unique site name
        SET @siteName = CONCAT('loadtest-',IdIn);
        -- save each of the default data elements
        CALL store_profile_data(@profileID, 89, 'false', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 16, '15', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 127, '0', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 61, '', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 117, '', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 100, 'true', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 125, 'false', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 84, '5537', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 12, '10', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 73, '', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 4, 'loadtest street', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 36, '', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 54, '', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 75, '', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 107, 'false', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 110, 'false', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 3, @siteName, 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 38, '', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 78, '', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 5, 'loadtest city', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 8, '', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 86, 'false', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 135, 'false', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 71, 'pp', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 32, '', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 91, 'English', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 111, '', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 118, 'false', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 116, '', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 92, '["English","French","Arabic"]', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 76, '', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 106, '', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 70, '', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 64, '', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 66, '10', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 101, 'true', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 83, 'false', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 136, '11111', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 113, '', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 93, '', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 59, '["void"]', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 41, '', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 29, '00000', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 69, '', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 129, '60', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 124, 'false', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 55, '', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 37, '', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 74, 'true', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 79, '', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 1, IdIn, 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 21, '', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 17, '20', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 85, '85', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 60, '60', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 95, 'false', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 81, 'true', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 77, '', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 114, 'false', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 31, '', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 82, 'true', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 115, 'false', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 80, '', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 112, '', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 128, '1', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 34, '', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 98, '', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 9, '', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 102, 'false', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 96, 'true', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 126, '[ { ""cardName"":""VISA"", ""transactionLimit"":0 }, { ""cardName"":""MASTER"", ""transactionLimit"":0 } ]', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 103, '', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 14, '00:01', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 52, '', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 94, '', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 33, '', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 39, '', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 43, '', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 58, 'true', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 97, 'AED', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 22, '["sale","refund","void","preAuth","gratuitySale","gratuityCompletion","alipay","upi","xls","visaQr","mastercardQr","eppVoid","balanceInquiry","PWCB"]', 'loadtester', 1, 1, false);
        CALL store_profile_data(@profileID, 67, '5', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 28, 'false', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 88, 'Asia/Dubai', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 10, '', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 53, 'true', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 35, '', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 130, 'false', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 13, '4', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 90, 'false', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 105, '', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 99, '', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 87, 'false', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 11, '', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 104, '', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 133, 'false', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 137, '[]', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 40, '', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 119, '[]', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 108, 'false', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 57, '', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 109, 'false', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 72, 'standalone', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 15, 'true', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 68, '0', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 62, '6,4', 'loadtester', 1, 1, false);
        CALL store_profile_data(@profileID, 65, '305', 'loadtester', 1, 1, false);
    END IF;
END