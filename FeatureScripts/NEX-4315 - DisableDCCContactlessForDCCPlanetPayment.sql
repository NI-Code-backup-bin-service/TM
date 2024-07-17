--multiline
UPDATE profile_data
SET
    profile_data.datavalue = 'false',
    profile_data.updated_at = NOW(),
    profile_data.updated_by = 'system'
WHERE profile_data.profile_data_id IN
    (SELECT
           DcContactless.profile_data_id
    FROM profile p
    INNER JOIN profile_type pt ON
        p.profile_type_id = pt.profile_type_id
    LEFT JOIN data_element DcProvDesc ON
        DcProvDesc.name ='dccProvider'
    LEFT JOIN data_element DcContactlessDesc ON
        DcContactlessDesc.name = 'dccContactless'
    LEFT JOIN (SELECT * FROM profile_data) DcProv ON
        DcProv.data_element_id = DcProvDesc.data_element_id
        AND
        DcProv.profile_id = p.profile_id
    LEFT JOIN (SELECT * FROM profile_data) DcContactless ON
        DcContactless.data_element_id = DcContactlessDesc.data_element_id
        AND
        DcContactless.profile_id = p.profile_id
    WHERE
        pt.name = 'site' #only for sites
        AND
        DcProv.datavalue = 'pp' #DCC Provider is PlanetPayments
        AND
        DcContactless.datavalue = 'true'); #DCC contactless is enabled