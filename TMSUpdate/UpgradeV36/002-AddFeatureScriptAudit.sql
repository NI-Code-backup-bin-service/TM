--multiline
CREATE PROCEDURE `add_feature_script_audit`(
IN script_name varchar(255),
IN s_date date,
In s_time time,
IN status int,
IN failure_reason varchar(250)
)
BEGIN
    insert into feature_script_audits(s_name, s_date, s_time, outcome, fail_reason) values(script_name, s_date, s_time, status, failure_reason);
END