--multiline;
INSERT INTO tid_updates(tid_update_id, tid_id, target_package_id, update_date) 
SELECT 0, tid_id, target_package_id, update_date
from tid where update_date is not null;