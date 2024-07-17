--multiline;
CREATE TABLE if NOT EXISTS third_party_apks (
    apk_id INT NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(255) UNIQUE,
    `version` VARCHAR(255),
    PRIMARY KEY(apk_id)
);