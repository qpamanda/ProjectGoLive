-- MySQL dump 10.13  Distrib 8.0.25, for Win64 (x86_64)
--
-- Host: 127.0.0.1    Database: dbProject
-- ------------------------------------------------------
-- Server version	8.0.25

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `Category`
--

DROP TABLE IF EXISTS `Category`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `Category` (
  `CategoryID` int NOT NULL,
  `Category` varchar(25) DEFAULT NULL,
  `CreatedBy` varchar(25) DEFAULT NULL,
  `Created_dt` datetime DEFAULT NULL,
  `LastModifiedBy` varchar(25) DEFAULT NULL,
  `LastModified_dt` datetime DEFAULT NULL,
  PRIMARY KEY (`CategoryID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Category`
--

LOCK TABLES `Category` WRITE;
/*!40000 ALTER TABLE `Category` DISABLE KEYS */;
INSERT INTO `Category` VALUES (1,'Donation (Monetary)','admin','2021-06-27 09:41:38','admin','2021-06-27 09:41:38'),(2,'Donation (Physical Items)','admin','2021-06-27 09:41:38','admin','2021-06-27 09:41:38'),(3,'Errands','admin','2021-06-27 09:41:38','admin','2021-06-27 09:41:38');
/*!40000 ALTER TABLE `Category` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `MemberType`
--

DROP TABLE IF EXISTS `MemberType`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `MemberType` (
  `MemberTypeID` int NOT NULL,
  `MemberType` varchar(25) DEFAULT NULL,
  `CreatedBy` varchar(25) DEFAULT NULL,
  `Created_dt` datetime DEFAULT NULL,
  `LastModifiedBy` varchar(25) DEFAULT NULL,
  `LastModified_dt` datetime DEFAULT NULL,
  PRIMARY KEY (`MemberTypeID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `MemberType`
--

LOCK TABLES `MemberType` WRITE;
/*!40000 ALTER TABLE `MemberType` DISABLE KEYS */;
INSERT INTO `MemberType` VALUES (1,'Admin','admin','2021-06-26 08:45:43','admin','2021-06-26 08:45:43'),(2,'Requester','admin','2021-06-26 08:45:43','admin','2021-06-26 08:45:43'),(3,'Helper','admin','2021-06-26 08:45:43','admin','2021-06-26 08:45:43');
/*!40000 ALTER TABLE `MemberType` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `RequestStatus`
--

DROP TABLE IF EXISTS `RequestStatus`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `RequestStatus` (
  `StatusCode` int NOT NULL,
  `Status` varchar(25) DEFAULT NULL,
  `CreatedBy` varchar(25) DEFAULT NULL,
  `Created_dt` datetime DEFAULT NULL,
  `LastModifiedBy` varchar(25) DEFAULT NULL,
  `LastModified_dt` datetime DEFAULT NULL,
  PRIMARY KEY (`StatusCode`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `RequestStatus`
--

LOCK TABLES `RequestStatus` WRITE;
/*!40000 ALTER TABLE `RequestStatus` DISABLE KEYS */;
INSERT INTO `RequestStatus` VALUES (0,'Pending/To be Assigned','admin','2021-06-27 11:11:55','admin','2021-06-27 11:11:55'),(1,'Being Handled/Assigned','admin','2021-06-27 11:11:56','admin','2021-06-27 11:11:56'),(2,'Completed','admin','2021-06-27 11:11:56','admin','2021-06-27 11:11:56');
/*!40000 ALTER TABLE `RequestStatus` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2021-06-28 16:39:50
