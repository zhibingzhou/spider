/*
SQLyog Ultimate v12.09 (64 bit)
MySQL - 8.0.19 : Database - xinhe
*********************************************************************
*/

/*!40101 SET NAMES utf8 */;

/*!40101 SET SQL_MODE=''*/;

/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
CREATE DATABASE /*!32312 IF NOT EXISTS*/`xinhe` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin */ /*!80016 DEFAULT ENCRYPTION='N' */;

USE `xinhe`;

/*Table structure for table `country` */

DROP TABLE IF EXISTS `country`;

CREATE TABLE `country` (
  `id` int NOT NULL AUTO_INCREMENT,
  `country_name` varchar(50) NOT NULL COMMENT '国家名称',
  PRIMARY KEY (`id`),
  UNIQUE KEY `index_title` (`country_name`)
) ENGINE=InnoDB AUTO_INCREMENT=257 DEFAULT CHARSET=utf8 COMMENT='影片出处';

/*Table structure for table `film` */

DROP TABLE IF EXISTS `film`;

CREATE TABLE `film` (
  `id` int NOT NULL AUTO_INCREMENT,
  `en_name` varchar(50) NOT NULL COMMENT '中文名',
  `first_name` varchar(50) NOT NULL COMMENT '原名',
  `second_name` varchar(500) DEFAULT NULL COMMENT '又名',
  `year` int DEFAULT NULL COMMENT '时间',
  `country` varchar(50) DEFAULT NULL COMMENT '国家',
  `score` decimal(28,2) DEFAULT '0.00' COMMENT '评分',
  `show_time` varchar(500) DEFAULT NULL COMMENT '上映时间',
  `title` varchar(1500) DEFAULT NULL COMMENT '简介',
  `url_image` varchar(500) DEFAULT NULL COMMENT '图片地址',
  `video_url` varchar(500) DEFAULT NULL COMMENT '播放地址',
  `video_type` int DEFAULT NULL COMMENT '类型，1电影，2电视剧，3综艺',
  PRIMARY KEY (`id`),
  KEY `index_year` (`year`),
  KEY `index_country` (`country`),
  KEY `index_en_name` (`en_name`),
  KEY `index_en_type` (`video_type`)
) ENGINE=InnoDB AUTO_INCREMENT=35461812 DEFAULT CHARSET=utf8 COMMENT='电影信息';

/*Table structure for table `film_person` */

DROP TABLE IF EXISTS `film_person`;

CREATE TABLE `film_person` (
  `id` int NOT NULL AUTO_INCREMENT,
  `film_id` int DEFAULT NULL COMMENT '电影id',
  `name` varchar(50) NOT NULL COMMENT '演员名称',
  `identity` int DEFAULT NULL COMMENT '身份,1,导演，2编剧，3主演',
  PRIMARY KEY (`id`),
  UNIQUE KEY `index_id` (`film_id`,`identity`,`name`) USING BTREE,
  KEY `index_film` (`film_id`),
  KEY `index_identity` (`identity`),
  KEY `index_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=61434 DEFAULT CHARSET=utf8 COMMENT='演员';

/*Table structure for table `film_title` */

DROP TABLE IF EXISTS `film_title`;

CREATE TABLE `film_title` (
  `id` int NOT NULL AUTO_INCREMENT,
  `film_id` int DEFAULT NULL COMMENT '电影id',
  `title_id` int DEFAULT NULL COMMENT '标签id',
  PRIMARY KEY (`id`),
  UNIQUE KEY `index_id_title` (`film_id`,`title_id`) USING BTREE,
  KEY `index_film` (`film_id`),
  KEY `index_id` (`title_id`)
) ENGINE=InnoDB AUTO_INCREMENT=26702 DEFAULT CHARSET=utf8 COMMENT='电影标签';

/*Table structure for table `film_type` */

DROP TABLE IF EXISTS `film_type`;

CREATE TABLE `film_type` (
  `id` int NOT NULL AUTO_INCREMENT,
  `film_id` int DEFAULT NULL COMMENT '电影id',
  `type_id` int DEFAULT NULL COMMENT '类型id',
  PRIMARY KEY (`id`),
  KEY `index_film` (`film_id`),
  KEY `index_id` (`type_id`)
) ENGINE=InnoDB AUTO_INCREMENT=34927 DEFAULT CHARSET=utf8 COMMENT='电影类型';

/*Table structure for table `title` */

DROP TABLE IF EXISTS `title`;

CREATE TABLE `title` (
  `id` int NOT NULL AUTO_INCREMENT,
  `title_name` varchar(50) NOT NULL COMMENT '标签名称',
  PRIMARY KEY (`id`),
  UNIQUE KEY `index_title` (`title_name`)
) ENGINE=InnoDB AUTO_INCREMENT=4159 DEFAULT CHARSET=utf8 COMMENT='标签';

/*Table structure for table `type_all` */

DROP TABLE IF EXISTS `type_all`;

CREATE TABLE `type_all` (
  `id` int NOT NULL AUTO_INCREMENT,
  `type_name` varchar(50) NOT NULL COMMENT '类型名称',
  PRIMARY KEY (`id`),
  UNIQUE KEY `index_name` (`type_name`)
) ENGINE=InnoDB AUTO_INCREMENT=70 DEFAULT CHARSET=utf8 COMMENT='影片类型';

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
