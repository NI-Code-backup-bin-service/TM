--multiline;
CREATE PROCEDURE `count_matched_values`(IN element_id INT, IN data_value varchar(255), IN profileId INT)
BEGIN

SELECT COUNT(*)
FROM profile_data  pd
WHERE pd.data_element_id = element_id
AND pd.datavalue = data_value
AND pd.profile_id != profileId
AND pd.version = (SELECT MAX(version) 
					FROM profile_data d 
                    WHERE d.data_element_id = element_id 
                    AND d.profile_id = pd.profile_id);

END