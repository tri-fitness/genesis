
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE `SUBSCRIPTION` (
	`ID`			INTEGER			NOT NULL,
	`ACCOUNT_UUID`	VARCHAR(36)		NOT NULL,
	`STARTED_AT`	DATETIME		NOT NULL,
	`ENDED_AT`		DATETIME		NULL,
	`TYPE`			VARCHAR(128)	NOT NULL,
	`STATUS`		VARCHAR(128)	NOT NULL,
	`INTERVAL`		VARCHAR(128)	NOT NULL,
	`CANCELLED_AT`	DATETIME		NULL,

	PRIMARY KEY (`ID`),
	FOREIGN KEY(`ACCOUNT_UUID`)
		REFERENCES `ACCOUNT` (`UUID`)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `SUBSCRIPTION`;

