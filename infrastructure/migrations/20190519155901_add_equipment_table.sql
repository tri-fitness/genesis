
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE `EQUIPMENT` (
	`UUID`		VARCHAR(36)		NOT NULL,
	`NAME`		VARCHAR(128)	NOT NULL,
	`TYPE`		VARCHAR(128)	NOT NULL,

	PRIMARY KEY (`UUID`),
	UNIQUE KEY (`NAME`)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `EQUIPMENT`;

