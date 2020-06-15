-- -----------------------------------------------------
-- Schema yofio-test
-- -----------------------------------------------------
DROP SCHEMA IF EXISTS `yofio-test` ;
CREATE SCHEMA IF NOT EXISTS `yofio-test` DEFAULT CHARACTER SET utf8 ;
USE `yofio-test` ;


-- -----------------------------------------------------
-- Table Users
-- -----------------------------------------------------
DROP TABLE IF EXISTS `yofio-test`.`Users` ;
CREATE  TABLE IF NOT EXISTS `yofio-test`.`Users` (
  `ID` INT  AUTOINCREMENT ,
  `FullName` VARCHAR(150) NOT NULL ,
  `Phone` VARCHAR(13) NOT NULL,
  `Email` VARCHAR(100) NOT NULL,
  `Pwd` VARCHAR(100) NOT NULL,
  `LoginFailures` INT(1) ,
  `IsBlocked` TINYINT(1) ,
  `InsertedAt` DATETIME ,
  `UpdatedAt` DATETIME ,
  `IsDeleted` TINYINT(1),
  PRIMARY KEY (`ID`) );



-- -----------------------------------------------------
-- Table Memberships
-- -----------------------------------------------------
DROP TABLE IF EXISTS `yofio-test`.`Memberships` ;
CREATE  TABLE IF NOT EXISTS `yofio-test`.`Memberships` (
  `ID` INT  AUTOINCREMENT ,
  `UserID` INT ,
  `CardNumber` VARCHAR(16) NOT NULL,
  `ExpDate` DATETIME ,
  `Owner` INT(1) NOT NULL,
  `CVV` varchar(3) ,
  `Tier` VARCHAR(10) ,
  `Pricing` INT ,
  `InsertedAt` DATETIME ,
  `UpdatedAt` DATETIME ,
  `IsBlocked` TINYINT(1) ,
  `BlockedAt` DATETIME ,
  `IsDeleted` TINYINT(1) ,
  PRIMARY KEY (`ID`) )
  CONSTRAINT `UserID`
  FOREIGN KEY (`UserID`)
  REFERENCES `yofio-test`.`Users` (`ID`);



-- -----------------------------------------------------
-- Table Payments
-- -----------------------------------------------------
DROP TABLE IF EXISTS `yofio-test`.`Payments` ;
  CREATE  TABLE IF NOT EXISTS `yofio-test`.`Payments` (
  `ID` INT  AUTOINCREMENT ,
  `MembershipID` INT ,
  `UserID` INT ,
  `Status` VARCHAR(16) ,
  `Amount` INT(8),
  `InsertedAt` DATETIME ,
  `UpdatedAt` DATETIME ,
  `IsDeleted` TINYINT(1) ,
  PRIMARY KEY (`ID`) )
  FOREIGN KEY (`MembershipID`)
  REFERENCES `yofio-test`.`Memberships` (`ID`);
  FOREIGN KEY (`UserID`)
  REFERENCES `yofio-test`.`Users` (`ID`);
