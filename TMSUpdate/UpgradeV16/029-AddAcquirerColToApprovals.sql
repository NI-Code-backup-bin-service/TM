ALTER TABLE approvals ADD COLUMN acquirer TEXT;
update approvals set acquirer = 'NI';