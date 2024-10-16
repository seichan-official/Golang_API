package main

import (
    "database/sql"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    // DSN (Data Source Name) 設定
    dsn := "user:password@tcp(127.0.0.1:3306)/dbname"

    // データベース接続
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        panic(err)
    }
    defer db.Close()

    // 接続テスト
    if err := db.Ping(); err != nil {
        panic(err)
    }

    fmt.Println("データベースに接続しました")
}
