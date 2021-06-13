create table comment (
     id bigint unsigned not null unique auto_increment,
     item_name varchar(255) not null comment '商品名',
     comment text comment '评论全文',
     hot_word_list text comment '高频词列表',
     status int not null default 0 comment '状态; 0.Unknown 1.爬取中 2.统计中 3.完成',
     primary key(id)
) comment = '评论表';