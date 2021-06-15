create table comment (
     id bigint unsigned not null auto_increment,
     item_name varchar(255) not null comment '商品名',
     comment text comment '评论全文',
     score int not null default 0 comment '评分',
     task_id bigint unsigned not null comment '对应任务id',
     primary key (id) using btree
) comment='评论' engine=InnoDB default charset=utf8;

create table task (
      id bigint unsigned not null unique auto_increment,
      item_name varchar(255) not null comment '商品名',
      url varchar(255) not null comment '商品url',
      status int not null comment '任务状态',
      report text comment '分析结果报告',
      word_cloud_url varchar(255) comment '词云图',
      hot_words text comment '高频词列表'
) comment='分析任务' engine=InnoDB default charset=utf8