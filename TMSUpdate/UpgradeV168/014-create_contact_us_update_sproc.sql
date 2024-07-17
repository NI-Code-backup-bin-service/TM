--multiline;
CREATE procedure set_contact_us_data(IN AcquirerId int,
                                     IN AcquirerName text,
                                     IN AcquirerPrimaryPhone text,
                                     IN AcquirerSecondaryPhone text,
                                     IN AcquirerEmail text,
                                     IN AcquirerAddressLineOne text,
                                     IN AcquirerAddressLineTwo text,
                                     IN AcquirerAddressLineThree text,
                                     IN FurtherInformation text)
BEGIN
    INSERT INTO contact_us value (AcquirerId,
                                  AcquirerName,
                                  AcquirerPrimaryPhone,
                                  AcquirerSecondaryPhone,
                                  AcquirerEmail,
                                  AcquirerAddressLineOne,
                                  AcquirerAddressLineTwo,
                                  AcquirerAddressLineThree,
                                  FurtherInformation)
    ON DUPLICATE KEY UPDATE AcquirerId               = AcquirerId,
                            AcquirerName             = AcquirerName,
                            AcquirerPrimaryPhone     = AcquirerPrimaryPhone,
                            AcquirerSecondaryPhone   = AcquirerSecondaryPhone,
                            AcquirerEmail            = AcquirerEmail,
                            AcquirerAddressLineOne   = AcquirerAddressLineOne,
                            AcquirerAddressLineTwo   = AcquirerAddressLineTwo,
                            AcquirerAddressLineThree = AcquirerAddressLineThree,
                            FurtherInformation       = FurtherInformation;
END