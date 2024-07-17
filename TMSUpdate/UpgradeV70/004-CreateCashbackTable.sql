--multiline;
CREATE TABLE `cashback` (
  `profile_id` INT NOT NULL,
  `cashback_data` TEXT NOT NULL,
  UNIQUE INDEX `profile_id_UNIQUE` (`profile_id` ASC),
  CONSTRAINT `profile_id`
    FOREIGN KEY (`profile_id`)
    REFERENCES `profile` (`profile_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8;