--multiline
CREATE PROCEDURE `get_profile_id_from_mid`(
  IN mid TEXT)
BEGIN
SELECT profile_id FROM profile_data WHERE datavalue = mid AND data_element_id = (
      SELECT data_element_id FROM data_element WHERE name = 'merchantNo'
    );
END