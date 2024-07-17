--multiline;
insert into contact_us(AcquirerId,
                       AcquirerName,
                       AcquirerPrimaryPhone,
                       AcquirerSecondaryPhone,
                       AcquirerEmail,
                       AcquirerAddressLineOne,
                       AcquirerAddressLineTwo,
                       AcquirerAddressLineThree,
                       FurtherInformation) values((SELECT profile_id FROM profile WHERE name = 'NI' AND profile_type_id = 2),'NI', '1111', '2222', 'ni@email.com', 'Address', 'Address2', 'Address3', 'FurtherInfo');