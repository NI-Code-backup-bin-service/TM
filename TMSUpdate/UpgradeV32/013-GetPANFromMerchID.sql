--multiline;
CREATE PROCEDURE `get_pan_perm_from_merch_id`(IN merchID varchar(255))
BEGIN
  SELECT
    pd.datavalue AS 'perm_bool'
  FROM `profile_data` AS pd
    LEFT JOIN `data_element` AS de
      ON de.data_element_id = pd.data_element_id
  WHERE pd.data_element_id = de.data_element_id
    AND de.name = 'allowGetPAN'
    AND pd.profile_id = (SELECT pdn.profile_id
                         FROM `profile_data` AS pdn
                           LEFT JOIN `data_element` AS den
                             ON den.data_element_id = pdn.data_element_id
                         WHERE pdn.datavalue = merchID
                           AND den.name = 'propertyId');
END