--multiline;
CREATE TABLE chain_profiles(
  `chain_id` INT NOT NULL AUTO_INCREMENT,
  `chain_profile_id` INT NOT NULL,
  `acquirer_id` INT NOT NULL,
  PRIMARY KEY (`chain_id`)
);