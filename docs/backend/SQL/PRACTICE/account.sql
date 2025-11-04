use SQLpractice;

# 两个表： accounts 表（包含字段 id 主键， balance 账户余额）和 transactions 表
# （包含字段 id 主键， from_account_id 转出账户ID， to_account_id 转入账户ID， amount 转账金额）。
create table accounts (
    id int auto_increment primary key ,
    balance decimal(10,2) not null DEFAULT 0.00,
    constraint check_balance_positive check (balance>=0)
);
create table transactions(
    id int primary key auto_increment,
    from_account_id int not null,
    to_account_id int not null,
    amount decimal(10,2) not null,
    -- 交易发生时间
    created_at timestamp default current_timestamp,
    -- 确保用户真实存在
    foreign key (from_account_id) references accounts(id),
    foreign key (to_account_id) references accounts(id),
    -- 确保转账金额必须大于0
    constraint check_amount_positive check (amount>0)
);

insert into accounts (id, balance) values (1, 500.00), (2, 200.00);

# 要求 ：
# 编写一个事务，实现从账户 A 向账户 B 转账 100 元的操作。
# 在事务中，需要先检查账户 A 的余额是否足够，如果足够则从账户 A 扣除 100 元，向账户 B 增加 100 元，
# 并在 transactions 表中记录该笔转账信息。如果余额不足，则回滚事务。

create procedure sp_transfer_money(
    in p_from_account_id int,
    in p_to_account_id int,
    in p_amount decimal(10, 2)
)
    begin
        -- 转入账户余额
        declare current_balance decimal(10,2);
        -- SQL错误时回滚
        declare exit handler for sqlexception
            begin
                rollback ;
            end ;
        start transaction ;
        -- 检查余额
        select accounts.balance into current_balance
            from accounts
                where id=p_from_account_id
            for update ;
        -- 余额足够，进行转账操作
        if current_balance>=p_amount then
            -- 扣除转出账户余额
            update accounts
                set balance=balance-p_amount
                where id=p_from_account_id;
            -- 增加转入账户资金
            update accounts
                set balance=balance+p_amount
                where id=p_to_account_id;
            -- 记录转账
            insert into transactions(from_account_id ,to_account_id,amount)
            values (p_from_account_id,p_to_account_id,p_amount);

            commit ;
        else
            -- 金额不足，转账失败
            rollback ;
        end if;

end ;


call sp_transfer_money(1,2,200.00);