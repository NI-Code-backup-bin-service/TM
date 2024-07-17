alter table data_element add column options text not null
update data_element set options = "sale, refund, void, preAuth, gratuitySale, gratuityCompletion" where data_element_id=22