--multiline
CREATE PROCEDURE AddColumn(IN P_TableName TEXT, IN P_ColumnName TEXT, IN P_ColumnType TEXT, IN P_Args TEXT)
BEGIN
    IF NOT EXISTS(SELECT NULL
                   FROM INFORMATION_SCHEMA.COLUMNS
                   WHERE table_name = P_TableName
                     AND table_schema = database()
                     AND column_name = P_ColumnName)  THEN
        SET @preparedStatement = CONCAT('ALTER TABLE ', P_TableName, ' ADD ', P_ColumnName, ' ', P_ColumnType, ' ', P_Args);
        PREPARE addColumnStatement FROM @preparedStatement;
        EXECUTE addColumnStatement;
        DEALLOCATE PREPARE addColumnStatement;
    END IF;
END;