--multiline;
CREATE PROCEDURE set_tid_target_package(IN tid int, target_package int)
BEGIN
	UPDATE tid t SET t.target_package_id = target_package
    WHERE t.tid_id = tid;
END;