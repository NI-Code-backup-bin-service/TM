--multiline;
CREATE PROCEDURE `get_thirdparty_partialpackagename_datavalue`(IN siteId int,IN tidId text)
BEGIN
SET @data_value=
(
       SELECT datavalue
       FROM   profile_data pd
       JOIN   tid_site td
       ON     pd.profile_id = td.tid_profile_id
       WHERE  td.tid_id = tidId
       AND    td.site_id = siteId
       AND    data_element_id =
              (
                     SELECT data_element_id
                     FROM   data_element
                     WHERE  name = 'partialPackageName')
			 );
	IF @data_value IS NULL then
	SELECT pd.datavalue
	FROM   profile_data pd
	JOIN   site_profiles sp
	ON     sp.profile_id = pd.profile_id
	JOIN   profile p
	ON     sp.profile_id = p.profile_id
       JOIN profile_type pt 
	ON p.profile_type_id = pt.profile_type_id
	WHERE  
       site_id = siteId
       AND
       data_element_id =
       (
              SELECT data_element_id
              FROM   data_element
              WHERE  name = 'partialPackageName' 
		) order by pt.priority limit 1;
	ELSE 
	SELECT @data_value limit 1;
	END IF;
END;