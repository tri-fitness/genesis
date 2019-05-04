
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE `ACCOUNT_WORKOUT`(
	`ACCOUNT_UUID`	VARCHAR(36)	NOT NULL,
	`WORKOUT_UUID`	VARCHAR(36)	NOT NULL,
	`COMPLETIONS`	INTEGER			NOT NULL,

	PRIMARY KEY(`ACCOUNT_UUID`, `WORKOUT_UUID`),
	FOREIGN KEY(`ACCOUNT_UUID`)
		REFERENCES `ACCOUNT` (`UUID`),
	FOREIGN KEY(`WORKOUT_UUID`)
		REFERENCES `WORKOUT` (`UUID`)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `ACCOUNT_WORKOUT`;

