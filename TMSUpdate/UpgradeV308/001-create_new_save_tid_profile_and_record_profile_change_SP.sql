--multiline;
CREATE PROCEDURE `create_tid_override_and_save_profile_change`(
    IN profileTypeId INT,
    IN updated_by varchar(255),
    IN change_type INT,
    IN In_dataValue varchar(255),
    IN tidId varchar(255),
	IN tidInt INT,
    IN approved INT,
    IN siteId INT
)
BEGIN

Declare profileId INT;

     #insert and get profile
      INSERT INTO profile(
      profile_type_id,
      name,
      version,
      updated_at,
      updated_by,
      created_at,
      created_by
    ) VALUES (
      profileTypeId,
      tidId,
      1,
      CURRENT_TIMESTAMP,
      updated_by,
      CURRENT_TIMESTAMP,
      updated_by
    );      
    set profileId=LAST_INSERT_ID();
 

    #Add tid profile link
	update tid_site set tid_profile_id = profileId where tid_id = tidInt and site_id = siteId;

	#get active/enabled datagroups and add to profile data group
	INSERT INTO profile_data_group (profile_id, data_group_id, version, updated_at, updated_by, created_at, created_by)
	(SELECT profileId,data_group_id, 1, NOW(), updated_by, NOW(), updated_by
		FROM data_group
		WHERE name IN (SELECT distinct(dg.name)
	FROM site_profiles sp
	INNER JOIN site_profiles sp2 ON sp.site_id = sp2.site_id
    INNER JOIN profile_data_group pdg ON pdg.profile_id=sp2.profile_id
	INNER JOIN data_group dg ON pdg.data_group_id=dg.data_group_id
    INNER JOIN data_element de ON dg.data_group_id=de.data_group_id
	WHERE sp.profile_id = (SELECT sp.profile_Id FROM site_profiles sp 
		LEFT JOIN profile p ON p.profile_id = sp.profile_id 
		LEFT JOIN profile_type pt ON pt.profile_type_id = p.profile_type_id
       WHERE sp.site_id = siteId  AND pt.priority = 2)
	AND sp2.profile_id != 1 AND de.tid_overridable=1));
	

    #add data elements details 
    INSERT INTO profile_data
    (profile_id, data_element_id, datavalue, version, updated_by, created_by, approved, overriden, is_encrypted)
    (WITH profileData AS (
        SELECT pd.data_element_id, pd.datavalue, pd.is_encrypted,
               ROW_NUMBER() OVER(PARTITION BY pd.data_element_id ORDER BY pt.priority) AS RowNum
        FROM profile_data pd
        JOIN data_element de
            ON de.data_element_id=pd.data_element_id
        JOIN site_profiles sp
            ON pd.profile_id = sp.profile_id
        JOIN profile p
            ON p.profile_id = sp.profile_id
        JOIN profile_type pt
            ON p.profile_type_id = pt.profile_type_id
        WHERE sp.site_id = siteId
        and de.tid_overridable=1
        and de.data_group_id IN (SELECT DISTINCT(de.data_group_id)
            FROM profile_data pd
                JOIN site_profiles sp
                    ON pd.profile_id= sp.profile_id
                JOIN profile p
                    ON p.profile_id=sp.profile_id
                JOIN profile_type pt
                    ON pt.profile_type_id=p.profile_type_id
                JOIN data_element de
                    ON de.data_element_id=pd.data_element_id
            WHERE sp.site_id=siteId and pt.name='site')
    )
    SELECT profileId, data_element_id, datavalue, 0, updated_by, updated_by, 1, 0, is_encrypted
    FROM profileData
    WHERE RowNum = 1);
	
    #record the changes
    SET @profileType = (SELECT profile_type_id FROM `profile` p WHERE p.profile_id = profileId limit 1);
    IF @profileType = 4 THEN
        SET @siteId = (SELECT site_id FROM site_profiles sp WHERE sp.profile_id = profileId limit 1);
        SET @acquirer = (SELECT DISTINCT p4.name FROM site_profiles tp4
            JOIN profile p4 ON p4.profile_id = tp4.profile_id
            JOIN profile_type pt4 ON pt4.profile_type_id = p4.profile_type_id AND pt4.priority = 4 WHERE tp4.site_id = @siteId);
    ELSEIF @profileType = 5 THEN
        SET @acquirer = (SELECT DISTINCT p4.name FROM profile p
                                                          JOIN tid_site ts ON ts.tid_profile_id = profileId
                                                          JOIN site t ON t.site_id = ts.site_id
                                                          JOIN (site_profiles tp4
            JOIN profile p4 ON p4.profile_id = tp4.profile_id
            JOIN profile_type pt4 ON pt4.profile_type_id = p4.profile_type_id AND pt4.priority = 4)
                                                               ON tp4.site_id = t.site_id
                         WHERE p.profile_type_id = (SELECT profile_type_id FROM profile_type WHERE profile_type.name = "tid"));
	END IF;
    SET @createChangeType = (SELECT approval_type_id FROM approval_type WHERE approval_type_name = "Create" LIMIT 1);

    SET @unapprovedRequestCount = (SELECT COUNT(*) FROM approvals a WHERE
            a.profile_id = profile_id AND a.change_type = change_type AND a.new_value = In_dataValue AND a.approved_at IS NULL AND a.approved = approved AND a.tid_id = tidId);
    IF  @unapprovedRequestCount < 1 OR change_type = @createChangeType THEN
        INSERT INTO approvals(profile_id, data_element_id, change_type, current_value, new_value, created_at, approved, created_by, approved_by, approved_at, tid_id, acquirer)
        VALUES (profileId, 1, change_type, current_value, In_dataValue, NOW(), approved, updated_by, updated_by, CASE WHEN approved > 0 THEN NOW() ELSE NULL END, tidId, @acquirer);
	END IF;
    END