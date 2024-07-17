INSERT IGNORE INTO data_group (name, version, updated_at, updated_by, created_at, created_by, displayname_en) VALUES ('meezaQR', 1, NOW(), 'system', NOW(), 'system', 'Meeza QR');
UPDATE data_element SET options = CONCAT(options, '|meezaQR') WHERE name = 'active' AND options not like '%meezaQR%';
UPDATE data_element SET options = CONCAT(options, '|meezaQRSale|meezaQRRTPSale|meezaQRRefund') WHERE name = 'PINRestrictedModules' AND options not like '%meezaQR%';
