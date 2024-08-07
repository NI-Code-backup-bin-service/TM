--multiline
CREATE PROCEDURE `updateMpans`()
BEGIN
	# Temp table to store the data.
	DROP TEMPORARY TABLE IF EXISTS visa_qr_temp;
	CREATE TEMPORARY TABLE visa_qr_temp (merchantNo VARCHAR(12), siteId INT(11), siteProfileId INT(11), currentActiveModulesValue TEXT, mpan VARCHAR(16), categoryCode VARCHAR(255));

	SET @siteProfileTypeId = (SELECT profile_type_id FROM profile_type WHERE `name` = 'site');

	SET @merchantNoDataElementId = (SELECT data_element_id FROM data_element WHERE `name` = 'merchantNo');
	SET @visaQrDataGroupId = (SELECT data_group_id FROM data_group WHERE `name` = 'visaQr');
	SET @mpanDataElementId = (SELECT data_element_id FROM data_element WHERE `name` = 'mpan' AND data_group_id = @visaQrDataGroupId);

	# Insert what we've been given.
	INSERT INTO visa_qr_temp (merchantNo, siteId, siteProfileId)
	VALUES  ('200600002455',NULL,NULL),
			('200600002604',NULL,NULL),
			('200600002927',NULL,NULL),
			('200600004287',NULL,NULL),
			('200600005177',NULL,NULL),
			('200600006290',NULL,NULL),
			('200600006621',NULL,NULL),
			('200600007264',NULL,NULL),
			('200600007280',NULL,NULL),
			('200600007538',NULL,NULL),
			('200600008098',NULL,NULL),
			('200600010367',NULL,NULL),
			('200600011480',NULL,NULL),
			('200600012223',NULL,NULL),
			('200600012272',NULL,NULL),
			('200600012363',NULL,NULL),
			('200600012843',NULL,NULL),
			('200600013304',NULL,NULL),
			('200600014286',NULL,NULL),
			('200600014377',NULL,NULL),
			('200600014443',NULL,NULL),
			('200600014930',NULL,NULL),
			('200600014948',NULL,NULL),
			('200600015325',NULL,NULL),
			('200600015713',NULL,NULL),
			('200600015796',NULL,NULL),
			('200600015911',NULL,NULL),
			('200600016133',NULL,NULL),
			('200600016489',NULL,NULL),
			('200600017198',NULL,NULL),
			('200600017230',NULL,NULL),
			('200600017719',NULL,NULL),
			('200600017735',NULL,NULL),
			('200600017933',NULL,NULL),
			('200600018105',NULL,NULL),
			('200600018758',NULL,NULL),
			('200600018998',NULL,NULL),
			('200600020135',NULL,NULL),
			('200600020606',NULL,NULL),
			('200600021919',NULL,NULL),
			('200600022339',NULL,NULL),
			('200600022941',NULL,NULL),
			('200600023675',NULL,NULL),
			('200600023766',NULL,NULL),
			('200600023881',NULL,NULL),
			('200600024939',NULL,NULL),
			('200600025035',NULL,NULL),
			('200600026777',NULL,NULL),
			('200600026785',NULL,NULL),
			('200600027205',NULL,NULL),
			('200600027270',NULL,NULL),
			('200600027353',NULL,NULL),
			('200600027445',NULL,NULL),
			('200600027494',NULL,NULL),
			('200600027643',NULL,NULL),
			('200600027882',NULL,NULL),
			('200600028633',NULL,NULL),
			('200600028757',NULL,NULL),
			('200600029003',NULL,NULL),
			('200600029359',NULL,NULL),
			('200600029946',NULL,NULL),
			('200600029953',NULL,NULL),
			('200600031058',NULL,NULL),
			('200600031769',NULL,NULL),
			('200600031900',NULL,NULL),
			('200600031991',NULL,NULL),
			('200600032981',NULL,NULL),
			('200600033476',NULL,NULL),
			('200600033609',NULL,NULL),
			('200600033930',NULL,NULL),
			('200600034656',NULL,NULL),
			('200600034714',NULL,NULL),
			('200600034953',NULL,NULL),
			('200600035620',NULL,NULL),
			('200600036545',NULL,NULL),
			('200600036636',NULL,NULL),
			('200600036768',NULL,NULL),
			('200600037063',NULL,NULL),
			('200600037113',NULL,NULL),
			('200600037725',NULL,NULL),
			('200600037832',NULL,NULL),
			('200600037980',NULL,NULL),
			('200600038137',NULL,NULL),
			('200600038194',NULL,NULL),
			('200600038483',NULL,NULL),
			('200600038756',NULL,NULL),
			('200600039218',NULL,NULL),
			('200600039267',NULL,NULL),
			('200600039507',NULL,NULL),
			('200600039523',NULL,NULL),
			('200600039762',NULL,NULL),
			('200600039788',NULL,NULL),
			('200600040588',NULL,NULL),
			('200600041370',NULL,NULL),
			('200600041792',NULL,NULL),
			('200600042006',NULL,NULL),
			('200600042022',NULL,NULL),
			('200600042030',NULL,NULL),
			('200600042923',NULL,NULL),
			('200600043046',NULL,NULL),
			('200600049126',NULL,NULL),
			('200600051023',NULL,NULL),
			('200600051825',NULL,NULL),
			('200600051965',NULL,NULL),
			('200600052401',NULL,NULL),
			('200600052435',NULL,NULL),
			('200600052773',NULL,NULL),
			('200600053185',NULL,NULL),
			('200600053961',NULL,NULL),
			('200600054225',NULL,NULL),
			('200600054282',NULL,NULL),
			('200600054555',NULL,NULL),
			('200600055305',NULL,NULL),
			('200600055701',NULL,NULL),
			('200600055792',NULL,NULL),
			('200600055974',NULL,NULL),
			('200600056113',NULL,NULL),
			('200600056402',NULL,NULL),
			('200600056659',NULL,NULL),
			('200600056931',NULL,NULL),
			('200600057129',NULL,NULL),
			('200600057228',NULL,NULL),
			('200600057483',NULL,NULL),
			('200600057491',NULL,NULL),
			('200600057541',NULL,NULL),
			('200600057681',NULL,NULL),
			('200600058135',NULL,NULL),
			('200600058606',NULL,NULL),
			('200600058705',NULL,NULL),
			('200600059075',NULL,NULL),
			('200600059158',NULL,NULL),
			('200600059786',NULL,NULL),
			('200600060339',NULL,NULL),
			('200600060354',NULL,NULL),
			('200600060636',NULL,NULL),
			('200600060800',NULL,NULL),
			('200600060818',NULL,NULL),
			('200600060842',NULL,NULL),
			('200600061352',NULL,NULL),
			('200600061626',NULL,NULL),
			('200600061758',NULL,NULL),
			('200600062012',NULL,NULL),
			('200600062558',NULL,NULL),
			('200600062780',NULL,NULL),
			('200600063564',NULL,NULL),
			('200600063796',NULL,NULL),
			('200600064158',NULL,NULL),
			('200600065221',NULL,NULL),
			('200600065569',NULL,NULL),
			('200600065882',NULL,NULL),
			('200600066153',NULL,NULL),
			('200600066765',NULL,NULL),
			('200600066823',NULL,NULL),
			('200600066914',NULL,NULL),
			('200600067243',NULL,NULL),
			('200600067300',NULL,NULL),
			('200600067631',NULL,NULL),
			('200600067664',NULL,NULL),
			('200600067698',NULL,NULL),
			('200600067896',NULL,NULL),
			('200600068118',NULL,NULL),
			('200600068456',NULL,NULL),
			('200600068621',NULL,NULL),
			('200600069215',NULL,NULL),
			('200600069371',NULL,NULL),
			('200600069587',NULL,NULL),
			('200600069678',NULL,NULL),
			('200600069744',NULL,NULL),
			('200600069777',NULL,NULL),
			('200600069835',NULL,NULL),
			('200600069850',NULL,NULL),
			('200600069900',NULL,NULL),
			('200600070098',NULL,NULL),
			('200600070122',NULL,NULL),
			('200600070130',NULL,NULL),
			('200600070262',NULL,NULL),
			('200600070296',NULL,NULL),
			('200600070379',NULL,NULL),
			('200600070478',NULL,NULL),
			('200600070684',NULL,NULL),
			('200600070767',NULL,NULL),
			('200600070874',NULL,NULL),
			('200600070882',NULL,NULL),
			('200600070890',NULL,NULL),
			('200600070924',NULL,NULL),
			('200600070981',NULL,NULL),
			('200600071153',NULL,NULL),
			('200600071351',NULL,NULL),
			('200600071377',NULL,NULL),
			('200600071476',NULL,NULL),
			('200600071559',NULL,NULL),
			('200600071567',NULL,NULL),
			('200600071609',NULL,NULL),
			('200600071633',NULL,NULL),
			('200600071765',NULL,NULL),
			('200600071799',NULL,NULL),
			('200600071948',NULL,NULL),
			('200600071997',NULL,NULL),
			('200600072193',NULL,NULL),
			('200600072375',NULL,NULL),
			('200600072458',NULL,NULL),
			('200600072474',NULL,NULL),
			('200600072565',NULL,NULL),
			('200600072722',NULL,NULL),
			('200600072789',NULL,NULL),
			('200600072813',NULL,NULL),
			('200600072870',NULL,NULL),
			('200600072961',NULL,NULL),
			('200600073027',NULL,NULL),
			('200600073068',NULL,NULL),
			('200600073274',NULL,NULL),
			('200600073589',NULL,NULL),
			('200600073977',NULL,NULL),
			('200600074306',NULL,NULL),
			('200600074330',NULL,NULL),
			('200600074512',NULL,NULL),
			('200600074785',NULL,NULL),
			('200600074868',NULL,NULL),
			('200600074983',NULL,NULL),
			('200600075089',NULL,NULL),
			('200600075121',NULL,NULL),
			('200600075220',NULL,NULL),
			('200600075238',NULL,NULL),
			('200600075303',NULL,NULL),
			('200600075501',NULL,NULL),
			('200600075741',NULL,NULL),
			('200600075915',NULL,NULL),
			('200600076079',NULL,NULL),
			('200600076111',NULL,NULL),
			('200600076129',NULL,NULL),
			('200600076186',NULL,NULL),
			('200600076293',NULL,NULL),
			('200600076335',NULL,NULL),
			('200600076475',NULL,NULL),
			('200600076624',NULL,NULL),
			('200600076707',NULL,NULL),
			('200600076731',NULL,NULL),
			('200600076756',NULL,NULL),
			('200600076772',NULL,NULL),
			('200600076871',NULL,NULL),
			('200600076897',NULL,NULL),
			('200600076962',NULL,NULL),
			('200600077002',NULL,NULL),
			('200600077101',NULL,NULL),
			('200600077143',NULL,NULL),
			('200600077150',NULL,NULL),
			('200600077176',NULL,NULL),
			('200600077218',NULL,NULL),
			('200600077341',NULL,NULL),
			('200600077481',NULL,NULL),
			('200600077564',NULL,NULL),
			('200600077606',NULL,NULL),
			('200600077630',NULL,NULL),
			('200600077648',NULL,NULL),
			('200600077739',NULL,NULL),
			('200600077804',NULL,NULL),
			('200600077812',NULL,NULL),
			('200600077887',NULL,NULL),
			('200600077994',NULL,NULL),
			('200600078000',NULL,NULL),
			('200600078083',NULL,NULL),
			('200600078117',NULL,NULL),
			('200600078174',NULL,NULL),
			('200600078182',NULL,NULL),
			('200600078190',NULL,NULL),
			('200600078216',NULL,NULL),
			('200600078240',NULL,NULL),
			('200600078273',NULL,NULL),
			('200600078281',NULL,NULL),
			('200600078489',NULL,NULL),
			('200600078562',NULL,NULL),
			('200600078588',NULL,NULL),
			('200600078638',NULL,NULL),
			('200600078646',NULL,NULL),
			('200600078703',NULL,NULL),
			('200600078729',NULL,NULL),
			('200600078760',NULL,NULL),
			('200600078794',NULL,NULL),
			('200600078877',NULL,NULL),
			('200600078919',NULL,NULL),
			('200600078927',NULL,NULL),
			('200600078943',NULL,NULL),
			('200600079024',NULL,NULL),
			('200600079040',NULL,NULL),
			('200600079099',NULL,NULL),
			('200600079131',NULL,NULL),
			('200600079248',NULL,NULL),
			('200600079321',NULL,NULL),
			('200600079388',NULL,NULL),
			('200600079420',NULL,NULL),
			('200600079479',NULL,NULL),
			('200600079495',NULL,NULL),
			('200600079511',NULL,NULL),
			('200600079545',NULL,NULL),
			('200600079602',NULL,NULL),
			('200600079610',NULL,NULL),
			('200600079636',NULL,NULL),
			('200600079644',NULL,NULL),
			('200600079677',NULL,NULL),
			('200600079685',NULL,NULL),
			('200600079693',NULL,NULL),
			('200600079701',NULL,NULL),
			('200600079727',NULL,NULL),
			('200600079743',NULL,NULL),
			('200600079776',NULL,NULL),
			('200600079875',NULL,NULL),
			('200600080006',NULL,NULL),
			('200600080030',NULL,NULL),
			('200600080055',NULL,NULL),
			('200600080139',NULL,NULL),
			('200600080147',NULL,NULL),
			('200600080170',NULL,NULL),
			('200600080188',NULL,NULL),
			('200600080196',NULL,NULL),
			('200600080212',NULL,NULL),
			('200600080253',NULL,NULL),
			('200600080261',NULL,NULL),
			('200600080279',NULL,NULL),
			('200600080295',NULL,NULL),
			('200600080303',NULL,NULL),
			('200600080311',NULL,NULL),
			('200600080329',NULL,NULL),
			('200600080345',NULL,NULL),
			('200600080410',NULL,NULL);

	# Work out site IDs from merchant number data element entries.
	UPDATE visa_qr_temp
	SET siteId = (SELECT sp.site_id
				  FROM site_profiles sp
				  JOIN `profile` p ON sp.profile_id = p.profile_id AND p.profile_type_id = @siteProfileTypeId
				  JOIN profile_data pd ON p.profile_id = pd.profile_id AND pd.data_element_id = @merchantNoDataElementId
				  WHERE pd.datavalue = merchantNo
				  ORDER BY sp.updated_at DESC
				  LIMIT 1);
				  
	# Work out profile IDs from site IDs.
	UPDATE visa_qr_temp
	SET siteProfileId = (SELECT sp.profile_id
						 FROM site_profiles sp
						 JOIN `profile` p ON sp.profile_id = p.profile_id AND p.profile_type_id = @siteProfileTypeId
						 WHERE site_id = siteId
						 ORDER BY sp.updated_at DESC
						 LIMIT 1);

	# Log and then remove any sites which don't actually exist.
	SELECT * FROM visa_qr_temp WHERE siteId IS NULL OR siteProfileId IS NULL;
	DELETE FROM visa_qr_temp WHERE siteId IS NULL OR siteProfileId IS NULL;

	# Update the timestamp of the MPAN field so that config download will see this as new.
	UPDATE profile_data pd
	SET updated_at = CURRENT_TIMESTAMP,
		updated_by = 'system'
	WHERE pd.profile_id IN (SELECT siteProfileId FROM visa_qr_temp)
	AND data_element_id = @mpanDataElementId;

	DROP TEMPORARY TABLE IF EXISTS visa_qr_temp;
END