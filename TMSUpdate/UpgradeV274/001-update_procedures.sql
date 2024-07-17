--multiline;
CREATE PROCEDURE get_element_value_v2(IN profile_id int, IN element_id int)
    BEGIN
        declare site_id, chain_id, acquirer_id int default 0;
        select ifnull(sp.site_id, 0) into site_id from profile p inner join profile_type pt on pt.profile_type_id=p.profile_type_id inner join site_profiles sp on sp.profile_id=p.profile_id where p.profile_id=profile_id and pt.name='site';
        IF (site_id > 0)
            THEN
        select p.profile_id into acquirer_id from site_profiles sp inner join profile p on p.profile_id=sp.profile_id inner join profile_type pt ON pt.profile_type_id=p.profile_type_id where sp.site_id=site_id and pt.name='acquirer';
        select p.profile_id into chain_id from site_profiles sp inner join profile p on p.profile_id=sp.profile_id inner join profile_type pt ON pt.profile_type_id=p.profile_type_id where sp.site_id=site_id and pt.name='chain';
        with
            acquirer_p as (
                SELECT distinct sd.datavalue, e.is_encrypted, e.is_password
                FROM site_data_elements e
                INNER JOIN data_element de ON de.data_element_id = e.data_element_id
                INNER JOIN site_data sd ON sd.site_id = e.site_id AND sd.data_element_id = e.data_element_id
                WHERE e.site_id = site_id AND e.data_element_id =element_id AND e.profile_id=profile_id AND sd.level="acquirer" limit 1
            ),
            chain_p as (
                SELECT pd.datavalue, pd.is_encrypted, de.is_password
                FROM profile_data pd inner join data_element de on de.data_element_id=pd.data_element_id
                WHERE pd.data_element_id = element_id and pd.profile_id = chain_id
                ORDER BY pd.version DESC LIMIT 1
            ),
            site_p as (
                SELECT pd.datavalue, pd.is_encrypted, de.is_password
                FROM profile_data pd inner join data_element de on de.data_element_id=pd.data_element_id
                WHERE pd.data_element_id = element_id AND pd.profile_id = profile_id
                ORDER BY pd.version DESC LIMIT 1
            ),
            results as (
                select datavalue, is_encrypted, is_password from site_p
                UNION
                select datavalue, is_encrypted, is_password from chain_p
                UNION
                select datavalue, is_encrypted, is_password from acquirer_p
           )
           select r.datavalue, IFNULL(r.is_encrypted, 0) as is_encrypted, IFNULL(r.is_password, 0) as isPassword from results as r LIMIT 1;
        ELSE
            SELECT pd.datavalue, IFNULL(pd.is_encrypted,0) as is_encrypted, IFNULL(de.is_password, 0) AS isPassword
            FROM profile_data pd inner join data_element de on de.data_element_id=pd.data_element_id WHERE pd.data_element_id = element_id AND pd.profile_id = profile_id ORDER BY pd.version DESC LIMIT 1;
    END IF;
END;