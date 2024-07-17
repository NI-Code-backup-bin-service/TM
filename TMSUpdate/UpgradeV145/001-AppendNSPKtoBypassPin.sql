-- update Bypass PIN Drop down to include MIR
UPDATE data_element SET options = 'UNIONPAY|NSPK' WHERE name = "bypassPIN";