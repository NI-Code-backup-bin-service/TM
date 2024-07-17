/*
 * This can only be run if the stored procedure add_generic_data_to_profile has been added to the database first. This
 * can be done by running Create_Sproc_Generic_Update.sql, which is in this directory. This should be used to bulk-add
 * configuration to TIDs or MIDs. This is not a replacement for adding new datagroups or making data elements
 * overridable - you must use existing methods for that.
 *
 *
 * Use NextGen_TMS.add_generic_data_to_profile(tid, mid, dataElement, dataGroup, newValue)
 * If you pass in a TID and a MID, the SPROC will use the TID and ignore the MID. If both TID and MID are null, an error
 * is thrown.
 *
 * This SPROC will insert values as new into the profile if they do not exist, or update the value if it already exists
 * in the profile. This WILL add data to a profile even if datagroups are not enabled or are not overridable so use with
 * caution.
 *
 * The below sections add a value to TID Level Override for ENBD Cash Desk ID. Then it adds values for the 4 configurable
 * fields in the site-level configuration. All fields exist within the ENBD data group. The value to set is the FINAL
 * parameter, which is currently a blank string. This is where you must update the value. The TID and MID are placeholder
 * values. You must update them accordingly.
 */
call NextGen_TMS.add_generic_data_to_profile(88880001, NULL, 'cashDeskId', 'IPP', '');
call NextGen_TMS.add_generic_data_to_profile(NULL, '111122223333', 'bankUserId', 'IPP', '');
call NextGen_TMS.add_generic_data_to_profile(NULL, '111122223333', 'paymentLocationId', 'IPP', '');
call NextGen_TMS.add_generic_data_to_profile(NULL, '111122223333', 'merchantTag', 'IPP', '');
call NextGen_TMS.add_generic_data_to_profile(NULL, '111122223333', 'city', 'store', '');

call NextGen_TMS.add_generic_data_to_profile(88880001, NULL, 'cashDeskId', 'IPP', '');
call NextGen_TMS.add_generic_data_to_profile(NULL, '111122223333', 'bankUserId', 'IPP', '');
call NextGen_TMS.add_generic_data_to_profile(NULL, '111122223333', 'paymentLocationId', 'IPP', '');
call NextGen_TMS.add_generic_data_to_profile(NULL, '111122223333', 'merchantTag', 'IPP', '');
call NextGen_TMS.add_generic_data_to_profile(NULL, '111122223333', 'city', 'store', '');

call NextGen_TMS.add_generic_data_to_profile(88880001, NULL, 'cashDeskId', 'IPP', '');
call NextGen_TMS.add_generic_data_to_profile(NULL, '111122223333', 'bankUserId', 'IPP', '');
call NextGen_TMS.add_generic_data_to_profile(NULL, '111122223333', 'paymentLocationId', 'IPP', '');
call NextGen_TMS.add_generic_data_to_profile(NULL, '111122223333', 'merchantTag', 'IPP', '');
call NextGen_TMS.add_generic_data_to_profile(NULL, '111122223333', 'city', 'store', '');
