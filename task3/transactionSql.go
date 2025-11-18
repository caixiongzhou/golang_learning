package main

import (
	"database/sql"
	"fmt"
)

/*
题目2：事务语句
假设有两个表： accounts 表（包含字段 id 主键， balance 账户余额）和 transactions 表
（包含字段 id 主键， from_account_id 转出账户ID， to_account_id 转入账户ID， amount 转账金额）。
要求 ：
编写一个事务，实现从账户 A 向账户 B 转账 100 元的操作。在事务中，需要先检查账户 A 的余额是否足够，如果足够则从账户 A 扣除 100 元，
向账户 B 增加 100 元，并在 transactions 表中记录该笔转账信息。如果余额不足，则回滚事务。
*/

func TransaferMoney(db *sql.DB, fromAccountID int, toAccountID int, amount float64) error {
	//开始事务
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback() // 如果报错就回滚
			panic(r)      // 重新抛出panic，不让程序静默失败，因为这是未预期的错误
		} else if err != nil {
			// 处理普通错误情况
			tx.Rollback()
		}
	}()
	// 1. 检查转出账户余额
	var currentBalance float64
	err = tx.QueryRow(`SELECT balance FROM accounts WHERE id = ? FOR UPDATE`, fromAccountID).Scan(&currentBalance)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("查询账户余额失败：%v", err)
	}
	if currentBalance < amount {
		tx.Rollback()
		return fmt.Errorf("余额不足：当前余额 %.2f，需要 %.2f", currentBalance, amount)
	}
	// 2.扣除转出账户余额
	_, err = tx.Exec(`UPDATE accounts SET balance = balance - ? WHERE id = ? `, amount, fromAccountID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("扣除金额失败：%v", err)
	}
	// 3.增加转入账户余额
	_, err = tx.Exec(`UPDATE accounts SET balance = balance + ? WHERE id = ?`, amount, toAccountID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("增加金额失败：%v", err)
	}
	// 4.记录交易
	_, err = tx.Exec(`INSERT INTO transaction (from_account_id,to_account_id,amount,transaction_time)
					VALUES (?,?,?,NOW())`, fromAccountID, toAccountID, amount)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("记录交易失败：%v", err)
	}
	// 5.提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败：%v", err)
	}
	return nil
}
