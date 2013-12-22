DROP TABLE Users;
DROP TABLE Macs;

CREATE TABLE Macs (
	mac CHAR(17),
	secondsseen INT UNSIGNED, # The total amount, in seconds, that the client has been observed
	timesseen INT UNSIGNED, # The total number of frames captured originating from this MAC
	CONSTRAINT Macs_PK
        PRIMARY KEY(mac),
	CONSTRAINT MacLength # This is silently ignored in mysql/mariadb
        CHECK (length(mac) = 17)
);

CREATE TABLE Users (
	mac CHAR(17), # Foriegn key
	cid VARCHAR(15) NOT NULL,
	devicetype VARCHAR(255) NOT NULL, # E.g. "Phone" or "Laptop"

	CONSTRAINT Users_PK
		PRIMARY KEY(mac, cid),
	CONSTRAINT Users_ref_Macs
		FOREIGN KEY(mac) REFERENCES Macs(mac),
	CONSTRAINT Unique_Mac
		UNIQUE(mac)
);

CREATE OR REPLACE VIEW Statistics AS
    SELECT cid, devicetype, secondsseen, timesseen, Users.mac
    FROM Users LEFT OUTER JOIN Macs
                ON Users.mac = Macs.mac
    ORDER BY secondsseen, cid, timesseen DESC;