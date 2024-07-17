--multiline
CREATE PROCEDURE UpdateTIDTemp()
BEGIN
    DECLARE done INT DEFAULT 0;
    DECLARE prof_id INT;
    DECLARE data_json JSON;
    DECLARE cur CURSOR FOR
        SELECT profile_data_id, datavalue
        FROM profile_data
        WHERE data_element_id=(SELECT data_element_id from data_element WHERE name ='paymentServicesConfigs' AND data_group_id=(SELECT dg.data_group_id FROM data_group dg WHERE dg.name='paymentServices'))
        AND JSON_VALID(datavalue) AND LEFT(datavalue, 1) = '[' AND RIGHT(datavalue, 1) = ']';
    DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = 1;
    OPEN cur;
    read_loop: LOOP
        FETCH cur INTO prof_id, data_json;
        IF done = 1 THEN
            LEAVE read_loop;
        END IF;

        SET @i = 0;

        WHILE JSON_EXTRACT(data_json, CONCAT('$[', @i, '].TID')) IS NOT NULL DO
                IF JSON_TYPE(JSON_EXTRACT(data_json, CONCAT('$[', @i, '].TID'))) = 'INTEGER' THEN
                    SET @tidValue = JSON_UNQUOTE(JSON_EXTRACT(data_json, CONCAT('$[', @i, '].TID')));
                    SET @tidValue = LPAD(@tidValue, 8, '0');
                    SET data_json = JSON_REPLACE(data_json, CONCAT('$[', @i, '].TID'), @tidValue);
                END IF;

                SET @i = @i + 1;
            END WHILE;
        UPDATE profile_data SET datavalue = data_json WHERE profile_data_id = prof_id;
    END LOOP;

    CLOSE cur;
END
