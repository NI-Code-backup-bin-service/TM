DROP INDEX `UQ_Username` ON `tid_user_override`;
DROP INDEX `UQ_PIN` ON `tid_user_override`;
CREATE UNIQUE INDEX `UQ_Username` ON `tid_user_override` (Username, tid_id);
CREATE UNIQUE INDEX `UQ_PIN` ON `tid_user_override` (PIN, tid_id);
