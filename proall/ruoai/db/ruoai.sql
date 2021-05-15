/*
SQLyog Ultimate v12.09 (64 bit)
MySQL - 8.0.19 : Database - ruoai
*********************************************************************
*/

/*!40101 SET NAMES utf8 */;

/*!40101 SET SQL_MODE=''*/;

/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
CREATE DATABASE /*!32312 IF NOT EXISTS*/`ruoai` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;

USE `ruoai`;

/*Table structure for table `city` */

DROP TABLE IF EXISTS `city`;

CREATE TABLE `city` (
  `id` int NOT NULL AUTO_INCREMENT,
  `title` varchar(50) NOT NULL COMMENT '城市名称',
  `province_id` int NOT NULL COMMENT '省份',
  PRIMARY KEY (`id`),
  UNIQUE KEY `index_city` (`title`),
  UNIQUE KEY `index_pro_title` (`title`,`province_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=411 DEFAULT CHARSET=utf8 COMMENT='城市';

/*Table structure for table `education` */

DROP TABLE IF EXISTS `education`;

CREATE TABLE `education` (
  `id` int NOT NULL AUTO_INCREMENT,
  `title` varchar(50) NOT NULL COMMENT '学历名称',
  PRIMARY KEY (`id`),
  UNIQUE KEY `index_pro` (`title`)
) ENGINE=InnoDB AUTO_INCREMENT=59 DEFAULT CHARSET=utf8 COMMENT='学历';

/*Table structure for table `province` */

DROP TABLE IF EXISTS `province`;

CREATE TABLE `province` (
  `id` int NOT NULL AUTO_INCREMENT,
  `title` varchar(50) NOT NULL COMMENT '省份名称',
  PRIMARY KEY (`id`),
  UNIQUE KEY `index_pro` (`title`)
) ENGINE=InnoDB AUTO_INCREMENT=41 DEFAULT CHARSET=utf8 COMMENT='省份';

/*Table structure for table `user` */

DROP TABLE IF EXISTS `user`;

CREATE TABLE `user` (
  `id` int NOT NULL AUTO_INCREMENT,
  `nickname` varchar(50) NOT NULL COMMENT '网名',
  `city_code` int DEFAULT NULL COMMENT '城市',
  `age` int DEFAULT NULL COMMENT '年龄',
  `gender` varchar(20) DEFAULT NULL COMMENT '性别',
  `education` int DEFAULT NULL COMMENT '学历',
  `birthday_time` datetime DEFAULT NULL COMMENT '出身日期',
  `height` varchar(20) DEFAULT NULL COMMENT '身高',
  `weight` varchar(20) DEFAULT NULL COMMENT '体重',
  `married` varchar(20) DEFAULT NULL COMMENT '婚姻',
  `house` varchar(20) DEFAULT NULL COMMENT '是否购房',
  `salary_up` int DEFAULT NULL COMMENT '最高工资',
  `work` int DEFAULT NULL COMMENT '职业',
  `salary_down` int DEFAULT NULL COMMENT '最低工资',
  `hometown` int DEFAULT NULL COMMENT '家乡',
  `title` varchar(500) DEFAULT NULL COMMENT '爱情宣言',
  `url_image` varchar(300) DEFAULT NULL COMMENT '图片地址',
  PRIMARY KEY (`id`),
  KEY `index_city` (`city_code`),
  KEY `index_education` (`education`),
  KEY `index_work` (`work`),
  KEY `index_gen` (`gender`)
) ENGINE=InnoDB AUTO_INCREMENT=674579 DEFAULT CHARSET=utf8 COMMENT='用户信息';

/*Table structure for table `work` */

DROP TABLE IF EXISTS `work`;

CREATE TABLE `work` (
  `id` int NOT NULL AUTO_INCREMENT,
  `title` varchar(50) NOT NULL COMMENT '职业名称',
  PRIMARY KEY (`id`),
  UNIQUE KEY `index_pro` (`title`)
) ENGINE=InnoDB AUTO_INCREMENT=11591 DEFAULT CHARSET=utf8 COMMENT='职业';

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
