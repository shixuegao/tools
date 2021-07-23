CREATE TABLE `bitmap_test` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `topic` varchar(100) NOT NULL,
  `offset` bigint(20) unsigned NOT NULL,
  `bitmap` longblob NOT NULL,
  `stamp` int(10) unsigned NOT NULL,
  `create_time` timestamp NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8