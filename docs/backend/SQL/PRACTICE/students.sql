CREATE DATABASE IF NOT EXISTS SQLpractice CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
use SQLpractice;
show tables;

# 一个名为 students 的表，包含字段 id （主键，自增）、 name （学生姓名，字符串类型）、 age （学生年龄，整数类型）、 grade （学生年级，字符串类型）。
create table students (
    id int primary key auto_increment,
    name varchar(20) not null ,
    age tinyint,
    grade ENUM('一年级', '二年级','三年级','四年级') not null
    );

# 编写SQL语句向 students 表中插入一条新记录，学生姓名为 "张三"，年龄为 20，年级为 "三年级"。

insert into students
set name='张三',age=20,grade='三年级';
select *from students where name='张三';

INSERT INTO students (name, age, grade) VALUES
('李四', 19, '二年级'),
('王五', 21, '四年级'),
('赵六', 17, '二年级'),
('孙七', 22, '四年级'),
('周八', 14, '一年级'),
('吴九', 18, '三年级'),
('郑十', 13, '一年级');

select *from students;

# 编写SQL语句查询 students 表中所有年龄大于 18 岁的学生信息。
select * from students where age>18;

# 编写SQL语句将 students 表中姓名为 "张三" 的学生年级更新为 "四年级"。
update students
    set grade='四年级'
    where name='张三';
select *from students where name='张三';

# 编写SQL语句删除 students 表中年龄小于 15 岁的学生记录。
delete  from students
    where age<15;
select *from students;