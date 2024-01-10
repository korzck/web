create index price_type on items (price, type);

create index title_idx on items (title);



begin;
drop index title_idx;
explain analyze select * from items where title = '%рез%';
rollback;

explain analyze select * from items where price < 10 and type = 'alloy';

update items set img_url = replace(img_url, '127.0.0.1', '192.168.160.14')