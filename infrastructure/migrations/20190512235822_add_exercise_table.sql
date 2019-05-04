
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE `EXERCISE` (
	`UUID`			VARCHAR(128)	NOT NULL,
	`VIDEO_ID`		VARCHAR(128)	NOT NULL,
	`DIFFICULTY`	VARCHAR(128)	NOT NULL,
	`WATCHES`		INTEGER			NOT NULL,
	`DURATION`		INTEGER			NOT NULL,
	`AUTHOR_UUID`	VARCHAR(128)	NOT NULL,

	PRIMARY KEY(`UUID`),
	FOREIGN KEY(`AUTHOR_UUID`)
		REFERENCES `ACCOUNT` (`UUID`)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `EXERCISE`;

