--multiline
CREATE PROCEDURE UpdateProfileName()
BEGIN
    DECLARE done INT DEFAULT 0;
    DECLARE profileID INT;
    DECLARE profileName VARCHAR(100);
    DECLARE result CURSOR FOR
        SELECT pd.profile_id ,pd.datavalue from profile_data pd
        LEFT JOIN `profile` p ON p.profile_id = pd.profile_id
        WHERE pd.datavalue <> p.name
        AND p.profile_type_id = 4
        AND data_element_id = (select data_element_id from data_element where name='name');
    DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = 1;
    OPEN result;
        result_loop: LOOP
            FETCH result INTO profileID, profileName;
            IF done = 1 THEN
                LEAVE result_loop;
            END IF;

            UPDATE profile SET name = profileName WHERE profile_id = profileID;
        END LOOP;
    CLOSE result;
END
