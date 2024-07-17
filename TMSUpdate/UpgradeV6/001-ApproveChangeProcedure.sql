--multiline;
CREATE PROCEDURE approve_change(IN profile_data_id int)
BEGIN
	UPDATE profile_data pd SET pd.approved = 1
	  WHERE pd.profile_data_id = profile_data_id;
END;