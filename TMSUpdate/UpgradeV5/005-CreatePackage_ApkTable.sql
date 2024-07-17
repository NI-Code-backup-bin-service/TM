--multiline
CREATE TABLE if not exists package_apk (
  `package_id` INT NOT NULL,
  `apk_id` INT NOT NULL,
  PRIMARY KEY (`package_id`, `apk_id`),
  INDEX `apk_id_idx` (`apk_id` ASC),
  CONSTRAINT `package_id`
    FOREIGN KEY (`package_id`)
    REFERENCES `package` (`package_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `apk_id`
    FOREIGN KEY (`apk_id`)
    REFERENCES `apk` (`apk_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION);

