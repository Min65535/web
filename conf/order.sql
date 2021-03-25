CREATE TABLE IF NOT EXISTS `demo_order` (
	`id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '默认自增id',
	`order_no` VARCHAR(255) NOT NULL UNIQUE KEY COMMENT '订单id',
	`user_name` VARCHAR(255)  COMMENT '用户名' ,
	`amount` DECIMAL(20,5) COMMENT '金额',
	`status` VARCHAR(20) COMMENT '状态',
	`file_url` VARCHAR(255) COMMENT '文件url',
	PRIMARY KEY(`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;


select * from demo_order
 where
 order_no =  '' or user_name like '' or amount = xx or status = xx order by create_name,amount