
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE `ACCOUNT_EXERCISE`(
	`ACCOUNT_UUID`	VARCHAR(36)		NOT NULL,
	`EXERCISE_UUID`	VARCHAR(36)		NOT NULL,
	`WATCHES`		INTEGER			NOT NULL,

	PRIMARY KEY(`ACCOUNT_UUID`, `EXERCISE_UUID`),
	FOREIGN KEY(`ACCOUNT_UUID`)
		REFERENCES `ACCOUNT` (`UUID`),
	FOREIGN KEY(`EXERCISE_UUID`)
		REFERENCES `EXERCISE` (`UUID`)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `ACCOUNT_EXERCISE`;

