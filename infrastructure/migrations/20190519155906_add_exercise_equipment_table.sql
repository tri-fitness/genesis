
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE `EXERCISE_EQUIPMENT` (
	`EXERCISE_UUID`		VARCHAR(36)		NOT NULL,
	`EQUIPMENT_UUID`	VARCHAR(36)		NOT NULL,

	PRIMARY KEY (`EXERCISE_UUID`, `EQUIPMENT_UUID`)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `EXERCISE_EQUIPMENT`;

