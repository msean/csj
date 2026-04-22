-- MySQL dump 10.13  Distrib 8.0.43, for macos15 (arm64)
--
-- Host: 20.2.233.88    Database: caishuji
-- ------------------------------------------------------
-- Server version	8.0.19

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
-- Table structure for table `ban_record`
--

DROP TABLE IF EXISTS `ban_record`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ban_record` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `bot_id` bigint DEFAULT NULL COMMENT '机器人ID',
  `user_id` bigint DEFAULT NULL COMMENT '用户ID',
  `user_name` varchar(128) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '用户名',
  `full_name` varchar(128) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '用户名',
  `chat_id` bigint DEFAULT NULL COMMENT '所在群聊',
  `ban_duration` bigint DEFAULT NULL COMMENT '封禁时长',
  `ban_type` bigint DEFAULT NULL COMMENT '封禁时长',
  `remark` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '封禁时长',
  `msg` text COLLATE utf8mb4_general_ci COMMENT '禁用信息',
  `lifting_time` datetime(3) DEFAULT NULL COMMENT '解禁时间',
  PRIMARY KEY (`id`),
  KEY `idx_ban_record_deleted_at` (`deleted_at`),
  KEY `idx_ban_record_user_name` (`user_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ban_record`
--

LOCK TABLES `ban_record` WRITE;
/*!40000 ALTER TABLE `ban_record` DISABLE KEYS */;
/*!40000 ALTER TABLE `ban_record` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `bot`
--

DROP TABLE IF EXISTS `bot`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `bot` (
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '机器人名称',
  `is_for_ledger` bigint DEFAULT '2' COMMENT '是否用来记账(1开2关)',
  `is_for_ledger2` bigint DEFAULT '2' COMMENT '是否用来记账(1开2关)',
  `is_for_msg_mgr` bigint DEFAULT '2' COMMENT '是否用来记账(1开2关)',
  `is_for_msg_mass` bigint DEFAULT '2' COMMENT '是否用来记账(1开2关)',
  `is_for_ad_publish` bigint DEFAULT '2' COMMENT '是否用来广告自动发布',
  `bot_id` bigint DEFAULT NULL COMMENT '机器人ID',
  `token` varchar(256) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'token',
  KEY `idx_bot_deleted_at` (`deleted_at`),
  KEY `idx_bot_bot_id` (`bot_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `bot`
--

LOCK TABLES `bot` WRITE;
/*!40000 ALTER TABLE `bot` DISABLE KEYS */;
/*!40000 ALTER TABLE `bot` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `bot_ad_publish_record`
--

DROP TABLE IF EXISTS `bot_ad_publish_record`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `bot_ad_publish_record` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `bot_id` bigint DEFAULT NULL,
  `bot_name` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `channel_name` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `publish_times` bigint DEFAULT NULL,
  `user_id` bigint DEFAULT NULL,
  `user_name` varchar(256) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `channel_id` bigint DEFAULT NULL,
  `price` decimal(10,3) DEFAULT NULL,
  `content` text COLLATE utf8mb4_general_ci COMMENT '发布内容',
  PRIMARY KEY (`id`),
  KEY `idx_bot_ad_publish_record_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `bot_ad_publish_record`
--

LOCK TABLES `bot_ad_publish_record` WRITE;
/*!40000 ALTER TABLE `bot_ad_publish_record` DISABLE KEYS */;
/*!40000 ALTER TABLE `bot_ad_publish_record` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `bot_ban_content`
--

DROP TABLE IF EXISTS `bot_ban_content`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `bot_ban_content` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `ban_content` varchar(1024) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'ban_content',
  `bot_id` bigint DEFAULT NULL COMMENT '机器人ID',
  PRIMARY KEY (`id`),
  KEY `idx_bot_ban_content_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `bot_ban_content`
--

LOCK TABLES `bot_ban_content` WRITE;
/*!40000 ALTER TABLE `bot_ban_content` DISABLE KEYS */;
/*!40000 ALTER TABLE `bot_ban_content` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `bot_ban_group_mem`
--

DROP TABLE IF EXISTS `bot_ban_group_mem`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `bot_ban_group_mem` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `bot_id` bigint DEFAULT NULL,
  `chat_group_id` bigint DEFAULT NULL COMMENT 'chatGroupID',
  `ban_mem_content` varchar(256) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '封禁成员名称(成员还有该字段就封禁)',
  PRIMARY KEY (`id`),
  KEY `idx_bot_ban_group_mem_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `bot_ban_group_mem`
--

LOCK TABLES `bot_ban_group_mem` WRITE;
/*!40000 ALTER TABLE `bot_ban_group_mem` DISABLE KEYS */;
/*!40000 ALTER TABLE `bot_ban_group_mem` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `bot_channel`
--

DROP TABLE IF EXISTS `bot_channel`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `bot_channel` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `bot_id` bigint DEFAULT NULL COMMENT '机器人id',
  `channel_id` bigint DEFAULT NULL COMMENT '频道ID',
  `channel_name` varchar(256) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '渠道名称',
  PRIMARY KEY (`id`),
  KEY `idx_bot_channel_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `bot_channel`
--

LOCK TABLES `bot_channel` WRITE;
/*!40000 ALTER TABLE `bot_channel` DISABLE KEYS */;
/*!40000 ALTER TABLE `bot_channel` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `bot_chat_group`
--

DROP TABLE IF EXISTS `bot_chat_group`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `bot_chat_group` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `bot_id` bigint DEFAULT NULL,
  `chat_group_id` bigint DEFAULT NULL,
  `chat_group_name` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `ban_forward` bigint DEFAULT '1',
  `max_words` bigint DEFAULT '-1',
  `sync_message` bigint DEFAULT '2',
  `must_join_channels` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `invaid_channel_fold_link` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `members` text COLLATE utf8mb4_general_ci,
  PRIMARY KEY (`id`),
  KEY `idx_bot_chat_group_deleted_at` (`deleted_at`),
  KEY `idx_bot_group` (`bot_id`,`chat_group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `bot_chat_group`
--

LOCK TABLES `bot_chat_group` WRITE;
/*!40000 ALTER TABLE `bot_chat_group` DISABLE KEYS */;
/*!40000 ALTER TABLE `bot_chat_group` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `bot_chat_group_classify`
--

DROP TABLE IF EXISTS `bot_chat_group_classify`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `bot_chat_group_classify` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `title` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `chat_groups` text COLLATE utf8mb4_general_ci COMMENT '群组列表',
  `users` text COLLATE utf8mb4_general_ci COMMENT '群组列表',
  PRIMARY KEY (`id`),
  KEY `idx_bot_chat_group_classify_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `bot_chat_group_classify`
--

LOCK TABLES `bot_chat_group_classify` WRITE;
/*!40000 ALTER TABLE `bot_chat_group_classify` DISABLE KEYS */;
/*!40000 ALTER TABLE `bot_chat_group_classify` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `bot_cmd_config`
--

DROP TABLE IF EXISTS `bot_cmd_config`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `bot_cmd_config` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `bot_id` bigint DEFAULT NULL COMMENT '机器人ID',
  `title` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '标题',
  `content` text COLLATE utf8mb4_general_ci COMMENT '开始设置内容',
  `cmd` varchar(512) COLLATE utf8mb4_general_ci DEFAULT '/start' COMMENT '开始设置内容',
  `cmd_buttons` text COLLATE utf8mb4_general_ci COMMENT '命令按钮配置',
  `type` bigint DEFAULT '1' COMMENT '类型',
  PRIMARY KEY (`id`),
  KEY `idx_bot_cmd_config_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `bot_cmd_config`
--

LOCK TABLES `bot_cmd_config` WRITE;
/*!40000 ALTER TABLE `bot_cmd_config` DISABLE KEYS */;
/*!40000 ALTER TABLE `bot_cmd_config` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `bot_group_related_channnel_follow`
--

DROP TABLE IF EXISTS `bot_group_related_channnel_follow`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `bot_group_related_channnel_follow` (
  `user_id` bigint DEFAULT NULL,
  `bot_id` bigint DEFAULT NULL,
  `chat_group_id` bigint DEFAULT NULL,
  `check_time` datetime(3) DEFAULT NULL,
  KEY `idx_bot_group_related_channnel_follow_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `bot_group_related_channnel_follow`
--

LOCK TABLES `bot_group_related_channnel_follow` WRITE;
/*!40000 ALTER TABLE `bot_group_related_channnel_follow` DISABLE KEYS */;
/*!40000 ALTER TABLE `bot_group_related_channnel_follow` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `bot_ledger`
--

DROP TABLE IF EXISTS `bot_ledger`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `bot_ledger` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `opr_user_id` bigint DEFAULT NULL COMMENT '操作用户ID',
  `opr_first_name` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `opr_last_name` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `opr_user_nick_name` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '操作人昵称',
  `action_type` bigint DEFAULT NULL COMMENT '操作类型',
  `amount` decimal(32,2) DEFAULT NULL COMMENT '操作金额',
  `amount_with_fee` decimal(32,2) DEFAULT NULL COMMENT '操作金额',
  `bot_id` bigint DEFAULT NULL COMMENT '机器人ID',
  `chat_group_id` bigint DEFAULT NULL COMMENT '所在群组',
  `message_id` bigint DEFAULT NULL COMMENT '消息ID',
  `raw_input` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '原始输入',
  `remark` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '备注',
  `current_fee_rate` double DEFAULT NULL COMMENT '当前费率',
  PRIMARY KEY (`id`),
  KEY `idx_bot_ledger_deleted_at` (`deleted_at`),
  KEY `idx_bot_ledger_message_id` (`message_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `bot_ledger`
--

LOCK TABLES `bot_ledger` WRITE;
/*!40000 ALTER TABLE `bot_ledger` DISABLE KEYS */;
/*!40000 ALTER TABLE `bot_ledger` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `bot_ledger2`
--

DROP TABLE IF EXISTS `bot_ledger2`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `bot_ledger2` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `opr_user_id` bigint DEFAULT NULL COMMENT '操作用户ID',
  `opr_first_name` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `opr_last_name` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `opr_user_nick_name` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '操作人昵称',
  `action_type` bigint DEFAULT NULL COMMENT '操作类型',
  `amount` decimal(18,2) DEFAULT NULL COMMENT '操作金额',
  `user_name` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '姓名',
  `account` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '账号',
  `message_id` bigint DEFAULT NULL COMMENT '消息ID',
  `raw_input` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '原始输入',
  `remark` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '备注',
  `bot_id` bigint DEFAULT NULL COMMENT '机器人ID',
  `chat_group_id` bigint DEFAULT NULL COMMENT '群ID',
  `work_date` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '账单日期',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_msg` (`message_id`),
  KEY `idx_bot_ledger2_deleted_at` (`deleted_at`),
  KEY `idx_bot_chat_date` (`bot_id`,`chat_group_id`,`work_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `bot_ledger2`
--

LOCK TABLES `bot_ledger2` WRITE;
/*!40000 ALTER TABLE `bot_ledger2` DISABLE KEYS */;
/*!40000 ALTER TABLE `bot_ledger2` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `bot_ledger2_permission`
--

DROP TABLE IF EXISTS `bot_ledger2_permission`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `bot_ledger2_permission` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `bot_id` bigint DEFAULT NULL COMMENT '机器人ID',
  `chat_group_id` bigint DEFAULT NULL COMMENT '群聊ID',
  `opr_users` text COLLATE utf8mb4_general_ci COMMENT '操作人',
  PRIMARY KEY (`id`),
  KEY `idx_bot_ledger2_permission_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `bot_ledger2_permission`
--

LOCK TABLES `bot_ledger2_permission` WRITE;
/*!40000 ALTER TABLE `bot_ledger2_permission` DISABLE KEYS */;
/*!40000 ALTER TABLE `bot_ledger2_permission` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `bot_ledger_permission`
--

DROP TABLE IF EXISTS `bot_ledger_permission`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `bot_ledger_permission` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `current_fee_rate` double DEFAULT NULL COMMENT '当前费率',
  `bot_id` bigint DEFAULT NULL COMMENT '机器人ID',
  `chat_group_id` bigint DEFAULT NULL COMMENT '群聊ID',
  `opr_users` text COLLATE utf8mb4_general_ci COMMENT '操作人',
  PRIMARY KEY (`id`),
  KEY `idx_bot_ledger_permission_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `bot_ledger_permission`
--

LOCK TABLES `bot_ledger_permission` WRITE;
/*!40000 ALTER TABLE `bot_ledger_permission` DISABLE KEYS */;
/*!40000 ALTER TABLE `bot_ledger_permission` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `bot_mass_msg_permission`
--

DROP TABLE IF EXISTS `bot_mass_msg_permission`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `bot_mass_msg_permission` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `bot_id` bigint DEFAULT NULL COMMENT '机器人ID',
  `chat_group_id` bigint DEFAULT NULL COMMENT '群聊ID',
  `opr_users` text COLLATE utf8mb4_general_ci COMMENT '操作人',
  PRIMARY KEY (`id`),
  KEY `idx_bot_mass_msg_permission_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `bot_mass_msg_permission`
--

LOCK TABLES `bot_mass_msg_permission` WRITE;
/*!40000 ALTER TABLE `bot_mass_msg_permission` DISABLE KEYS */;
/*!40000 ALTER TABLE `bot_mass_msg_permission` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `bot_mass_msg_record`
--

DROP TABLE IF EXISTS `bot_mass_msg_record`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `bot_mass_msg_record` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `bot_id` bigint DEFAULT NULL,
  `chat_group_id` bigint DEFAULT NULL,
  `msg` text COLLATE utf8mb4_general_ci,
  `members` text COLLATE utf8mb4_general_ci,
  `remark` text COLLATE utf8mb4_general_ci,
  PRIMARY KEY (`id`),
  KEY `idx_bot_mass_msg_record_deleted_at` (`deleted_at`),
  KEY `idx_bot_group` (`bot_id`,`chat_group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `bot_mass_msg_record`
--

LOCK TABLES `bot_mass_msg_record` WRITE;
/*!40000 ALTER TABLE `bot_mass_msg_record` DISABLE KEYS */;
/*!40000 ALTER TABLE `bot_mass_msg_record` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `bot_recharge_config`
--

DROP TABLE IF EXISTS `bot_recharge_config`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `bot_recharge_config` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `bot_id` bigint DEFAULT NULL COMMENT '机器人ID',
  `publish_times` bigint DEFAULT '1' COMMENT '发布次数',
  `price` double DEFAULT NULL COMMENT '价格',
  PRIMARY KEY (`id`),
  KEY `idx_bot_recharge_config_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `bot_recharge_config`
--

LOCK TABLES `bot_recharge_config` WRITE;
/*!40000 ALTER TABLE `bot_recharge_config` DISABLE KEYS */;
/*!40000 ALTER TABLE `bot_recharge_config` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `bot_task`
--

DROP TABLE IF EXISTS `bot_task`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `bot_task` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `title` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '发送标题',
  `bot_id` bigint DEFAULT NULL COMMENT '机器人ID',
  `chat_group_id` bigint DEFAULT NULL COMMENT '群ID',
  `send_group_type` bigint DEFAULT '1' COMMENT '组类型(1、群聊 2、频道)',
  `group_id` bigint DEFAULT NULL COMMENT 'groupID',
  `task_send_type` bigint DEFAULT NULL COMMENT '发送类型',
  `content` text COLLATE utf8mb4_general_ci,
  `extrend_button` text COLLATE utf8mb4_general_ci COMMENT '扩展按钮',
  `send_interval` bigint DEFAULT NULL COMMENT '发送间隔',
  `next_send_time` datetime(3) DEFAULT NULL COMMENT '下一次发送时间',
  `stop_time` datetime(3) DEFAULT NULL COMMENT '任务结束时间',
  `stop_time_text` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '-',
  `pre_send_time` datetime(3) DEFAULT NULL COMMENT '上一次发送时间',
  `status` bigint DEFAULT NULL COMMENT '状态(1开 2关)',
  PRIMARY KEY (`id`),
  KEY `idx_bot_task_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `bot_task`
--

LOCK TABLES `bot_task` WRITE;
/*!40000 ALTER TABLE `bot_task` DISABLE KEYS */;
/*!40000 ALTER TABLE `bot_task` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `bot_user_recharge_record`
--

DROP TABLE IF EXISTS `bot_user_recharge_record`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `bot_user_recharge_record` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `bot_id` bigint DEFAULT NULL COMMENT '机器人ID',
  `publish_times` bigint DEFAULT '1' COMMENT '发布次数',
  `start_time` datetime(3) DEFAULT NULL COMMENT '发布开始时间',
  `publish_interval` bigint DEFAULT NULL COMMENT '发布间隔',
  `publish_content` text COLLATE utf8mb4_general_ci COMMENT '发布内容',
  `status` bigint DEFAULT NULL COMMENT '状态(1、创建 2、支付成功 3、支付超时失败)',
  `msg_id` bigint DEFAULT NULL,
  `channel_id` bigint DEFAULT NULL,
  `user_name` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `user_id` bigint DEFAULT NULL,
  `price` decimal(10,3) DEFAULT NULL,
  `chat_id` bigint DEFAULT NULL,
  `tx_id` varchar(256) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `payment_addr` varchar(256) COLLATE utf8mb4_general_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_bot_user_recharge_record_deleted_at` (`deleted_at`),
  KEY `idx_bot_user_recharge_record_status` (`status`),
  KEY `idx_bot_user_recharge_record_user_name` (`user_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `bot_user_recharge_record`
--

LOCK TABLES `bot_user_recharge_record` WRITE;
/*!40000 ALTER TABLE `bot_user_recharge_record` DISABLE KEYS */;
/*!40000 ALTER TABLE `bot_user_recharge_record` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `bot_user_wallet`
--

DROP TABLE IF EXISTS `bot_user_wallet`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `bot_user_wallet` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `bot_id` bigint DEFAULT NULL COMMENT '机器人名称',
  `bot_name` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `user_id` bigint DEFAULT NULL COMMENT '用户ID',
  `user_name` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '用户名',
  `balance` decimal(10,3) DEFAULT NULL COMMENT '余额',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_bot_user` (`bot_id`,`user_id`),
  KEY `idx_bot_user_wallet_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `bot_user_wallet`
--

LOCK TABLES `bot_user_wallet` WRITE;
/*!40000 ALTER TABLE `bot_user_wallet` DISABLE KEYS */;
/*!40000 ALTER TABLE `bot_user_wallet` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `casbin_rule`
--

DROP TABLE IF EXISTS `casbin_rule`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `casbin_rule` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `ptype` varchar(100) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `v0` varchar(100) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `v1` varchar(100) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `v2` varchar(100) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `v3` varchar(100) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `v4` varchar(100) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `v5` varchar(100) COLLATE utf8mb4_general_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_casbin_rule` (`ptype`,`v0`,`v1`,`v2`,`v3`,`v4`,`v5`)
) ENGINE=InnoDB AUTO_INCREMENT=217 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `casbin_rule`
--

LOCK TABLES `casbin_rule` WRITE;
/*!40000 ALTER TABLE `casbin_rule` DISABLE KEYS */;
INSERT INTO `casbin_rule` VALUES (2,'p','888','/api/createApi','POST','','',''),(5,'p','888','/api/deleteApi','POST','','',''),(8,'p','888','/api/deleteApisByIds','DELETE','','',''),(11,'p','888','/api/enterSyncApi','POST','','',''),(7,'p','888','/api/getAllApis','POST','','',''),(4,'p','888','/api/getApiById','POST','','',''),(10,'p','888','/api/getApiGroups','GET','','',''),(3,'p','888','/api/getApiList','POST','','',''),(12,'p','888','/api/ignoreApi','POST','','',''),(9,'p','888','/api/syncApi','GET','','',''),(6,'p','888','/api/updateApi','POST','','',''),(129,'p','888','/attachmentCategory/addCategory','POST','','',''),(130,'p','888','/attachmentCategory/deleteCategory','POST','','',''),(128,'p','888','/attachmentCategory/getCategoryList','GET','','',''),(13,'p','888','/authority/copyAuthority','POST','','',''),(15,'p','888','/authority/createAuthority','POST','','',''),(16,'p','888','/authority/deleteAuthority','POST','','',''),(17,'p','888','/authority/getAuthorityList','POST','','',''),(18,'p','888','/authority/setDataAuthority','POST','','',''),(14,'p','888','/authority/updateAuthority','PUT','','',''),(105,'p','888','/authorityBtn/canRemoveAuthorityBtn','POST','','',''),(104,'p','888','/authorityBtn/getAuthorityBtn','POST','','',''),(103,'p','888','/authorityBtn/setAuthorityBtn','POST','','',''),(74,'p','888','/autoCode/addFunc','POST','','',''),(67,'p','888','/autoCode/createPackage','POST','','',''),(71,'p','888','/autoCode/createPlug','POST','','',''),(64,'p','888','/autoCode/createTemp','POST','','',''),(70,'p','888','/autoCode/delPackage','POST','','',''),(65,'p','888','/autoCode/delSysHistory','POST','','',''),(62,'p','888','/autoCode/getColumn','GET','','',''),(58,'p','888','/autoCode/getDB','GET','','',''),(59,'p','888','/autoCode/getMeta','POST','','',''),(69,'p','888','/autoCode/getPackage','POST','','',''),(66,'p','888','/autoCode/getSysHistory','POST','','',''),(61,'p','888','/autoCode/getTables','GET','','',''),(68,'p','888','/autoCode/getTemplates','GET','','',''),(72,'p','888','/autoCode/installPlugin','POST','','',''),(75,'p','888','/autoCode/mcp','POST','','',''),(77,'p','888','/autoCode/mcpList','POST','','',''),(76,'p','888','/autoCode/mcpTest','POST','','',''),(60,'p','888','/autoCode/preview','POST','','',''),(73,'p','888','/autoCode/pubPlug','POST','','',''),(63,'p','888','/autoCode/rollback','POST','','',''),(48,'p','888','/casbin/getPolicyPathByAuthorityId','POST','','',''),(47,'p','888','/casbin/updateCasbin','POST','','',''),(56,'p','888','/customer/customer','DELETE','','',''),(53,'p','888','/customer/customer','GET','','',''),(55,'p','888','/customer/customer','POST','','',''),(54,'p','888','/customer/customer','PUT','','',''),(57,'p','888','/customer/customerList','GET','','',''),(98,'p','888','/email/emailTest','POST','','',''),(99,'p','888','/email/sendEmail','POST','','',''),(40,'p','888','/fileUploadAndDownload/breakpointContinue','POST','','',''),(39,'p','888','/fileUploadAndDownload/breakpointContinueFinish','POST','','',''),(43,'p','888','/fileUploadAndDownload/deleteFile','POST','','',''),(44,'p','888','/fileUploadAndDownload/editFileName','POST','','',''),(38,'p','888','/fileUploadAndDownload/findFile','GET','','',''),(45,'p','888','/fileUploadAndDownload/getFileList','POST','','',''),(46,'p','888','/fileUploadAndDownload/importURL','POST','','',''),(41,'p','888','/fileUploadAndDownload/removeChunk','POST','','',''),(42,'p','888','/fileUploadAndDownload/upload','POST','','',''),(115,'p','888','/info/createInfo','POST','','',''),(116,'p','888','/info/deleteInfo','DELETE','','',''),(117,'p','888','/info/deleteInfoByIds','DELETE','','',''),(119,'p','888','/info/findInfo','GET','','',''),(120,'p','888','/info/getInfoList','GET','','',''),(118,'p','888','/info/updateInfo','PUT','','',''),(49,'p','888','/jwt/jsonInBlacklist','POST','','',''),(21,'p','888','/menu/addBaseMenu','POST','','',''),(23,'p','888','/menu/addMenuAuthority','POST','','',''),(25,'p','888','/menu/deleteBaseMenu','POST','','',''),(27,'p','888','/menu/getBaseMenuById','POST','','',''),(22,'p','888','/menu/getBaseMenuTree','POST','','',''),(19,'p','888','/menu/getMenu','POST','','',''),(24,'p','888','/menu/getMenuAuthority','POST','','',''),(20,'p','888','/menu/getMenuList','POST','','',''),(26,'p','888','/menu/updateBaseMenu','POST','','',''),(101,'p','888','/simpleUploader/checkFileMd5','GET','','',''),(102,'p','888','/simpleUploader/mergeFileMd5','GET','','',''),(100,'p','888','/simpleUploader/upload','POST','','',''),(90,'p','888','/sysDictionary/createSysDictionary','POST','','',''),(91,'p','888','/sysDictionary/deleteSysDictionary','DELETE','','',''),(87,'p','888','/sysDictionary/findSysDictionary','GET','','',''),(89,'p','888','/sysDictionary/getSysDictionaryList','GET','','',''),(88,'p','888','/sysDictionary/updateSysDictionary','PUT','','',''),(80,'p','888','/sysDictionaryDetail/createSysDictionaryDetail','POST','','',''),(82,'p','888','/sysDictionaryDetail/deleteSysDictionaryDetail','DELETE','','',''),(78,'p','888','/sysDictionaryDetail/findSysDictionaryDetail','GET','','',''),(85,'p','888','/sysDictionaryDetail/getDictionaryDetailsByParent','GET','','',''),(86,'p','888','/sysDictionaryDetail/getDictionaryPath','GET','','',''),(83,'p','888','/sysDictionaryDetail/getDictionaryTreeList','GET','','',''),(84,'p','888','/sysDictionaryDetail/getDictionaryTreeListByType','GET','','',''),(81,'p','888','/sysDictionaryDetail/getSysDictionaryDetailList','GET','','',''),(79,'p','888','/sysDictionaryDetail/updateSysDictionaryDetail','PUT','','',''),(106,'p','888','/sysExportTemplate/createSysExportTemplate','POST','','',''),(107,'p','888','/sysExportTemplate/deleteSysExportTemplate','DELETE','','',''),(108,'p','888','/sysExportTemplate/deleteSysExportTemplateByIds','DELETE','','',''),(112,'p','888','/sysExportTemplate/exportExcel','GET','','',''),(113,'p','888','/sysExportTemplate/exportTemplate','GET','','',''),(110,'p','888','/sysExportTemplate/findSysExportTemplate','GET','','',''),(111,'p','888','/sysExportTemplate/getSysExportTemplateList','GET','','',''),(114,'p','888','/sysExportTemplate/importExcel','POST','','',''),(109,'p','888','/sysExportTemplate/updateSysExportTemplate','PUT','','',''),(94,'p','888','/sysOperationRecord/createSysOperationRecord','POST','','',''),(96,'p','888','/sysOperationRecord/deleteSysOperationRecord','DELETE','','',''),(97,'p','888','/sysOperationRecord/deleteSysOperationRecordByIds','DELETE','','',''),(92,'p','888','/sysOperationRecord/findSysOperationRecord','GET','','',''),(95,'p','888','/sysOperationRecord/getSysOperationRecordList','GET','','',''),(93,'p','888','/sysOperationRecord/updateSysOperationRecord','PUT','','',''),(121,'p','888','/sysParams/createSysParams','POST','','',''),(122,'p','888','/sysParams/deleteSysParams','DELETE','','',''),(123,'p','888','/sysParams/deleteSysParamsByIds','DELETE','','',''),(125,'p','888','/sysParams/findSysParams','GET','','',''),(127,'p','888','/sysParams/getSysParam','GET','','',''),(126,'p','888','/sysParams/getSysParamsList','GET','','',''),(124,'p','888','/sysParams/updateSysParams','PUT','','',''),(52,'p','888','/system/getServerInfo','POST','','',''),(50,'p','888','/system/getSystemConfig','POST','','',''),(51,'p','888','/system/setSystemConfig','POST','','',''),(136,'p','888','/sysVersion/deleteSysVersion','DELETE','','',''),(137,'p','888','/sysVersion/deleteSysVersionByIds','DELETE','','',''),(133,'p','888','/sysVersion/downloadVersionJson','GET','','',''),(134,'p','888','/sysVersion/exportVersion','POST','','',''),(131,'p','888','/sysVersion/findSysVersion','GET','','',''),(132,'p','888','/sysVersion/getSysVersionList','GET','','',''),(135,'p','888','/sysVersion/importVersion','POST','','',''),(1,'p','888','/user/admin_register','POST','','',''),(33,'p','888','/user/changePassword','POST','','',''),(32,'p','888','/user/deleteUser','DELETE','','',''),(28,'p','888','/user/getUserInfo','GET','','',''),(31,'p','888','/user/getUserList','POST','','',''),(36,'p','888','/user/resetPassword','POST','','',''),(30,'p','888','/user/setSelfInfo','PUT','','',''),(37,'p','888','/user/setSelfSetting','PUT','','',''),(35,'p','888','/user/setUserAuthorities','POST','','',''),(34,'p','888','/user/setUserAuthority','POST','','',''),(29,'p','888','/user/setUserInfo','PUT','','',''),(139,'p','8881','/api/createApi','POST','','',''),(142,'p','8881','/api/deleteApi','POST','','',''),(144,'p','8881','/api/getAllApis','POST','','',''),(141,'p','8881','/api/getApiById','POST','','',''),(140,'p','8881','/api/getApiList','POST','','',''),(143,'p','8881','/api/updateApi','POST','','',''),(145,'p','8881','/authority/createAuthority','POST','','',''),(146,'p','8881','/authority/deleteAuthority','POST','','',''),(147,'p','8881','/authority/getAuthorityList','POST','','',''),(148,'p','8881','/authority/setDataAuthority','POST','','',''),(167,'p','8881','/casbin/getPolicyPathByAuthorityId','POST','','',''),(166,'p','8881','/casbin/updateCasbin','POST','','',''),(173,'p','8881','/customer/customer','DELETE','','',''),(174,'p','8881','/customer/customer','GET','','',''),(171,'p','8881','/customer/customer','POST','','',''),(172,'p','8881','/customer/customer','PUT','','',''),(175,'p','8881','/customer/customerList','GET','','',''),(163,'p','8881','/fileUploadAndDownload/deleteFile','POST','','',''),(164,'p','8881','/fileUploadAndDownload/editFileName','POST','','',''),(162,'p','8881','/fileUploadAndDownload/getFileList','POST','','',''),(165,'p','8881','/fileUploadAndDownload/importURL','POST','','',''),(161,'p','8881','/fileUploadAndDownload/upload','POST','','',''),(168,'p','8881','/jwt/jsonInBlacklist','POST','','',''),(151,'p','8881','/menu/addBaseMenu','POST','','',''),(153,'p','8881','/menu/addMenuAuthority','POST','','',''),(155,'p','8881','/menu/deleteBaseMenu','POST','','',''),(157,'p','8881','/menu/getBaseMenuById','POST','','',''),(152,'p','8881','/menu/getBaseMenuTree','POST','','',''),(149,'p','8881','/menu/getMenu','POST','','',''),(154,'p','8881','/menu/getMenuAuthority','POST','','',''),(150,'p','8881','/menu/getMenuList','POST','','',''),(156,'p','8881','/menu/updateBaseMenu','POST','','',''),(169,'p','8881','/system/getSystemConfig','POST','','',''),(170,'p','8881','/system/setSystemConfig','POST','','',''),(138,'p','8881','/user/admin_register','POST','','',''),(158,'p','8881','/user/changePassword','POST','','',''),(176,'p','8881','/user/getUserInfo','GET','','',''),(159,'p','8881','/user/getUserList','POST','','',''),(160,'p','8881','/user/setUserAuthority','POST','','',''),(178,'p','9528','/api/createApi','POST','','',''),(181,'p','9528','/api/deleteApi','POST','','',''),(183,'p','9528','/api/getAllApis','POST','','',''),(180,'p','9528','/api/getApiById','POST','','',''),(179,'p','9528','/api/getApiList','POST','','',''),(182,'p','9528','/api/updateApi','POST','','',''),(184,'p','9528','/authority/createAuthority','POST','','',''),(185,'p','9528','/authority/deleteAuthority','POST','','',''),(186,'p','9528','/authority/getAuthorityList','POST','','',''),(187,'p','9528','/authority/setDataAuthority','POST','','',''),(215,'p','9528','/autoCode/createTemp','POST','','',''),(206,'p','9528','/casbin/getPolicyPathByAuthorityId','POST','','',''),(205,'p','9528','/casbin/updateCasbin','POST','','',''),(213,'p','9528','/customer/customer','DELETE','','',''),(211,'p','9528','/customer/customer','GET','','',''),(212,'p','9528','/customer/customer','POST','','',''),(210,'p','9528','/customer/customer','PUT','','',''),(214,'p','9528','/customer/customerList','GET','','',''),(202,'p','9528','/fileUploadAndDownload/deleteFile','POST','','',''),(203,'p','9528','/fileUploadAndDownload/editFileName','POST','','',''),(201,'p','9528','/fileUploadAndDownload/getFileList','POST','','',''),(204,'p','9528','/fileUploadAndDownload/importURL','POST','','',''),(200,'p','9528','/fileUploadAndDownload/upload','POST','','',''),(207,'p','9528','/jwt/jsonInBlacklist','POST','','',''),(190,'p','9528','/menu/addBaseMenu','POST','','',''),(192,'p','9528','/menu/addMenuAuthority','POST','','',''),(194,'p','9528','/menu/deleteBaseMenu','POST','','',''),(196,'p','9528','/menu/getBaseMenuById','POST','','',''),(191,'p','9528','/menu/getBaseMenuTree','POST','','',''),(188,'p','9528','/menu/getMenu','POST','','',''),(193,'p','9528','/menu/getMenuAuthority','POST','','',''),(189,'p','9528','/menu/getMenuList','POST','','',''),(195,'p','9528','/menu/updateBaseMenu','POST','','',''),(208,'p','9528','/system/getSystemConfig','POST','','',''),(209,'p','9528','/system/setSystemConfig','POST','','',''),(177,'p','9528','/user/admin_register','POST','','',''),(197,'p','9528','/user/changePassword','POST','','',''),(216,'p','9528','/user/getUserInfo','GET','','',''),(198,'p','9528','/user/getUserList','POST','','',''),(199,'p','9528','/user/setUserAuthority','POST','','','');
/*!40000 ALTER TABLE `casbin_rule` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `gva_announcements_info`
--

DROP TABLE IF EXISTS `gva_announcements_info`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `gva_announcements_info` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `title` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '公告标题',
  `content` text COLLATE utf8mb4_general_ci COMMENT '公告内容',
  `user_id` bigint DEFAULT NULL COMMENT '发布者',
  `attachments` json DEFAULT NULL COMMENT '相关附件',
  PRIMARY KEY (`id`),
  KEY `idx_gva_announcements_info_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `gva_announcements_info`
--

LOCK TABLES `gva_announcements_info` WRITE;
/*!40000 ALTER TABLE `gva_announcements_info` DISABLE KEYS */;
/*!40000 ALTER TABLE `gva_announcements_info` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `jwt_blacklists`
--

DROP TABLE IF EXISTS `jwt_blacklists`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `jwt_blacklists` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `jwt` text COLLATE utf8mb4_general_ci COMMENT 'jwt',
  PRIMARY KEY (`id`),
  KEY `idx_jwt_blacklists_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `jwt_blacklists`
--

LOCK TABLES `jwt_blacklists` WRITE;
/*!40000 ALTER TABLE `jwt_blacklists` DISABLE KEYS */;
INSERT INTO `jwt_blacklists` VALUES (1,'2026-04-21 15:04:10.222','2026-04-21 15:04:10.222',NULL,'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVVUlEIjoiNmE3NTQ3YzktN2VlOC00YWNiLTlmOTEtOWNmOGU4ZmNhYTZlIiwiSUQiOjEsIlVzZXJuYW1lIjoiYWRtaW4iLCJOaWNrTmFtZSI6Ik1yLuWlh-a3vCIsIkF1dGhvcml0eUlkIjo4ODgsIkJ1ZmZlclRpbWUiOjg2NDAwLCJpc3MiOiJxbVBsdXMiLCJhdWQiOlsiR1ZBIl0sImV4cCI6MTc3NzM2MjA0NSwibmJmIjoxNzc2NzU3MjQ1fQ.mZsbjfeGko4toXb2KiZ410qYJVcsaMMq7o8HaVkGGzA');
/*!40000 ALTER TABLE `jwt_blacklists` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ledger_account_group`
--

DROP TABLE IF EXISTS `ledger_account_group`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ledger_account_group` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `account_group` text COLLATE utf8mb4_general_ci,
  `title` varchar(256) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '标题',
  PRIMARY KEY (`id`),
  KEY `idx_ledger_account_group_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ledger_account_group`
--

LOCK TABLES `ledger_account_group` WRITE;
/*!40000 ALTER TABLE `ledger_account_group` DISABLE KEYS */;
/*!40000 ALTER TABLE `ledger_account_group` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ledger_sessions`
--

DROP TABLE IF EXISTS `ledger_sessions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ledger_sessions` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `bot_id` bigint DEFAULT NULL,
  `chat_group_id` bigint DEFAULT NULL,
  `work_date` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `is_active` bigint DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `idx_ledger_sessions_deleted_at` (`deleted_at`),
  KEY `idx_ledger_sessions_bot_id` (`bot_id`),
  KEY `idx_ledger_sessions_chat_group_id` (`chat_group_id`),
  KEY `idx_ledger_sessions_work_date` (`work_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ledger_sessions`
--

LOCK TABLES `ledger_sessions` WRITE;
/*!40000 ALTER TABLE `ledger_sessions` DISABLE KEYS */;
/*!40000 ALTER TABLE `ledger_sessions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_apis`
--

DROP TABLE IF EXISTS `sys_apis`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_apis` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `path` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'api路径',
  `description` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'api中文描述',
  `api_group` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'api组',
  `method` varchar(191) COLLATE utf8mb4_general_ci DEFAULT 'POST' COMMENT '方法',
  PRIMARY KEY (`id`),
  KEY `idx_sys_apis_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=136 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_apis`
--

LOCK TABLES `sys_apis` WRITE;
/*!40000 ALTER TABLE `sys_apis` DISABLE KEYS */;
INSERT INTO `sys_apis` VALUES (1,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/jwt/jsonInBlacklist','jwt加入黑名单(退出，必选)','jwt','POST'),(2,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/user/deleteUser','删除用户','系统用户','DELETE'),(3,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/user/admin_register','用户注册','系统用户','POST'),(4,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/user/getUserList','获取用户列表','系统用户','POST'),(5,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/user/setUserInfo','设置用户信息','系统用户','PUT'),(6,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/user/setSelfInfo','设置自身信息(必选)','系统用户','PUT'),(7,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/user/getUserInfo','获取自身信息(必选)','系统用户','GET'),(8,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/user/setUserAuthorities','设置权限组','系统用户','POST'),(9,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/user/changePassword','修改密码（建议选择)','系统用户','POST'),(10,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/user/setUserAuthority','修改用户角色(必选)','系统用户','POST'),(11,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/user/resetPassword','重置用户密码','系统用户','POST'),(12,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/user/setSelfSetting','用户界面配置','系统用户','PUT'),(13,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/api/createApi','创建api','api','POST'),(14,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/api/deleteApi','删除Api','api','POST'),(15,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/api/updateApi','更新Api','api','POST'),(16,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/api/getApiList','获取api列表','api','POST'),(17,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/api/getAllApis','获取所有api','api','POST'),(18,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/api/getApiById','获取api详细信息','api','POST'),(19,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/api/deleteApisByIds','批量删除api','api','DELETE'),(20,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/api/syncApi','获取待同步API','api','GET'),(21,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/api/getApiGroups','获取路由组','api','GET'),(22,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/api/enterSyncApi','确认同步API','api','POST'),(23,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/api/ignoreApi','忽略API','api','POST'),(24,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/authority/copyAuthority','拷贝角色','角色','POST'),(25,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/authority/createAuthority','创建角色','角色','POST'),(26,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/authority/deleteAuthority','删除角色','角色','POST'),(27,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/authority/updateAuthority','更新角色信息','角色','PUT'),(28,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/authority/getAuthorityList','获取角色列表','角色','POST'),(29,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/authority/setDataAuthority','设置角色资源权限','角色','POST'),(30,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/casbin/updateCasbin','更改角色api权限','casbin','POST'),(31,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/casbin/getPolicyPathByAuthorityId','获取权限列表','casbin','POST'),(32,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/menu/addBaseMenu','新增菜单','菜单','POST'),(33,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/menu/getMenu','获取菜单树(必选)','菜单','POST'),(34,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/menu/deleteBaseMenu','删除菜单','菜单','POST'),(35,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/menu/updateBaseMenu','更新菜单','菜单','POST'),(36,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/menu/getBaseMenuById','根据id获取菜单','菜单','POST'),(37,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/menu/getMenuList','分页获取基础menu列表','菜单','POST'),(38,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/menu/getBaseMenuTree','获取用户动态路由','菜单','POST'),(39,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/menu/getMenuAuthority','获取指定角色menu','菜单','POST'),(40,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/menu/addMenuAuthority','增加menu和角色关联关系','菜单','POST'),(41,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/fileUploadAndDownload/findFile','寻找目标文件（秒传）','分片上传','GET'),(42,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/fileUploadAndDownload/breakpointContinue','断点续传','分片上传','POST'),(43,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/fileUploadAndDownload/breakpointContinueFinish','断点续传完成','分片上传','POST'),(44,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/fileUploadAndDownload/removeChunk','上传完成移除文件','分片上传','POST'),(45,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/fileUploadAndDownload/upload','文件上传（建议选择）','文件上传与下载','POST'),(46,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/fileUploadAndDownload/deleteFile','删除文件','文件上传与下载','POST'),(47,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/fileUploadAndDownload/editFileName','文件名或者备注编辑','文件上传与下载','POST'),(48,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/fileUploadAndDownload/getFileList','获取上传文件列表','文件上传与下载','POST'),(49,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/fileUploadAndDownload/importURL','导入URL','文件上传与下载','POST'),(50,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/system/getServerInfo','获取服务器信息','系统服务','POST'),(51,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/system/getSystemConfig','获取配置文件内容','系统服务','POST'),(52,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/system/setSystemConfig','设置配置文件内容','系统服务','POST'),(53,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/customer/customer','更新客户','客户','PUT'),(54,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/customer/customer','创建客户','客户','POST'),(55,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/customer/customer','删除客户','客户','DELETE'),(56,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/customer/customer','获取单一客户','客户','GET'),(57,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/customer/customerList','获取客户列表','客户','GET'),(58,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/autoCode/getDB','获取所有数据库','代码生成器','GET'),(59,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/autoCode/getTables','获取数据库表','代码生成器','GET'),(60,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/autoCode/createTemp','自动化代码','代码生成器','POST'),(61,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/autoCode/preview','预览自动化代码','代码生成器','POST'),(62,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/autoCode/getColumn','获取所选table的所有字段','代码生成器','GET'),(63,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/autoCode/installPlugin','安装插件','代码生成器','POST'),(64,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/autoCode/pubPlug','打包插件','代码生成器','POST'),(65,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/autoCode/mcp','自动生成 MCP Tool 模板','代码生成器','POST'),(66,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/autoCode/mcpTest','MCP Tool 测试','代码生成器','POST'),(67,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/autoCode/mcpList','获取 MCP ToolList','代码生成器','POST'),(68,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/autoCode/createPackage','配置模板','模板配置','POST'),(69,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/autoCode/getTemplates','获取模板文件','模板配置','GET'),(70,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/autoCode/getPackage','获取所有模板','模板配置','POST'),(71,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/autoCode/delPackage','删除模板','模板配置','POST'),(72,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/autoCode/getMeta','获取meta信息','代码生成器历史','POST'),(73,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/autoCode/rollback','回滚自动生成代码','代码生成器历史','POST'),(74,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/autoCode/getSysHistory','查询回滚记录','代码生成器历史','POST'),(75,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/autoCode/delSysHistory','删除回滚记录','代码生成器历史','POST'),(76,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/autoCode/addFunc','增加模板方法','代码生成器历史','POST'),(77,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/sysDictionaryDetail/updateSysDictionaryDetail','更新字典内容','系统字典详情','PUT'),(78,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/sysDictionaryDetail/createSysDictionaryDetail','新增字典内容','系统字典详情','POST'),(79,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/sysDictionaryDetail/deleteSysDictionaryDetail','删除字典内容','系统字典详情','DELETE'),(80,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/sysDictionaryDetail/findSysDictionaryDetail','根据ID获取字典内容','系统字典详情','GET'),(81,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/sysDictionaryDetail/getSysDictionaryDetailList','获取字典内容列表','系统字典详情','GET'),(82,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/sysDictionaryDetail/getDictionaryTreeList','获取字典数列表','系统字典详情','GET'),(83,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/sysDictionaryDetail/getDictionaryTreeListByType','根据分类获取字典数列表','系统字典详情','GET'),(84,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/sysDictionaryDetail/getDictionaryDetailsByParent','根据父级ID获取字典详情','系统字典详情','GET'),(85,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/sysDictionaryDetail/getDictionaryPath','获取字典详情的完整路径','系统字典详情','GET'),(86,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/sysDictionary/createSysDictionary','新增字典','系统字典','POST'),(87,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/sysDictionary/deleteSysDictionary','删除字典','系统字典','DELETE'),(88,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/sysDictionary/updateSysDictionary','更新字典','系统字典','PUT'),(89,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/sysDictionary/findSysDictionary','根据ID获取字典（建议选择）','系统字典','GET'),(90,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/sysDictionary/getSysDictionaryList','获取字典列表','系统字典','GET'),(91,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/sysOperationRecord/createSysOperationRecord','新增操作记录','操作记录','POST'),(92,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/sysOperationRecord/findSysOperationRecord','根据ID获取操作记录','操作记录','GET'),(93,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/sysOperationRecord/getSysOperationRecordList','获取操作记录列表','操作记录','GET'),(94,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/sysOperationRecord/deleteSysOperationRecord','删除操作记录','操作记录','DELETE'),(95,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/sysOperationRecord/deleteSysOperationRecordByIds','批量删除操作历史','操作记录','DELETE'),(96,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/simpleUploader/upload','插件版分片上传','断点续传(插件版)','POST'),(97,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/simpleUploader/checkFileMd5','文件完整度验证','断点续传(插件版)','GET'),(98,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/simpleUploader/mergeFileMd5','上传完成合并文件','断点续传(插件版)','GET'),(99,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/email/emailTest','发送测试邮件','email','POST'),(100,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/email/sendEmail','发送邮件','email','POST'),(101,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/authorityBtn/setAuthorityBtn','设置按钮权限','按钮权限','POST'),(102,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/authorityBtn/getAuthorityBtn','获取已有按钮权限','按钮权限','POST'),(103,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/authorityBtn/canRemoveAuthorityBtn','删除按钮','按钮权限','POST'),(104,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/sysExportTemplate/createSysExportTemplate','新增导出模板','导出模板','POST'),(105,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/sysExportTemplate/deleteSysExportTemplate','删除导出模板','导出模板','DELETE'),(106,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/sysExportTemplate/deleteSysExportTemplateByIds','批量删除导出模板','导出模板','DELETE'),(107,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/sysExportTemplate/updateSysExportTemplate','更新导出模板','导出模板','PUT'),(108,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/sysExportTemplate/findSysExportTemplate','根据ID获取导出模板','导出模板','GET'),(109,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/sysExportTemplate/getSysExportTemplateList','获取导出模板列表','导出模板','GET'),(110,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/sysExportTemplate/exportExcel','导出Excel','导出模板','GET'),(111,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/sysExportTemplate/exportTemplate','下载模板','导出模板','GET'),(112,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/sysExportTemplate/importExcel','导入Excel','导出模板','POST'),(113,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/info/createInfo','新建公告','公告','POST'),(114,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/info/deleteInfo','删除公告','公告','DELETE'),(115,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/info/deleteInfoByIds','批量删除公告','公告','DELETE'),(116,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/info/updateInfo','更新公告','公告','PUT'),(117,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/info/findInfo','根据ID获取公告','公告','GET'),(118,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/info/getInfoList','获取公告列表','公告','GET'),(119,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/sysParams/createSysParams','新建参数','参数管理','POST'),(120,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/sysParams/deleteSysParams','删除参数','参数管理','DELETE'),(121,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/sysParams/deleteSysParamsByIds','批量删除参数','参数管理','DELETE'),(122,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/sysParams/updateSysParams','更新参数','参数管理','PUT'),(123,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/sysParams/findSysParams','根据ID获取参数','参数管理','GET'),(124,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/sysParams/getSysParamsList','获取参数列表','参数管理','GET'),(125,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/sysParams/getSysParam','获取参数列表','参数管理','GET'),(126,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/attachmentCategory/getCategoryList','分类列表','媒体库分类','GET'),(127,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/attachmentCategory/addCategory','添加/编辑分类','媒体库分类','POST'),(128,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/attachmentCategory/deleteCategory','删除分类','媒体库分类','POST'),(129,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/sysVersion/findSysVersion','获取单一版本','版本控制','GET'),(130,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/sysVersion/getSysVersionList','获取版本列表','版本控制','GET'),(131,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/sysVersion/downloadVersionJson','下载版本json','版本控制','GET'),(132,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/sysVersion/exportVersion','创建版本','版本控制','POST'),(133,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/sysVersion/importVersion','同步版本','版本控制','POST'),(134,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/sysVersion/deleteSysVersion','删除版本','版本控制','DELETE'),(135,'2026-04-21 14:16:22.088','2026-04-21 14:16:22.088',NULL,'/sysVersion/deleteSysVersionByIds','批量删除版本','版本控制','DELETE');
/*!40000 ALTER TABLE `sys_apis` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_authorities`
--

DROP TABLE IF EXISTS `sys_authorities`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_authorities` (
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `authority_id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '角色ID',
  `authority_name` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '角色名',
  `parent_id` bigint unsigned DEFAULT NULL COMMENT '父角色ID',
  `default_router` varchar(191) COLLATE utf8mb4_general_ci DEFAULT 'dashboard' COMMENT '默认菜单',
  PRIMARY KEY (`authority_id`),
  UNIQUE KEY `uni_sys_authorities_authority_id` (`authority_id`)
) ENGINE=InnoDB AUTO_INCREMENT=9529 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_authorities`
--

LOCK TABLES `sys_authorities` WRITE;
/*!40000 ALTER TABLE `sys_authorities` DISABLE KEYS */;
INSERT INTO `sys_authorities` VALUES ('2026-04-21 14:16:23.902','2026-04-21 15:32:06.137',NULL,888,'普通用户',0,'user'),('2026-04-21 14:16:23.902','2026-04-21 14:16:33.240',NULL,8881,'普通用户子角色',888,'dashboard'),('2026-04-21 14:16:23.902','2026-04-21 14:16:32.602',NULL,9528,'测试角色',0,'dashboard');
/*!40000 ALTER TABLE `sys_authorities` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_authority_btns`
--

DROP TABLE IF EXISTS `sys_authority_btns`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_authority_btns` (
  `authority_id` bigint unsigned DEFAULT NULL COMMENT '角色ID',
  `sys_menu_id` bigint unsigned DEFAULT NULL COMMENT '菜单ID',
  `sys_base_menu_btn_id` bigint unsigned DEFAULT NULL COMMENT '菜单按钮ID'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_authority_btns`
--

LOCK TABLES `sys_authority_btns` WRITE;
/*!40000 ALTER TABLE `sys_authority_btns` DISABLE KEYS */;
/*!40000 ALTER TABLE `sys_authority_btns` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_authority_menus`
--

DROP TABLE IF EXISTS `sys_authority_menus`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_authority_menus` (
  `sys_base_menu_id` bigint unsigned NOT NULL,
  `sys_authority_authority_id` bigint unsigned NOT NULL COMMENT '角色ID',
  PRIMARY KEY (`sys_base_menu_id`,`sys_authority_authority_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_authority_menus`
--

LOCK TABLES `sys_authority_menus` WRITE;
/*!40000 ALTER TABLE `sys_authority_menus` DISABLE KEYS */;
INSERT INTO `sys_authority_menus` VALUES (1,888),(1,8881),(1,9528),(3,888),(3,8881),(4,888),(4,8881),(4,9528),(6,888),(6,8881),(10,888),(11,888),(12,888),(13,888),(14,888),(15,888),(16,888),(20,888),(20,8881),(21,888),(21,8881),(22,888),(22,8881),(23,888),(23,8881),(24,888),(24,8881),(25,888),(25,8881),(26,888),(26,8881),(27,888),(27,8881),(28,888),(28,8881),(29,888),(29,8881),(30,888),(30,8881);
/*!40000 ALTER TABLE `sys_authority_menus` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_auto_code_histories`
--

DROP TABLE IF EXISTS `sys_auto_code_histories`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_auto_code_histories` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `table_name` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '表名',
  `package` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '模块名/插件名',
  `request` text COLLATE utf8mb4_general_ci COMMENT '前端传入的结构化信息',
  `struct_name` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '结构体名称',
  `abbreviation` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '结构体名称缩写',
  `business_db` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '业务库',
  `description` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'Struct中文名称',
  `templates` text COLLATE utf8mb4_general_ci COMMENT '模板信息',
  `Injections` text COLLATE utf8mb4_general_ci COMMENT '注入路径',
  `flag` bigint DEFAULT NULL COMMENT '[0:创建,1:回滚]',
  `api_ids` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'api表注册内容',
  `menu_id` bigint unsigned DEFAULT NULL COMMENT '菜单ID',
  `export_template_id` bigint unsigned DEFAULT NULL COMMENT '导出模板ID',
  `package_id` bigint unsigned DEFAULT NULL COMMENT '包ID',
  PRIMARY KEY (`id`),
  KEY `idx_sys_auto_code_histories_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_auto_code_histories`
--

LOCK TABLES `sys_auto_code_histories` WRITE;
/*!40000 ALTER TABLE `sys_auto_code_histories` DISABLE KEYS */;
/*!40000 ALTER TABLE `sys_auto_code_histories` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_auto_code_packages`
--

DROP TABLE IF EXISTS `sys_auto_code_packages`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_auto_code_packages` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `desc` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '描述',
  `label` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '展示名',
  `template` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '模版',
  `package_name` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '包名',
  `module` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_sys_auto_code_packages_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_auto_code_packages`
--

LOCK TABLES `sys_auto_code_packages` WRITE;
/*!40000 ALTER TABLE `sys_auto_code_packages` DISABLE KEYS */;
/*!40000 ALTER TABLE `sys_auto_code_packages` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_base_menu_btns`
--

DROP TABLE IF EXISTS `sys_base_menu_btns`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_base_menu_btns` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '按钮关键key',
  `desc` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `sys_base_menu_id` bigint unsigned DEFAULT NULL COMMENT '菜单ID',
  PRIMARY KEY (`id`),
  KEY `idx_sys_base_menu_btns_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_base_menu_btns`
--

LOCK TABLES `sys_base_menu_btns` WRITE;
/*!40000 ALTER TABLE `sys_base_menu_btns` DISABLE KEYS */;
/*!40000 ALTER TABLE `sys_base_menu_btns` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_base_menu_parameters`
--

DROP TABLE IF EXISTS `sys_base_menu_parameters`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_base_menu_parameters` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `sys_base_menu_id` bigint unsigned DEFAULT NULL,
  `type` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '地址栏携带参数为params还是query',
  `key` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '地址栏携带参数的key',
  `value` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '地址栏携带参数的值',
  PRIMARY KEY (`id`),
  KEY `idx_sys_base_menu_parameters_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_base_menu_parameters`
--

LOCK TABLES `sys_base_menu_parameters` WRITE;
/*!40000 ALTER TABLE `sys_base_menu_parameters` DISABLE KEYS */;
/*!40000 ALTER TABLE `sys_base_menu_parameters` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_base_menus`
--

DROP TABLE IF EXISTS `sys_base_menus`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_base_menus` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `menu_level` bigint unsigned DEFAULT NULL,
  `parent_id` bigint unsigned DEFAULT NULL COMMENT '父菜单ID',
  `path` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '路由path',
  `name` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '路由name',
  `hidden` tinyint(1) DEFAULT NULL COMMENT '是否在列表隐藏',
  `component` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '对应前端文件路径',
  `sort` bigint DEFAULT NULL COMMENT '排序标记',
  `active_name` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '附加属性',
  `keep_alive` tinyint(1) DEFAULT NULL COMMENT '附加属性',
  `default_menu` tinyint(1) DEFAULT NULL COMMENT '附加属性',
  `title` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '附加属性',
  `icon` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '附加属性',
  `close_tab` tinyint(1) DEFAULT NULL COMMENT '附加属性',
  `transition_type` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '附加属性',
  PRIMARY KEY (`id`),
  KEY `idx_sys_base_menus_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=37 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_base_menus`
--

LOCK TABLES `sys_base_menus` WRITE;
/*!40000 ALTER TABLE `sys_base_menus` DISABLE KEYS */;
INSERT INTO `sys_base_menus` VALUES (1,'2026-04-21 14:16:29.252','2026-04-21 14:16:29.252',NULL,0,0,'dashboard','dashboard',0,'view/dashboard/index.vue',1,'',0,0,'仪表盘','odometer',0,''),(2,'2026-04-21 14:16:29.252','2026-04-21 14:16:29.252','2026-04-21 15:29:14.093',0,0,'about','about',0,'view/about/index.vue',9,'',0,0,'关于我们','info-filled',0,''),(3,'2026-04-21 14:16:29.252','2026-04-21 14:16:29.252',NULL,0,0,'admin','superAdmin',0,'view/superAdmin/index.vue',3,'',0,0,'超级管理员','user',0,''),(4,'2026-04-21 14:16:29.252','2026-04-21 14:16:29.252',NULL,0,0,'person','person',1,'view/person/person.vue',4,'',0,0,'个人信息','message',0,''),(5,'2026-04-21 14:16:29.252','2026-04-21 14:16:29.252','2026-04-21 15:26:16.164',0,0,'example','example',0,'view/example/index.vue',7,'',0,0,'示例文件','management',0,''),(6,'2026-04-21 14:16:29.252','2026-04-21 14:16:29.252',NULL,0,0,'systemTools','systemTools',0,'view/systemTools/index.vue',5,'',0,0,'系统工具','tools',0,''),(7,'2026-04-21 14:16:29.252','2026-04-21 14:16:29.252','2026-04-21 15:33:29.887',0,0,'https://www.gin-vue-admin.com','https://www.gin-vue-admin.com',0,'/',0,'',0,0,'官方网站','customer-gva',0,''),(8,'2026-04-21 14:16:29.252','2026-04-21 14:16:29.252','2026-04-21 15:29:32.337',0,0,'state','state',0,'view/system/state.vue',8,'',0,0,'服务器状态','cloudy',0,''),(9,'2026-04-21 14:16:29.252','2026-04-21 14:16:29.252','2026-04-21 15:31:04.954',0,0,'plugin','plugin',0,'view/routerHolder.vue',6,'',0,0,'插件系统','cherry',0,''),(10,'2026-04-21 14:16:29.463','2026-04-21 14:16:29.463',NULL,1,3,'authority','authority',0,'view/superAdmin/authority/authority.vue',1,'',0,0,'角色管理','avatar',0,''),(11,'2026-04-21 14:16:29.463','2026-04-21 14:16:29.463',NULL,1,3,'menu','menu',0,'view/superAdmin/menu/menu.vue',2,'',1,0,'菜单管理','tickets',0,''),(12,'2026-04-21 14:16:29.463','2026-04-21 14:16:29.463',NULL,1,3,'api','api',0,'view/superAdmin/api/api.vue',3,'',1,0,'api管理','platform',0,''),(13,'2026-04-21 14:16:29.463','2026-04-21 14:16:29.463',NULL,1,3,'user','user',0,'view/superAdmin/user/user.vue',4,'',0,0,'用户管理','coordinate',0,''),(14,'2026-04-21 14:16:29.463','2026-04-21 14:16:29.463',NULL,1,3,'dictionary','dictionary',0,'view/superAdmin/dictionary/sysDictionary.vue',5,'',0,0,'字典管理','notebook',0,''),(15,'2026-04-21 14:16:29.463','2026-04-21 14:16:29.463',NULL,1,3,'operation','operation',0,'view/superAdmin/operation/sysOperationRecord.vue',6,'',0,0,'操作历史','pie-chart',0,''),(16,'2026-04-21 14:16:29.463','2026-04-21 14:16:29.463',NULL,1,3,'sysParams','sysParams',0,'view/superAdmin/params/sysParams.vue',7,'',0,0,'参数管理','compass',0,''),(17,'2026-04-21 14:16:29.463','2026-04-21 14:16:29.463','2026-04-21 15:26:10.312',1,5,'upload','upload',0,'view/example/upload/upload.vue',5,'',0,0,'媒体库（上传下载）','upload',0,''),(18,'2026-04-21 14:16:29.463','2026-04-21 14:16:29.463','2026-04-21 15:26:06.142',1,5,'breakpoint','breakpoint',0,'view/example/breakpoint/breakpoint.vue',6,'',0,0,'断点续传','upload-filled',0,''),(19,'2026-04-21 14:16:29.463','2026-04-21 14:16:29.463','2026-04-21 15:26:01.515',1,5,'customer','customer',0,'view/example/customer/customer.vue',7,'',0,0,'客户列表（资源示例）','avatar',0,''),(20,'2026-04-21 14:16:29.463','2026-04-21 14:16:29.463',NULL,1,6,'autoCode','autoCode',0,'view/systemTools/autoCode/index.vue',1,'',1,0,'代码生成器','cpu',0,''),(21,'2026-04-21 14:16:29.463','2026-04-21 14:16:29.463',NULL,1,6,'formCreate','formCreate',0,'view/systemTools/formCreate/index.vue',3,'',1,0,'表单生成器','magic-stick',0,''),(22,'2026-04-21 14:16:29.463','2026-04-21 14:16:29.463',NULL,1,6,'system','system',0,'view/systemTools/system/system.vue',4,'',0,0,'系统配置','operation',0,''),(23,'2026-04-21 14:16:29.463','2026-04-21 14:16:29.463',NULL,1,6,'autoCodeAdmin','autoCodeAdmin',0,'view/systemTools/autoCodeAdmin/index.vue',2,'',0,0,'自动化代码管理','magic-stick',0,''),(24,'2026-04-21 14:16:29.463','2026-04-21 14:16:29.463',NULL,1,6,'autoCodeEdit/:id','autoCodeEdit',1,'view/systemTools/autoCode/index.vue',0,'',0,0,'自动化代码-${id}','magic-stick',0,''),(25,'2026-04-21 14:16:29.463','2026-04-21 14:16:29.463',NULL,1,6,'autoPkg','autoPkg',0,'view/systemTools/autoPkg/autoPkg.vue',0,'',0,0,'模板配置','folder',0,''),(26,'2026-04-21 14:16:29.463','2026-04-21 14:16:29.463',NULL,1,6,'exportTemplate','exportTemplate',0,'view/systemTools/exportTemplate/exportTemplate.vue',5,'',0,0,'导出模板','reading',0,''),(27,'2026-04-21 14:16:29.463','2026-04-21 14:16:29.463',NULL,1,6,'picture','picture',0,'view/systemTools/autoCode/picture.vue',6,'',0,0,'AI页面绘制','picture-filled',0,''),(28,'2026-04-21 14:16:29.463','2026-04-21 14:16:29.463',NULL,1,6,'mcpTool','mcpTool',0,'view/systemTools/autoCode/mcp.vue',7,'',0,0,'Mcp Tools模板','magnet',0,''),(29,'2026-04-21 14:16:29.463','2026-04-21 14:16:29.463',NULL,1,6,'mcpTest','mcpTest',0,'view/systemTools/autoCode/mcpTest.vue',7,'',0,0,'Mcp Tools测试','partly-cloudy',0,''),(30,'2026-04-21 14:16:29.463','2026-04-21 14:16:29.463',NULL,1,6,'sysVersion','sysVersion',0,'view/systemTools/version/version.vue',8,'',0,0,'版本管理','server',0,''),(31,'2026-04-21 14:16:29.463','2026-04-21 14:16:29.463','2026-04-21 15:30:39.476',1,9,'https://plugin.gin-vue-admin.com/','https://plugin.gin-vue-admin.com/',0,'https://plugin.gin-vue-admin.com/',0,'',0,0,'插件市场','shop',0,''),(32,'2026-04-21 14:16:29.463','2026-04-21 14:16:29.463','2026-04-21 15:30:46.853',1,9,'installPlugin','installPlugin',0,'view/systemTools/installPlugin/index.vue',1,'',0,0,'插件安装','box',0,''),(33,'2026-04-21 14:16:29.463','2026-04-21 14:16:29.463','2026-04-21 15:30:51.003',1,9,'pubPlug','pubPlug',0,'view/systemTools/pubPlug/pubPlug.vue',3,'',0,0,'打包插件','files',0,''),(34,'2026-04-21 14:16:29.463','2026-04-21 14:16:29.463','2026-04-21 15:30:55.684',1,9,'plugin-email','plugin-email',0,'plugin/email/view/index.vue',4,'',0,0,'邮件插件','message',0,''),(35,'2026-04-21 14:16:29.463','2026-04-21 14:16:29.463','2026-04-21 15:31:00.636',1,9,'anInfo','anInfo',0,'plugin/announcement/view/info.vue',5,'',0,0,'公告管理[示例]','scaleToOriginal',0,''),(36,'2026-04-21 16:27:39.167','2026-04-21 16:27:39.167',NULL,0,24,'anInfo','anInfo',0,'plugin/announcement/view/info.vue',5,'',0,0,'公告管理','box',0,'');
/*!40000 ALTER TABLE `sys_base_menus` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_data_authority_id`
--

DROP TABLE IF EXISTS `sys_data_authority_id`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_data_authority_id` (
  `sys_authority_authority_id` bigint unsigned NOT NULL COMMENT '角色ID',
  `data_authority_id_authority_id` bigint unsigned NOT NULL COMMENT '角色ID',
  PRIMARY KEY (`sys_authority_authority_id`,`data_authority_id_authority_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_data_authority_id`
--

LOCK TABLES `sys_data_authority_id` WRITE;
/*!40000 ALTER TABLE `sys_data_authority_id` DISABLE KEYS */;
INSERT INTO `sys_data_authority_id` VALUES (888,888),(888,8881),(888,9528),(9528,8881),(9528,9528);
/*!40000 ALTER TABLE `sys_data_authority_id` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_dictionaries`
--

DROP TABLE IF EXISTS `sys_dictionaries`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_dictionaries` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '字典名（中）',
  `type` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '字典名（英）',
  `status` tinyint(1) DEFAULT NULL COMMENT '状态',
  `desc` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '描述',
  `parent_id` bigint unsigned DEFAULT NULL COMMENT '父级字典ID',
  PRIMARY KEY (`id`),
  KEY `idx_sys_dictionaries_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_dictionaries`
--

LOCK TABLES `sys_dictionaries` WRITE;
/*!40000 ALTER TABLE `sys_dictionaries` DISABLE KEYS */;
INSERT INTO `sys_dictionaries` VALUES (1,'2026-04-21 14:16:25.519','2026-04-21 14:16:25.945',NULL,'性别','gender',1,'性别字典',NULL),(2,'2026-04-21 14:16:25.519','2026-04-21 14:16:26.480',NULL,'数据库int类型','int',1,'int类型对应的数据库类型',NULL),(3,'2026-04-21 14:16:25.519','2026-04-21 14:16:27.013',NULL,'数据库时间日期类型','time.Time',1,'数据库时间日期类型',NULL),(4,'2026-04-21 14:16:25.519','2026-04-21 14:16:27.544',NULL,'数据库浮点型','float64',1,'数据库浮点型',NULL),(5,'2026-04-21 14:16:25.519','2026-04-21 14:16:28.077',NULL,'数据库字符串','string',1,'数据库字符串',NULL),(6,'2026-04-21 14:16:25.519','2026-04-21 14:16:28.613',NULL,'数据库bool类型','bool',1,'数据库bool类型',NULL);
/*!40000 ALTER TABLE `sys_dictionaries` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_dictionary_details`
--

DROP TABLE IF EXISTS `sys_dictionary_details`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_dictionary_details` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `label` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '展示值',
  `value` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '字典值',
  `extend` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '扩展值',
  `status` tinyint(1) DEFAULT NULL COMMENT '启用状态',
  `sort` bigint DEFAULT NULL COMMENT '排序标记',
  `sys_dictionary_id` bigint unsigned DEFAULT NULL COMMENT '关联标记',
  `parent_id` bigint unsigned DEFAULT NULL COMMENT '父级字典详情ID',
  `level` bigint DEFAULT NULL COMMENT '层级深度',
  `path` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '层级路径',
  PRIMARY KEY (`id`),
  KEY `idx_sys_dictionary_details_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=34 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_dictionary_details`
--

LOCK TABLES `sys_dictionary_details` WRITE;
/*!40000 ALTER TABLE `sys_dictionary_details` DISABLE KEYS */;
INSERT INTO `sys_dictionary_details` VALUES (1,'2026-04-21 14:16:26.052','2026-04-21 14:16:26.052',NULL,'男','1','',1,1,1,NULL,0,''),(2,'2026-04-21 14:16:26.052','2026-04-21 14:16:26.052',NULL,'女','2','',1,2,1,NULL,0,''),(3,'2026-04-21 14:16:26.587','2026-04-21 14:16:26.587',NULL,'smallint','1','mysql',1,1,2,NULL,0,''),(4,'2026-04-21 14:16:26.587','2026-04-21 14:16:26.587',NULL,'mediumint','2','mysql',1,2,2,NULL,0,''),(5,'2026-04-21 14:16:26.587','2026-04-21 14:16:26.587',NULL,'int','3','mysql',1,3,2,NULL,0,''),(6,'2026-04-21 14:16:26.587','2026-04-21 14:16:26.587',NULL,'bigint','4','mysql',1,4,2,NULL,0,''),(7,'2026-04-21 14:16:26.587','2026-04-21 14:16:26.587',NULL,'int2','5','pgsql',1,5,2,NULL,0,''),(8,'2026-04-21 14:16:26.587','2026-04-21 14:16:26.587',NULL,'int4','6','pgsql',1,6,2,NULL,0,''),(9,'2026-04-21 14:16:26.587','2026-04-21 14:16:26.587',NULL,'int6','7','pgsql',1,7,2,NULL,0,''),(10,'2026-04-21 14:16:26.587','2026-04-21 14:16:26.587',NULL,'int8','8','pgsql',1,8,2,NULL,0,''),(11,'2026-04-21 14:16:27.118','2026-04-21 14:16:27.118',NULL,'date','0','mysql',1,0,3,NULL,0,''),(12,'2026-04-21 14:16:27.118','2026-04-21 14:16:27.118',NULL,'time','1','mysql',1,1,3,NULL,0,''),(13,'2026-04-21 14:16:27.118','2026-04-21 14:16:27.118',NULL,'year','2','mysql',1,2,3,NULL,0,''),(14,'2026-04-21 14:16:27.118','2026-04-21 14:16:27.118',NULL,'datetime','3','mysql',1,3,3,NULL,0,''),(15,'2026-04-21 14:16:27.118','2026-04-21 14:16:27.118',NULL,'timestamp','5','mysql',1,5,3,NULL,0,''),(16,'2026-04-21 14:16:27.118','2026-04-21 14:16:27.118',NULL,'timestamptz','6','pgsql',1,5,3,NULL,0,''),(17,'2026-04-21 14:16:27.651','2026-04-21 14:16:27.651',NULL,'float','0','mysql',1,0,4,NULL,0,''),(18,'2026-04-21 14:16:27.651','2026-04-21 14:16:27.651',NULL,'double','1','mysql',1,1,4,NULL,0,''),(19,'2026-04-21 14:16:27.651','2026-04-21 14:16:27.651',NULL,'decimal','2','mysql',1,2,4,NULL,0,''),(20,'2026-04-21 14:16:27.651','2026-04-21 14:16:27.651',NULL,'numeric','3','pgsql',1,3,4,NULL,0,''),(21,'2026-04-21 14:16:27.651','2026-04-21 14:16:27.651',NULL,'smallserial','4','pgsql',1,4,4,NULL,0,''),(22,'2026-04-21 14:16:28.183','2026-04-21 14:16:28.183',NULL,'char','0','mysql',1,0,5,NULL,0,''),(23,'2026-04-21 14:16:28.183','2026-04-21 14:16:28.183',NULL,'varchar','1','mysql',1,1,5,NULL,0,''),(24,'2026-04-21 14:16:28.183','2026-04-21 14:16:28.183',NULL,'tinyblob','2','mysql',1,2,5,NULL,0,''),(25,'2026-04-21 14:16:28.183','2026-04-21 14:16:28.183',NULL,'tinytext','3','mysql',1,3,5,NULL,0,''),(26,'2026-04-21 14:16:28.183','2026-04-21 14:16:28.183',NULL,'text','4','mysql',1,4,5,NULL,0,''),(27,'2026-04-21 14:16:28.183','2026-04-21 14:16:28.183',NULL,'blob','5','mysql',1,5,5,NULL,0,''),(28,'2026-04-21 14:16:28.183','2026-04-21 14:16:28.183',NULL,'mediumblob','6','mysql',1,6,5,NULL,0,''),(29,'2026-04-21 14:16:28.183','2026-04-21 14:16:28.183',NULL,'mediumtext','7','mysql',1,7,5,NULL,0,''),(30,'2026-04-21 14:16:28.183','2026-04-21 14:16:28.183',NULL,'longblob','8','mysql',1,8,5,NULL,0,''),(31,'2026-04-21 14:16:28.183','2026-04-21 14:16:28.183',NULL,'longtext','9','mysql',1,9,5,NULL,0,''),(32,'2026-04-21 14:16:28.720','2026-04-21 14:16:28.720',NULL,'tinyint','1','mysql',1,0,6,NULL,0,''),(33,'2026-04-21 14:16:28.720','2026-04-21 14:16:28.720',NULL,'bool','2','pgsql',1,0,6,NULL,0,'');
/*!40000 ALTER TABLE `sys_dictionary_details` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_export_template_condition`
--

DROP TABLE IF EXISTS `sys_export_template_condition`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_export_template_condition` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `template_id` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '模板标识',
  `from` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '条件取的key',
  `column` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '作为查询条件的字段',
  `operator` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '操作符',
  PRIMARY KEY (`id`),
  KEY `idx_sys_export_template_condition_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_export_template_condition`
--

LOCK TABLES `sys_export_template_condition` WRITE;
/*!40000 ALTER TABLE `sys_export_template_condition` DISABLE KEYS */;
/*!40000 ALTER TABLE `sys_export_template_condition` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_export_template_join`
--

DROP TABLE IF EXISTS `sys_export_template_join`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_export_template_join` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `template_id` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '模板标识',
  `joins` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '关联',
  `table` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '关联表',
  `on` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '关联条件',
  PRIMARY KEY (`id`),
  KEY `idx_sys_export_template_join_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_export_template_join`
--

LOCK TABLES `sys_export_template_join` WRITE;
/*!40000 ALTER TABLE `sys_export_template_join` DISABLE KEYS */;
/*!40000 ALTER TABLE `sys_export_template_join` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_export_templates`
--

DROP TABLE IF EXISTS `sys_export_templates`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_export_templates` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `db_name` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '数据库名称',
  `name` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '模板名称',
  `table_name` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '表名称',
  `template_id` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '模板标识',
  `template_info` text COLLATE utf8mb4_general_ci,
  `limit` bigint DEFAULT NULL COMMENT '导出限制',
  `order` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '排序',
  PRIMARY KEY (`id`),
  KEY `idx_sys_export_templates_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_export_templates`
--

LOCK TABLES `sys_export_templates` WRITE;
/*!40000 ALTER TABLE `sys_export_templates` DISABLE KEYS */;
INSERT INTO `sys_export_templates` VALUES (1,'2026-04-21 14:16:31.532','2026-04-21 14:16:31.532',NULL,'','api','sys_apis','api','{\n\"path\":\"路径\",\n\"method\":\"方法（大写）\",\n\"description\":\"方法介绍\",\n\"api_group\":\"方法分组\"\n}',NULL,'');
/*!40000 ALTER TABLE `sys_export_templates` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_ignore_apis`
--

DROP TABLE IF EXISTS `sys_ignore_apis`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_ignore_apis` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `path` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'api路径',
  `method` varchar(191) COLLATE utf8mb4_general_ci DEFAULT 'POST' COMMENT '方法',
  PRIMARY KEY (`id`),
  KEY `idx_sys_ignore_apis_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_ignore_apis`
--

LOCK TABLES `sys_ignore_apis` WRITE;
/*!40000 ALTER TABLE `sys_ignore_apis` DISABLE KEYS */;
INSERT INTO `sys_ignore_apis` VALUES (1,'2026-04-21 14:16:23.210','2026-04-21 14:16:23.210',NULL,'/swagger/*any','GET'),(2,'2026-04-21 14:16:23.210','2026-04-21 14:16:23.210',NULL,'/api/freshCasbin','GET'),(3,'2026-04-21 14:16:23.210','2026-04-21 14:16:23.210',NULL,'/uploads/file/*filepath','GET'),(4,'2026-04-21 14:16:23.210','2026-04-21 14:16:23.210',NULL,'/health','GET'),(5,'2026-04-21 14:16:23.210','2026-04-21 14:16:23.210',NULL,'/uploads/file/*filepath','HEAD'),(6,'2026-04-21 14:16:23.210','2026-04-21 14:16:23.210',NULL,'/autoCode/llmAuto','POST'),(7,'2026-04-21 14:16:23.210','2026-04-21 14:16:23.210',NULL,'/system/reloadSystem','POST'),(8,'2026-04-21 14:16:23.210','2026-04-21 14:16:23.210',NULL,'/base/login','POST'),(9,'2026-04-21 14:16:23.210','2026-04-21 14:16:23.210',NULL,'/base/captcha','POST'),(10,'2026-04-21 14:16:23.210','2026-04-21 14:16:23.210',NULL,'/init/initdb','POST'),(11,'2026-04-21 14:16:23.210','2026-04-21 14:16:23.210',NULL,'/init/checkdb','POST'),(12,'2026-04-21 14:16:23.210','2026-04-21 14:16:23.210',NULL,'/info/getInfoDataSource','GET'),(13,'2026-04-21 14:16:23.210','2026-04-21 14:16:23.210',NULL,'/info/getInfoPublic','GET');
/*!40000 ALTER TABLE `sys_ignore_apis` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_operation_records`
--

DROP TABLE IF EXISTS `sys_operation_records`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_operation_records` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `ip` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '请求ip',
  `method` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '请求方法',
  `path` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '请求路径',
  `status` bigint DEFAULT NULL COMMENT '请求状态',
  `latency` bigint DEFAULT NULL COMMENT '延迟',
  `agent` text COLLATE utf8mb4_general_ci COMMENT '代理',
  `error_message` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '错误信息',
  `body` text COLLATE utf8mb4_general_ci COMMENT '请求Body',
  `resp` text COLLATE utf8mb4_general_ci COMMENT '响应Body',
  `user_id` bigint unsigned DEFAULT NULL COMMENT '用户id',
  PRIMARY KEY (`id`),
  KEY `idx_sys_operation_records_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=23 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_operation_records`
--

LOCK TABLES `sys_operation_records` WRITE;
/*!40000 ALTER TABLE `sys_operation_records` DISABLE KEYS */;
INSERT INTO `sys_operation_records` VALUES (1,'2026-04-21 14:18:03.874','2026-04-21 14:18:03.874',NULL,'127.0.0.1','POST','/user/setUserAuthorities',200,513673709,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36','','{\"ID\":1,\"authorityIds\":[888,8881,9528]}','{\"code\":0,\"data\":{},\"msg\":\"修改成功\"}',1),(2,'2026-04-21 14:18:09.858','2026-04-21 14:18:09.858',NULL,'127.0.0.1','POST','/user/setUserAuthorities',200,515328250,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36','','{\"ID\":1,\"authorityIds\":[888,8881,9528]}','{\"code\":0,\"data\":{},\"msg\":\"修改成功\"}',1),(3,'2026-04-21 14:33:59.286','2026-04-21 14:33:59.286',NULL,'127.0.0.1','POST','/user/setUserAuthority',200,728760584,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36','','{\"authorityId\":9528}','{\"code\":0,\"data\":{},\"msg\":\"修改成功\"}',1),(4,'2026-04-21 14:34:08.855','2026-04-21 14:34:08.855',NULL,'127.0.0.1','POST','/user/setUserAuthority',200,868121500,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36','','{\"authorityId\":888}','{\"code\":0,\"data\":{},\"msg\":\"修改成功\"}',1),(5,'2026-04-21 15:25:20.359','2026-04-21 15:25:20.359',NULL,'127.0.0.1','POST','/menu/deleteBaseMenu',200,109961084,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36','','{\"ID\":5}','{\"code\":7,\"data\":{},\"msg\":\"删除失败:此菜单存在子菜单不可删除\"}',1),(6,'2026-04-21 15:26:02.130','2026-04-21 15:26:02.130',NULL,'127.0.0.1','POST','/menu/deleteBaseMenu',200,926282458,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36','','{\"ID\":19}','{\"code\":0,\"data\":{},\"msg\":\"删除成功\"}',1),(7,'2026-04-21 15:26:06.764','2026-04-21 15:26:06.764',NULL,'127.0.0.1','POST','/menu/deleteBaseMenu',200,932384042,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36','','{\"ID\":18}','{\"code\":0,\"data\":{},\"msg\":\"删除成功\"}',1),(8,'2026-04-21 15:26:10.925','2026-04-21 15:26:10.925',NULL,'127.0.0.1','POST','/menu/deleteBaseMenu',200,923858208,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36','','{\"ID\":17}','{\"code\":0,\"data\":{},\"msg\":\"删除成功\"}',1),(9,'2026-04-21 15:26:17.138','2026-04-21 15:26:17.138',NULL,'127.0.0.1','POST','/menu/deleteBaseMenu',200,1058023000,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36','','{\"ID\":5}','{\"code\":0,\"data\":{},\"msg\":\"删除成功\"}',1),(10,'2026-04-21 15:29:14.709','2026-04-21 15:29:14.709',NULL,'127.0.0.1','POST','/menu/deleteBaseMenu',200,1277356417,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36','','{\"ID\":2}','{\"code\":0,\"data\":{},\"msg\":\"删除成功\"}',1),(11,'2026-04-21 15:29:33.445','2026-04-21 15:29:33.445',NULL,'127.0.0.1','POST','/menu/deleteBaseMenu',200,2003741334,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36','','{\"ID\":8}','{\"code\":0,\"data\":{},\"msg\":\"删除成功\"}',1),(12,'2026-04-21 15:30:28.527','2026-04-21 15:30:28.527',NULL,'127.0.0.1','POST','/menu/deleteBaseMenu',200,256784750,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36','','{\"ID\":9}','{\"code\":7,\"data\":{},\"msg\":\"删除失败:此菜单存在子菜单不可删除\"}',1),(13,'2026-04-21 15:30:40.914','2026-04-21 15:30:40.914',NULL,'127.0.0.1','POST','/menu/deleteBaseMenu',200,1650915875,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36','','{\"ID\":31}','{\"code\":0,\"data\":{},\"msg\":\"删除成功\"}',1),(14,'2026-04-21 15:30:47.470','2026-04-21 15:30:47.470',NULL,'127.0.0.1','POST','/menu/deleteBaseMenu',200,927849083,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36','','{\"ID\":32}','{\"code\":0,\"data\":{},\"msg\":\"删除成功\"}',1),(15,'2026-04-21 15:30:51.618','2026-04-21 15:30:51.618',NULL,'127.0.0.1','POST','/menu/deleteBaseMenu',200,925071916,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36','','{\"ID\":33}','{\"code\":0,\"data\":{},\"msg\":\"删除成功\"}',1),(16,'2026-04-21 15:30:56.610','2026-04-21 15:30:56.610',NULL,'127.0.0.1','POST','/menu/deleteBaseMenu',200,1895835625,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36','','{\"ID\":34}','{\"code\":0,\"data\":{},\"msg\":\"删除成功\"}',1),(17,'2026-04-21 15:31:01.263','2026-04-21 15:31:01.263',NULL,'127.0.0.1','POST','/menu/deleteBaseMenu',200,939188834,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36','','{\"ID\":35}','{\"code\":0,\"data\":{},\"msg\":\"删除成功\"}',1),(18,'2026-04-21 15:31:05.573','2026-04-21 15:31:05.573',NULL,'127.0.0.1','POST','/menu/deleteBaseMenu',200,930500750,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36','','{\"ID\":9}','{\"code\":0,\"data\":{},\"msg\":\"删除成功\"}',1),(19,'2026-04-21 15:32:05.473','2026-04-21 15:32:05.473',NULL,'127.0.0.1','PUT','/authority/updateAuthority',200,1103266750,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36','','{\"authorityId\":888,\"AuthorityName\":\"普通用户\",\"parentId\":0,\"defaultRouter\":\"user\"}','{\"code\":0,\"data\":{\"authority\":{\"CreatedAt\":\"0001-01-01T00:00:00Z\",\"UpdatedAt\":\"0001-01-01T00:00:00Z\",\"DeletedAt\":null,\"authorityId\":888,\"authorityName\":\"普通用户\",\"parentId\":0,\"dataAuthorityId\":null,\"children\":null,\"menus\":null,\"defaultRouter\":\"user\"}},\"msg\":\"更新成功\"}',1),(20,'2026-04-21 15:32:06.762','2026-04-21 15:32:06.762',NULL,'127.0.0.1','POST','/menu/addMenuAuthority',200,1060034833,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36','','[超出记录长度]','{\"code\":0,\"data\":{},\"msg\":\"添加成功\"}',1),(21,'2026-04-21 15:32:52.508','2026-04-21 15:32:52.508',NULL,'127.0.0.1','POST','/menu/deleteBaseMenu',200,311914625,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36','','{\"ID\":1}','{\"code\":7,\"data\":{},\"msg\":\"删除失败:此菜单有角色正在作为首页，不可删除\"}',1),(22,'2026-04-21 15:33:30.505','2026-04-21 15:33:30.505',NULL,'127.0.0.1','POST','/menu/deleteBaseMenu',200,924103083,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36','','{\"ID\":7}','{\"code\":0,\"data\":{},\"msg\":\"删除成功\"}',1);
/*!40000 ALTER TABLE `sys_operation_records` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_params`
--

DROP TABLE IF EXISTS `sys_params`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_params` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '参数名称',
  `key` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '参数键',
  `value` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '参数值',
  `desc` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '参数说明',
  PRIMARY KEY (`id`),
  KEY `idx_sys_params_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_params`
--

LOCK TABLES `sys_params` WRITE;
/*!40000 ALTER TABLE `sys_params` DISABLE KEYS */;
/*!40000 ALTER TABLE `sys_params` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_user_authority`
--

DROP TABLE IF EXISTS `sys_user_authority`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_user_authority` (
  `sys_user_id` bigint unsigned NOT NULL,
  `sys_authority_authority_id` bigint unsigned NOT NULL COMMENT '角色ID',
  PRIMARY KEY (`sys_user_id`,`sys_authority_authority_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_user_authority`
--

LOCK TABLES `sys_user_authority` WRITE;
/*!40000 ALTER TABLE `sys_user_authority` DISABLE KEYS */;
INSERT INTO `sys_user_authority` VALUES (1,888),(1,8881),(1,9528),(2,888);
/*!40000 ALTER TABLE `sys_user_authority` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_users`
--

DROP TABLE IF EXISTS `sys_users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `uuid` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '用户UUID',
  `username` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '用户登录名',
  `password` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '用户登录密码',
  `nick_name` varchar(191) COLLATE utf8mb4_general_ci DEFAULT '系统用户' COMMENT '用户昵称',
  `header_img` varchar(191) COLLATE utf8mb4_general_ci DEFAULT 'https://qmplusimg.henrongyi.top/gva_header.jpg' COMMENT '用户头像',
  `authority_id` bigint unsigned DEFAULT '888' COMMENT '用户角色ID',
  `phone` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '用户手机号',
  `email` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '用户邮箱',
  `enable` bigint DEFAULT '1' COMMENT '用户是否被冻结 1正常 2冻结',
  `origin_setting` text COLLATE utf8mb4_general_ci COMMENT '配置',
  PRIMARY KEY (`id`),
  KEY `idx_sys_users_deleted_at` (`deleted_at`),
  KEY `idx_sys_users_uuid` (`uuid`),
  KEY `idx_sys_users_username` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_users`
--

LOCK TABLES `sys_users` WRITE;
/*!40000 ALTER TABLE `sys_users` DISABLE KEYS */;
INSERT INTO `sys_users` VALUES (1,'2026-04-21 14:16:29.937','2026-04-21 14:34:08.644',NULL,'6a7547c9-7ee8-4acb-9f91-9cf8e8fcaa6e','admin','$2a$10$jHeBRAgGMqRInhL/8BMlpejVjB4khA/.nSag71avouVYKfHopPWSa','Mr.奇淼','https://qmplusimg.henrongyi.top/gva_header.jpg',888,'17611111111','333333333@qq.com',1,NULL),(2,'2026-04-21 14:16:29.937','2026-04-21 14:16:30.784',NULL,'7b388649-271e-426d-89d5-e8650627903a','a303176530','$2a$10$t0Z2wg43W0VR.HkYL7HKPOdFgLgTQfDcvxabSIvTqHS5iMqRkHFy6','用户1','https://qmplusimg.henrongyi.top/1572075907logo.png',9528,'17611111111','333333333@qq.com',1,NULL);
/*!40000 ALTER TABLE `sys_users` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_versions`
--

DROP TABLE IF EXISTS `sys_versions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_versions` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `version_name` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '版本名称',
  `version_code` varchar(100) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '版本号',
  `description` varchar(500) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '版本描述',
  `version_data` text COLLATE utf8mb4_general_ci COMMENT '版本数据JSON',
  PRIMARY KEY (`id`),
  KEY `idx_sys_versions_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_versions`
--

LOCK TABLES `sys_versions` WRITE;
/*!40000 ALTER TABLE `sys_versions` DISABLE KEYS */;
/*!40000 ALTER TABLE `sys_versions` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-04-21 16:36:49
