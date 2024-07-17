--multiline;
CREATE PROCEDURE `get_keys` (IN key_ident varchar(45)) 
BEGIN 
	SELECT 
		lkey, 
        rkey 
	FROM keystore WHERE identifier = key_ident; 
END; 