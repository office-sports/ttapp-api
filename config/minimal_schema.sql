-- MySQL dump 10.13  Distrib 8.0.33, for macos13.3 (x86_64)
--
-- Host: 127.0.0.1    Database: ttapp
-- ------------------------------------------------------
-- Server version	8.0.27

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `game`
--

DROP TABLE IF EXISTS `game`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `game` (
  `id` int NOT NULL AUTO_INCREMENT,
  `game_mode_id` int DEFAULT NULL,
  `home_player_id` int NOT NULL,
  `away_player_id` int NOT NULL,
  `tournament_id` int DEFAULT NULL,
  `is_finished` tinyint(1) NOT NULL,
  `is_abandoned` tinyint(1) NOT NULL,
  `is_walkover` tinyint(1) NOT NULL,
  `date_of_match` datetime NOT NULL,
  `tournament_group_id` int DEFAULT NULL,
  `winner_id` int NOT NULL,
  `home_score` smallint NOT NULL,
  `away_score` smallint NOT NULL,
  `date_played` datetime DEFAULT NULL,
  `current_set` smallint NOT NULL DEFAULT '1',
  `server_id` int DEFAULT NULL,
  `play_order` int DEFAULT NULL,
  `stage` int DEFAULT NULL,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `level` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `playoff_home_player_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `playoff_away_player_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `old_home_elo` int DEFAULT NULL,
  `old_away_elo` int DEFAULT NULL,
  `new_home_elo` int DEFAULT NULL,
  `new_away_elo` int DEFAULT NULL,
  `office_id` int DEFAULT NULL,
  `announced` tinyint DEFAULT '1',
  `ts` varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `IDX_232B318CE227FA65` (`game_mode_id`),
  KEY `IDX_232B318C33D1A3E7` (`tournament_id`),
  KEY `IDX_232B318C1EEA3FA2` (`tournament_group_id`),
  KEY `IDX_232B318CE7328C9B` (`home_player_id`),
  KEY `IDX_232B318C6861DE1` (`away_player_id`),
  KEY `IDX_232B318CFFA0C224` (`office_id`),
  CONSTRAINT `FK_232B318C1EEA3FA2` FOREIGN KEY (`tournament_group_id`) REFERENCES `tournament_group` (`id`),
  CONSTRAINT `FK_232B318C33D1A3E7` FOREIGN KEY (`tournament_id`) REFERENCES `tournament` (`id`),
  CONSTRAINT `FK_232B318CE227FA65` FOREIGN KEY (`game_mode_id`) REFERENCES `game_mode` (`id`),
  CONSTRAINT `FK_232B318CFFA0C224` FOREIGN KEY (`office_id`) REFERENCES `office` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1727 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `game`
--

INSERT INTO `game` VALUES (1726,5,176,177,21,0,0,0,'2023-09-21 16:52:30',59,0,0,0,NULL,1,176,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,5,0,'0');

--
-- Table structure for table `game_mode`
--

DROP TABLE IF EXISTS `game_mode`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `game_mode` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `short_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `wins_required` smallint NOT NULL,
  `max_sets` smallint NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `game_mode`
--

INSERT INTO `game_mode` VALUES (1,'Best of 1','BO1',1,1),(2,'Best of 2','BO2',2,2),(3,'Best of 3','BO3',2,3),(4,'Best of 4','BO4',3,4),(5,'Best of 5','BO5',3,5),(6,'Best of 6','BO6',4,6),(7,'Best of 7','BO7',4,7);

--
-- Table structure for table `level`
--

DROP TABLE IF EXISTS `level`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `level` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `level`
--


--
-- Table structure for table `office`
--

DROP TABLE IF EXISTS `office`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `office` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `is_default` tinyint(1) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `office`
--

INSERT INTO `office` VALUES (5,'Office One',1),(6,'Office Two',0);

--
-- Table structure for table `player`
--

DROP TABLE IF EXISTS `player`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `player` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `nickname` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `tournament_elo` int NOT NULL DEFAULT '1500',
  `current_elo` int NOT NULL DEFAULT '1500',
  `display_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `slack_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `profile_pic_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `office_id` int DEFAULT NULL,
  `tournament_elo_previous` int DEFAULT NULL,
  `active` int NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`),
  KEY `IDX_98197A65FFA0C224` (`office_id`),
  CONSTRAINT `FK_98197A65FFA0C224` FOREIGN KEY (`office_id`) REFERENCES `office` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=178 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `player`
--

INSERT INTO `player` VALUES (176,'John Doe','JD',1500,1500,'John Doe','-','-',5,NULL,1),(177,'Jane Doe','JDO',1500,1500,'Jane Doe','-','-',5,NULL,1);

--
-- Table structure for table `player_tournament_group`
--

DROP TABLE IF EXISTS `player_tournament_group`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `player_tournament_group` (
  `id` int NOT NULL AUTO_INCREMENT,
  `player_id` int NOT NULL,
  `group_id` int NOT NULL,
  `tournament_id` int NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=341 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `player_tournament_group`
--

INSERT INTO `player_tournament_group` VALUES (339,176,59,21),(340,177,59,21);

--
-- Table structure for table `points`
--

DROP TABLE IF EXISTS `points`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `points` (
  `id` int NOT NULL AUTO_INCREMENT,
  `score_id` int NOT NULL,
  `is_home_point` tinyint(1) NOT NULL,
  `is_away_point` tinyint(1) NOT NULL,
  `time` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `IDX_27BA8E2912EB0A51` (`score_id`),
  CONSTRAINT `FK_27BA8E2912EB0A51` FOREIGN KEY (`score_id`) REFERENCES `scores` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `points`
--


--
-- Table structure for table `scores`
--

DROP TABLE IF EXISTS `scores`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `scores` (
  `id` int NOT NULL AUTO_INCREMENT,
  `game_id` int NOT NULL,
  `set_number` smallint NOT NULL,
  `home_points` int NOT NULL,
  `away_points` int NOT NULL,
  PRIMARY KEY (`id`),
  KEY `IDX_750375EE48FD905` (`game_id`),
  CONSTRAINT `FK_750375EE48FD905` FOREIGN KEY (`game_id`) REFERENCES `game` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4532 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `scores`
--


--
-- Table structure for table `spectators`
--

DROP TABLE IF EXISTS `spectators`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `spectators` (
  `id` int NOT NULL AUTO_INCREMENT,
  `game_id` int NOT NULL,
  `spectators` int NOT NULL,
  `pit` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `spectators`
--


--
-- Table structure for table `tournament`
--

DROP TABLE IF EXISTS `tournament`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `tournament` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `is_finished` tinyint(1) NOT NULL,
  `is_playoffs` tinyint(1) NOT NULL,
  `start_time` datetime DEFAULT NULL,
  `is_official` tinyint(1) NOT NULL,
  `parent_tournament` int DEFAULT NULL,
  `office_id` int DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `IDX_BD5FB8D9FFA0C224` (`office_id`),
  CONSTRAINT `FK_BD5FB8D9FFA0C224` FOREIGN KEY (`office_id`) REFERENCES `office` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=22 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `tournament`
--

INSERT INTO `tournament` VALUES (21,'First Tournament',0,0,'2023-09-21 16:45:58',1,NULL,5);

--
-- Table structure for table `tournament_group`
--

DROP TABLE IF EXISTS `tournament_group`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `tournament_group` (
  `id` int NOT NULL AUTO_INCREMENT,
  `tournament_id` int NOT NULL,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `priority` smallint NOT NULL,
  `abbreviation` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `is_official` tinyint(1) NOT NULL,
  `color_template` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `promotions` int NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `IDX_AC96D4DC33D1A3E7` (`tournament_id`),
  CONSTRAINT `FK_6DC044C533D1A3E7` FOREIGN KEY (`tournament_id`) REFERENCES `tournament` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=61 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `tournament_group`
--

INSERT INTO `tournament_group` VALUES (59,21,'Group One',1,'GRP1',1,'1.1.2.2',2),(60,21,'Group Two',1,'GRP2',1,'1.1.2.2',2);

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `login` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `hash` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2023-09-21 16:57:05
