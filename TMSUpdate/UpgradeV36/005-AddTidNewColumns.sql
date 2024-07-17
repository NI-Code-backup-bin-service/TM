--multiline;
ALTER TABLE `tid`
    ADD COLUMN `last_checked_time` BIGINT NULL AFTER `Presence`,
    ADD COLUMN `confirmed_time` BIGINT NULL AFTER `last_checked_time`;