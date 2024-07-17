--multiline
UPDATE profile_data
SET datavalue = '[
   "XAC_DF0B0100 C00102 9F0607A0000000031010 DF0C0141 DF0D050010000000 DF0E05DC4004F800 DF0F05DC4000A800 500456534443 9F09020105 9F15020011 9F1D084400800000000000",
   "XAC_DF0B0101 C00102 9F0607A0000000032010 DF0C0141 DF0D050010000000 DF0E05DC4004F800 DF0F05DC4000A800 500456534443 9F09020105 9F15020011 9F1D084400800000000000",
   "XAC_DF0B0102 C00102 9F0607A0000000033010 DF0C0141 DF0D050010000000 DF0E05DC4004F800 DF0F05DC4000A800 500456534443 9F09020105 9F15020011 9F1D084400800000000000",
   "XAC_DF0B0103 C00101 9F0607A0000000041010 DF0C0109 DF0D050000000000 DF0E05FC50BCF800 DF0F05FC50BCA000 500A4D617374657243617264 9F09020002 9F15020011 9F1D084400800000000000 E60F E7039C0120 E808DF0D05FFFFFFFFFF",
   "XAC_DF0B0104 C00101 9F0607A0000000043060 DF0C010B DF0D050000800000 DF0E05FC50BCF800 DF0F05FC50BCA000 50074D61657374726F 9F09020002 9F15020011 9F1D084400800000000000 E60F E7039C0120 E808DF0D05FFFFFFFFFF",
   "XAC_DF0B0105 C0010C 9F0607A0000000651010 DF0C0143 DF0D050050000000 DF0E05FC60ACF800 DF0F05FC6024A800 50034A4342 9F09020200 9F15020011 9F1D084400800000000000",
   "VF_9F0607A0000000031010 9F09020105 DF1105DC4000A800 DF1205DC4004F800 DF13050010000000",
   "VF_9F0607A0000000032010 9F09020105 DF1105DC4000A800 DF1205DC4004F800 DF13050010000000",
   "VF_9F0607A0000000033010 9F09020105 DF1105DC4000A800 DF1205DC4004F800 DF13050010000000",
   "VF_9F0607A0000000041010 9F09020002 DF1105FC50BCA000 DF1205FC50BCA000 DF13050000000000",
   "VF_9F0607A0000000043060 9F09020002 DF1105FC50BCA000 DF1205FC50BCA000 DF13050000800000"
]'
where profile_id = (SELECT profile_id from profile where name = 'global')
and data_element_id = (SELECT data_element_id from data_element where name = 'ctlsApplicationConfigs'
and data_group_id = (SELECT data_group_id from data_group where name = 'multiDeviceEmv'));
