ALTER TABLE bulk_approvals ADD COLUMN filetype varchar(45) NOT NULL AFTER filename;
UPDATE bulk_approvals SET filetype = 'TerminalFlagging';