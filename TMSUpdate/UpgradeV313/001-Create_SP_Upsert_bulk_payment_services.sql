--multiline;
CREATE PROCEDURE `upsert_bulk_payment_services`(In tidId TEXT,In_dataValue LONGTEXT)
BEGIN
Declare PaymentServicesCount Int;
select count(*) into PaymentServicesCount from profile_data as pd 
    	inner join tid_site as ts ON ts.tid_profile_id = pd.profile_id 
		inner join data_element de on pd.data_element_id = de.data_element_id 
		where ts.tid_id = tidId
		and de.name = 'paymentServicesConfigs';
 IF PaymentServicesCount > 0 THEN
 update profile_data as pd
		inner join tid_site as ts ON ts.tid_profile_id = pd.profile_id
		inner join data_element de on pd.data_element_id = de.data_element_id
		set datavalue = In_dataValue
		where ts.tid_id = tidId
		and de.name = 'paymentServicesConfigs';
ELSE
insert into profile_data (profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden, is_encrypted, not_overridable) values
		((select tid_profile_id from tid_site where tid_id = tidId limit 1),
		(select data_element_id from data_element where name = 'paymentServicesConfigs' limit 1),
		In_dataValue, 1, NOW(), 'system', NOW(), 'system', 1, 1, 0, 0);
END IF;
END