--multiline
CREATE PROCEDURE `get_site_id_from_mid`(
  IN mid TEXT)
BEGIN
SELECT site_id FROM site_profiles WHERE profile_id = (
    SELECT profile_id FROM profile_data WHERE datavalue = mid
      AND data_element_id = (
        SELECT data_element_id FROM data_element WHERE name = 'merchantNo'
      )
  );
END