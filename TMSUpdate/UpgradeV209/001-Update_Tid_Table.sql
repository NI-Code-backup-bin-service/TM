ALTER TABLE tid ADD COLUMN coordinates VARCHAR(50) DEFAULT NULL;
ALTER TABLE tid ADD COLUMN accuracy VARCHAR(50) DEFAULT NULL;
ALTER TABLE tid ADD COLUMN last_coordinate_time DATETIME NULL;