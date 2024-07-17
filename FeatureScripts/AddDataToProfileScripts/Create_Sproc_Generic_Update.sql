--multiline;
CREATE PROCEDURE `add_generic_data_to_profile`(IN tid INT, IN mid varchar(255), IN deName varchar(255), IN deGroup varchar(255), IN newvalue varchar(255))
sproc: BEGIN

DECLARE v_data_element_id int;
DECLARE v_profile_id int;
SELECT data_element_id INTO v_data_element_id FROM data_element WHERE name = deName AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = deGroup);

IF tid IS NOT NULL THEN
SELECT tid_profile_id INTO v_profile_id FROM tid_site WHERE tid_id = tid;
ELSEIF mid IS NOT NULL THEN
SELECT profile_id INTO v_profile_id FROM profile_data WHERE datavalue = mid AND data_element_id = (SELECT data_element_id FROM data_element WHERE name = 'merchantNo');
ELSE
	SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'No tid or mid provided';
    LEAVE sproc;
END IF;

IF v_profile_id IS NULL THEN
	SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Invalid profile id';
    LEAVE sproc;
END IF;

IF v_data_element_id IS NULL THEN
	SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Invalid data element';
    LEAVE sproc;
END IF;

INSERT INTO profile_data (
    profile_id, data_element_id, datavalue,
    version, updated_at, updated_by,
    created_at, created_by, approved,
    overriden, is_encrypted
) VALUE (
  v_profile_id, v_data_element_id,
  newvalue, 1, NOW(), 'system', NOW(),
  'system', 1, 0, 0
) ON DUPLICATE KEY
UPDATE
    datavalue = newvalue,
    updated_at = NOW(),
    updated_by = 'system';

END