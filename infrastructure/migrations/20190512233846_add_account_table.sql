
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE `ACCOUNT` (
	`UUID`					VARCHAR(36)		NOT NULL,
	`TYPE`					VARCHAR(128)	NOT NULL,
	`GIVEN_NAME`			VARCHAR(128)	NOT NULL,
	`SURNAME`				VARCHAR(128)	NOT NULL,
	`BIO`					TEXT			NULL,
	`EMAIL_ADDRESS`			VARCHAR(128)	NOT NULL,
	`PHONE_NUMBER`			VARCHAR(16)		NOT NULL,
	`PRIMARY_CREDENTIAL`	VARCHAR(128)	NOT NULL,
	`SECONDARY_CREDENTIAL`	VARCHAR(128)	NOT NULL,
	`TARGET_UUID`			VARCHAR(36)		NULL,
	`CREATED_AT`			DATETIME		NOT NULL,
	`UPDATED_AT`			DATETIME		NOT NULL,
	`DELETED_AT`			DATETIME		NULL,
	
	PRIMARY KEY (`UUID`),
	UNIQUE KEY (`PHONE_NUMBER`),
	UNIQUE KEY (`EMAIL_ADDRESS`)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `ACCOUNT`;

