--multiline;
CREATE PROCEDURE get_tid_datavalue_from_site(IN siteId int, IN dataElementName text, IN dataGroupName text)
BEGIN
    SELECT pd.datavalue, pd.is_encrypted
    FROM profile_data pd
    JOIN tid_site ts
        ON ts.tid_profile_id=pd.profile_id
    JOIN data_element de
        ON de.data_element_id=pd.data_element_id
    JOIN data_group dg
        ON dg.data_group_id=de.data_group_id
    WHERE ts.site_id=siteId AND de.name=dataElementName AND dg.name=dataGroupName;
END;