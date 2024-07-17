--multiline
CREATE PROCEDURE `remove_override`(IN siteId INT, IN elementId INT)
BEGIN
	DECLARE profileID INT;
    
   SET @profileId = (SELECT 
							sp.profile_id
						FROM site_profiles sp 
						LEFT JOIN profile p ON p.profile_id = sp.profile_id
						LEFT JOIN profile_type pt ON pt.profile_type_id = p.profile_type_id
                        WHERE sp.site_id = siteId
						ORDER BY pt.priority
						LIMIT 1);
						
	DELETE FROM profile_data WHERE profile_id = @profileID AND data_element_id = elementId AND overriden = 1;
END