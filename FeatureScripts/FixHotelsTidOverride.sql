ALTER TABLE `tid_user_override` DROP INDEX `UQ_Username`, ADD UNIQUE INDEX `UQ_Username` (`Username` ASC, `tid_id` ASC);
ALTER TABLE `tid_user_override` DROP INDEX `UQ_PIN` , ADD UNIQUE INDEX `UQ_PIN` (`PIN` ASC, `tid_id` ASC);
