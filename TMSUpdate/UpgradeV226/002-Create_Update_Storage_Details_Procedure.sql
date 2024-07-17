--multiline
CREATE PROCEDURE `update_storage_details`(IN tid int, IN update_free_internal_storage varchar(50), IN update_total_internal_storage varchar(50))
BEGIN
        UPDATE tid t SET t.free_internal_storage = update_free_internal_storage,
                         t.total_internal_storage = update_total_internal_storage
        WHERE t.tid_id = tid;
END