#To Create NEW Site and store MPAN in Enc Format
--multiline
Update data_element SET is_encrypted=1, is_password=1 where data_group_id in(
SELECT data_group_id from data_group where name in('visaQr','mastercardQr')
) and name='mpan';


