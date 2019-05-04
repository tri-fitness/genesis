
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE `WORKOUT` (
	`UUID`			VARCHAR(36)		NOT NULL,
	`NAME`			VARCHAR(128)	NOT NULL,
	`DESCRIPTION`	VARCHAR(128)	NOT NULL,
	`DURATION`		INTEGER			NOT NULL,
	`DIFFICULTY`	VARCHAR(128)	NOT NULL,
	`COMPLETIONS`	INTEGER			NOT NULL,
	`AUTHOR_UUID`	VARCHAR(36)		NOT NULL,

	PRIMARY KEY(`UUID`),
	FOREIGN KEY(`AUTHOR_UUID`)
		REFERENCES `ACCOUNT` (`UUID`)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `WORKOUT`;

