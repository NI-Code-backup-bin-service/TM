UPDATE data_element de SET export_display_index = '54' WHERE de.name ='feeEnabled' AND de.data_group_id=(SELECT dg.data_group_id FROM data_group dg WHERE dg.name='paymentFees');
UPDATE data_element de SET export_display_index = '55' WHERE de.name ='printFee' AND de.data_group_id=(SELECT dg.data_group_id FROM data_group dg WHERE dg.name='paymentFees');
UPDATE data_element de SET export_display_index = '56' WHERE de.name ='paymentServicesConfigs' AND de.data_group_id=(SELECT dg.data_group_id FROM data_group dg WHERE dg.name='paymentServices');
UPDATE data_element de SET export_display_index = '57' WHERE de.name ='enabled' AND de.data_group_id=(SELECT dg.data_group_id FROM data_group dg WHERE dg.name='digitalReceipts');
UPDATE data_element de SET export_display_index = '58' WHERE de.name ='onScreenReceipt' AND de.data_group_id=(SELECT dg.data_group_id FROM data_group dg WHERE dg.name='digitalReceipts');
UPDATE data_element de SET export_display_index = '59' WHERE de.name ='qrEnabled' AND de.data_group_id=(SELECT dg.data_group_id FROM data_group dg WHERE dg.name='digitalReceipts');
UPDATE data_element de SET export_display_index = '60' WHERE de.name ='smsEnabled' AND de.data_group_id=(SELECT dg.data_group_id FROM data_group dg WHERE dg.name='digitalReceipts');
UPDATE data_element de SET export_display_index = '61' WHERE de.name ='emailEnabled' AND de.data_group_id=(SELECT dg.data_group_id FROM data_group dg WHERE dg.name='digitalReceipts');
UPDATE data_element de SET export_display_index = '62' WHERE de.name ='txnHistoryEnabled' AND de.data_group_id=(SELECT dg.data_group_id FROM data_group dg WHERE dg.name='digitalReceipts');
UPDATE data_element de SET export_display_index = '63' where de.name ='printClearPan' AND de.data_group_id=(SELECT dg.data_group_id FROM data_group dg WHERE dg.name='receipt');