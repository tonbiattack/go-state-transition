# go-state-transition

GoでDBのステータスを独自型として表現し、状態遷移の制御を一箇所に集めるサンプル実装です。

## 概要

DBの `TINYINT` ステータスを `int` のまま扱うと、遷移チェックが各所に散らばりやすくなります。
このリポジトリでは、独自型・遷移マップ・ユースケース層での遷移チェックを組み合わせて、状態遷移の制御を整理する方法を示します。

## 構成

```
go-state-transition/
├── internal/domain/
│   ├── bank_status.go       # BankStatus型・メソッド・遷移マップ・変換関数
│   ├── bank_status_test.go
│   ├── management.go        # Managementエンティティ
│   └── management_test.go
└── usecase/
    ├── management_usecase.go      # 遷移チェックを伴うユースケース
    └── management_usecase_test.go
```

## ステータスと遷移ルール

```
受付(1) → 処理中(2) → 完了(3)
                    → 停止(4)
```

完了・停止は終端ステータスであり、そこからの遷移は許可されていません。

## 実行

```bash
# テストを実行する
go test ./...

# 詳細表示
go test -v ./...
```

## 動作環境

- Go 1.21 以上
- 外部ライブラリへの依存なし
