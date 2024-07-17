--multiline
CREATE PROCEDURE `loadtest_add_acquirer`(
    IN nameIn varchar(255)
)
BEGIN
    -- Check if the acquirer already exists
    SET @acquirercount = (SELECT COUNT(profile_id) FROM profile WHERE `name` = nameIn);
    IF @acquirercount < 1 THEN
        -- get the typeID for acquirer
        SET @typeID = (SELECT profile_type_id FROM profile_type WHERE `name` = 'acquirer');
        -- save the new acquirer profile
        CALL profile_store(-1, @typeID, nameIn, 1, 'loadtester');
        -- get the profileID of the new acquirer
        SET @profileID = (SELECT profile_id FROM profile WHERE name = nameIn);
        -- save each of the default data elements
        CALL store_profile_data(@profileID, 96, 'false', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 95, 'false', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 81, 'false', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 87, 'false', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 114, 'false', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 74, 'false', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 72, 'epos', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 125, 'false', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 58, 'false', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 82, 'false', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 68, '0', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 107, 'false', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 53, 'false', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 108, 'false', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 22, '[]', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 88, 'Asia/Dubai', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 15, 'false', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 133, 'false', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 101, 'false', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 130, 'false', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 137, '[]', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 86, 'false', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 97, 'AED', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 28, 'false', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 119, '[]', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 83, 'false', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 109, 'false', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 115, 'false', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 59, '[]', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 100, 'false', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 91, 'English', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 71, 'pp', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 92, '[]', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 118, 'false', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 124, 'false', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 102, 'false', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 89, 'false', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 62, '6,4', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 90, 'false', 'loadtester', 1, 0, false);
        CALL store_profile_data(@profileID, 135, 'false', 'loadtester', 1, 0, false);
        -- save each one of the default data groups
        CALL profile_data_group_store(-1, @profileID, 3, -1, 'loadtester');
        CALL profile_data_group_store(-1, @profileID, 2, -1, 'loadtester');
        CALL profile_data_group_store(-1, @profileID, 4, -1, 'loadtester');
        CALL profile_data_group_store(-1, @profileID, 7, -1, 'loadtester');
        CALL profile_data_group_store(-1, @profileID, 5, -1, 'loadtester');
        CALL profile_data_group_store(-1, @profileID, 1, -1, 'loadtester');
        CALL profile_data_group_store(-1, @profileID, 12, -1, 'loadtester');
        CALL profile_data_group_store(-1, @profileID, 10, -1, 'loadtester');
        CALL profile_data_group_store(-1, @profileID, 13, -1, 'loadtester');
        CALL profile_data_group_store(-1, @profileID, 6, -1, 'loadtester');
        CALL profile_data_group_store(-1, @profileID, 11, -1, 'loadtester');
        CALL profile_data_group_store(-1, @profileID, 9, -1, 'loadtester');
        CALL profile_data_group_store(-1, @profileID, 8, -1, 'loadtester');
        CALL profile_data_group_store(-1, @profileID, 17, -1, 'loadtester');

        -- check if the chain is already set up
        SET @chaincount = (SELECT COUNT(profile_id) FROM profile WHERE `name` = 'loadtesterchain');
        IF @chaincount < 1 THEN
            -- get typeID for chain
            SET @chainTypeID = (SELECT profile_type_id FROM profile_type WHERE `name` = 'chain');
            -- save the loadtester chain profile
            CALL profile_store(-1, @chainTypeID, 'loadtesterchain', 1, 'loadtester');
            -- get the profileID of the new chain
            SET @chainProfileID = (SELECT profile_id FROM profile WHERE name = 'loadtesterchain');
            -- insert into chain profiles
            INSERT INTO chain_profiles(chain_profile_id, acquirer_id) VALUES(@chainProfileID, @profileID);
        ELSE
            -- get the profileID of the new chain
            SET @chainProfileID = (SELECT profile_id FROM profile WHERE name = 'loadtesterchain');
            -- insert into chain profiles
            INSERT INTO chain_profiles(chain_profile_id, acquirer_id) VALUES(@chainProfileID, @profileID);
        END IF;
    END IF;
END