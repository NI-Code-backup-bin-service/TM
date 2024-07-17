--multiline;
CREATE PROCEDURE get_site_profile_payment_services(IN siteId INT)
BEGIN
    DECLARE serviceGroupId INT;
    DECLARE serviceGroupName TEXT;
    DECLARE dataGroupId INT;
    DECLARE profileId INT;
    SELECT dg.data_group_id into dataGroupId FROM data_group dg WHERE dg.name = 'paymentServices';
    SELECT sf.profile_id
    into profileId
    FROM site_profiles sf
             JOIN profile p ON sf.profile_id = p.profile_id
    WHERE sf.site_id = siteId
      AND p.profile_type_id = (SELECT pt.profile_type_id FROM profile_type pt WHERE pt.name = 'site');
    IF dataGroupId IS NOT NULL THEN
        SELECT pd.datavalue
        into serviceGroupName
        FROM profile_data pd
        WHERE pd.profile_id = profileId
          AND pd.data_element_id = (SELECT de.data_element_id
                                    FROM data_element de
                                    WHERE de.name = 'paymentServiceGroup' AND de.data_group_id = dataGroupId);
        SELECT psg.group_id into serviceGroupId FROM payment_service_group psg WHERE psg.name = serviceGroupName;
        SELECT ps.service_id, ps.name FROM payment_service ps WHERE ps.group_id=serviceGroupId;
    ELSE
        SELECT NULL;
    END IF;
END;