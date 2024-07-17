# Feature script to add NOL SAM card keys to the database
#
# In order to use this script, please replace the values sam_uid and encrypted_sam_key
# with the actual values to be inserted. 
# This insert statement can then be copied and pasted multiple times for each uid-key
# pair to bulk insert.
#
INSERT INTO sam_card_keys(sam_uid, sam_key) VALUES ('sam_uid', 'encrypted_sam_key');