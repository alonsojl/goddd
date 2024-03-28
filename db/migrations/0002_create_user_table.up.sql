CREATE TABLE IF NOT EXISTS `user` (
	`id` 		 INT PRIMARY KEY AUTO_INCREMENT,
	`name` 	     VARCHAR(20) NOT NULL,
	`first_name` VARCHAR(30) NOT NULL,
	`last_name`  VARCHAR(30) NOT NULL,
	`email`      VARCHAR(30) NOT NULL,
	`phone`      VARCHAR(10) NULL,
	`age`        TINYINT NOT NULL
);
