--multiline
INSERT IGNORE INTO data_element_locations_data_element
(location_id, data_element_id)
SELECT (select location_id from data_element_locations where location_name = 'tid_override'), data_element_id
FROM data_element
where displayname_en IN
      ('PED Language', 'Receipt Languages', 'active', 'EPP Enabled', 'TIPMax %', 'TIP', 'manualEntryEnabled', 'mode',
       'preAuthMax %', 'Reference No.', 'Third Party Package Name', 'workstationNumber', 'addressLine1', 'addressLine2',
       'name', 'available', 'superPIN')
   OR name IN
      ('PED Language', 'Receipt Languages', 'active', 'EPP Enabled', 'TIPMax %', 'TIP', 'manualEntryEnabled', 'mode',
       'preAuthMax %', 'Reference No.', 'Third Party Package Name', 'workstationNumber', 'addressLine1', 'addressLine2',
       'name', 'available', 'superPIN');
