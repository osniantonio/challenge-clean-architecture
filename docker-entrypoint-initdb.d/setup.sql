use orders;

create table if not exists `orders` (
    id char(36) primary key,
    price decimal(10, 5) not null,
    tax decimal(10, 5) not null,
    final_price decimal(10, 5) not null
);