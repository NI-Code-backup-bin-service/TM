--multiline
UPDATE profile_data
SET datavalue = '[
   "C00101 9F3303E06048 9F350122 9F4005F000A0A001 9F7E0101 DF811801F0",
   "C00102 9F3303E0F0C8 9F350122 9F4005F000A0A001 9F6604B4804000",
   "C00103 9F3303E0F0C8 9F350122 9F4005F000A0A001",
   "C00104 9F3303E0F0C8 9F350122 9F4005F000A0A001",
   "C0010C 9F3303E0F0C8 9F350122 9F4005F000A0A001"
]'
where profile_id = (SELECT profile_id from profile where name = 'global')
and data_element_id = (SELECT data_element_id from data_element where name = 'ctlsKernelConfigs'
and data_group_id = (SELECT data_group_id from data_group where name = 'emv'));
