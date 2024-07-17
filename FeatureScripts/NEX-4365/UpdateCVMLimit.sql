-- This feature script will only update the CVM Contactless Limit Tag, DF12, under ctlsApplicationConfigShared.
-- Change the first value to what is currently on TMS, e.g. DF1206000000030000 which is 300.00 AED
-- Change the second value to the value required, e.g. DF1206000000050000 which is 500.00 AED
-- Change the profile_id to match the ID of the Site, Chain, or Acquirer you want to update, e.g. profile_id = 2735.
-- You can also use this script to update the Global Default by setting "profile_id = 1".
-- If you want to update all profiles remove the "AND profile_id = 1;" segment and add a ";" after the data_element_id value.
UPDATE profile_data SET datavalue = REPLACE(datavalue, 'DF1206000000030000', 'DF1206000000050000') WHERE data_element_id = 40 AND profile_id = 1;