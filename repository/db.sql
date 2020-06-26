DROP DATABASE IF EXISTS `Shorters`;
CREATE DATABASE `Shorters`;

USE `Shorters`;

DROP TABLE IF EXISTS `Users`;
CREATE TABLE `Users` (
  `Email`      VARCHAR(256)  COLLATE latin1_general_cs,
  PRIMARY KEY  (`Email`)
);

DROP TABLE IF EXISTS `Links`;
CREATE TABLE `Links` (
  `Key`         CHAR(12)      COLLATE latin1_general_cs,
  `URL`         VARCHAR(1024) NOT NULL,
  `Visits`      INTEGER       DEFAULT 0,
  `Creator`     VARCHAR(256)  COLLATE latin1_general_cs,
  `CreatedTime` TIMESTAMP     DEFAULT CURRENT_TIMESTAMP,
  `ExpiredTime` TIMESTAMP,
  PRIMARY KEY   (`Key`),
  FOREIGN KEY   (`Creator`)   REFERENCES `Users`(`Email`)
);
