--multiline;
CREATE PROCEDURE get_data_element_value(IN profileId INT,IN dataElementName text, IN dataGroupName VARCHAR(255))
BEGIN
    SELECT pd.datavalue, pd.is_encrypted, de.is_password
    FROM profile_data pd
    JOIN data_element de
        ON de.data_element_id=pd.data_element_id
    JOIN data_group dg
        ON dg.data_group_id=de.data_group_id
    WHERE pd.profile_id=profileId AND de.name=dataElementName AND dg.name=dataGroupName;
END;