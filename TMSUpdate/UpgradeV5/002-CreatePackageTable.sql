--multiline;
create table if not exists package (
    package_id int(11) NOT NULL AUTO_INCREMENT, 
    version varchar(45) unique, 
    PRIMARY KEY(package_id)
);
