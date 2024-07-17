# script parameters to be changed by user
SET @p_MID = "00000002";
SET @p_provider = "terraPay";
SET @p_menuFile = "mainMenu-v0.011.json";
call enable_terraPay_by_MID(@p_MID, @p_provider, @p_menuFile)