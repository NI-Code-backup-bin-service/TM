--multiline;
create table if not exists contact_us( AcquirerId int(11) PRIMARY KEY,
                                       AcquirerName VARCHAR(280),
                                       AcquirerPrimaryPhone VARCHAR(280),
                                       AcquirerSecondaryPhone VARCHAR(280),
                                       AcquirerEmail VARCHAR(280),
                                       AcquirerAddressLineOne VARCHAR(280),
                                       AcquirerAddressLineTwo VARCHAR(280),
                                       AcquirerAddressLineThree VARCHAR(280),
                                       FurtherInformation VARCHAR(280));
