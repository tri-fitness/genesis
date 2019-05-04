
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE `WORKOUT_EXERCISE` (
	`WORKOUT_UUID`	VARCHAR(36)		NOT NULL,
	`EXERCISE_UUID`	VARCHAR(36)		NOT NULL,

	PRIMARY KEY(`WORKOUT_UUID`, `EXERCISE_UUID`)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `WORKOUT_EXERCISE`;

