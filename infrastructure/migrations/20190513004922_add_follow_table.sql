
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE `FOLLOW`(
	`FOLLOWER_UUID` VARCHAR(36) NOT NULL,
	`FOLLOWEE_UUID` VARCHAR(36) NOT NULL,

	PRIMARY KEY(`FOLLOWEE_UUID`, `FOLLOWER_UUID`),
	FOREIGN KEY(`FOLLOWER_UUID`)
		REFERENCES `ACCOUNT` (`UUID`),
	FOREIGN KEY(`FOLLOWEE_UUID`)
		REFERENCES `ACCOUNT` (`UUID`)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `FOLLOW`;

