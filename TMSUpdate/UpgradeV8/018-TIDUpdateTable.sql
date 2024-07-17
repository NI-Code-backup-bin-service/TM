--multiline
create table if not exists tid_updates (
tid_update_id int(11) NOT NULL,
tid_id int(11) NOT NULL,
target_package_id int(11) NOT NULL,
update_date DATETIME NOT NULL
);