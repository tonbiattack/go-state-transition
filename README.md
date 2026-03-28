# go-state-transition

GoでDBのステータスを独自型として表現し、状態遷移の制御を一箇所に集めるサンプル実装です。

## 概要

DBの `TINYINT` ステータスを `int` のまま扱うと、遷移チェックや値の妥当性確認が各所に散らばりやすくなります。  
このリポジトリでは、以下を組み合わせて状態遷移の制御を整理する方法を示します。

- 独自型 `BankStatus` による状態表現
- 遷移マップによる許可ルールの集中管理
- `usecase` 層での状態遷移制御
- `persistence` 層での DB レコードから中間データを経由した変換

特に、DB から取得したレコードをいきなりドメイン層へマッピングせず、永続化層で一度中間データへ詰め替えてから `domain.Management` に変換する構成にしています。

## 構成

```text
go-state-transition/
├── docker/
│   └── mysql/
│       └── init/
│           └── 01_schema.sql      # MySQL初期化SQL
├── infrastructure/
│   └── persistence/
│       ├── management_repository.go      # GORMによる取得・保存
│       ├── management_repository_test.go
│       ├── management_mapper.go          # DBレコード -> 中間データ -> ドメイン変換
│       └── management_mapper_test.go
├── internal/domain/
│   ├── bank_status.go       # BankStatus型・メソッド・遷移マップ・変換関数
│   ├── bank_status_test.go
│   ├── management.go        # Managementエンティティ・生成関数
│   └── management_test.go
├── testhelper/
│   └── db.go                # テスト用DB接続とデータ投入ヘルパー
└── usecase/
    ├── management_usecase.go      # 遷移チェックを伴うユースケース
    └── management_usecase_test.go
```

## ステータスと遷移ルール

```text
受付(1) → 処理中(2) → 完了(3)
                    → 停止(4)
```

完了・停止は終端ステータスであり、そこからの遷移は許可されていません。

## レイヤ構成

### `internal/domain`

ドメインルールを表現する層です。

- `BankStatus` がステータス値の意味と妥当性を持つ
- `Management` が業務エンティティを表す
- `NewBankStatus` で DB の数値を業務上有効な値へ変換する

### `usecase`

ユースケースを表現する層です。

- `StartProcessing`
- `Complete`
- `Suspend`

各ユースケースは、取得した `Management` に対して遷移可否を確認し、許可された場合のみ保存します。

### `infrastructure/persistence`

MySQL + GORM による永続化層です。

- `managementRecord`: DB テーブルに対応するレコード構造体
- `managementData`: 永続化層内で使う中間データ
- `ManagementRepository`: DB 取得・保存の実装

変換の流れは次のとおりです。

```text
managements テーブル
  -> managementRecord
  -> managementData
  -> domain.Management
```

この構成により、DB スキーマの都合をそのままドメイン層へ持ち込まないようにしています。

## DBスキーマ

MySQL の初期スキーマは [`docker/mysql/init/01_schema.sql`](./docker/mysql/init/01_schema.sql) にあります。

`managements` テーブルは `bank_status` を `TINYINT` で保持します。

- `1`: 受付
- `2`: 処理中
- `3`: 完了
- `4`: 停止

## 実行

```bash
# 全テストを実行する
go test ./...

# 詳細表示
go test -v ./...
```

## MySQL を使った統合テスト

このプロジェクトには MySQL を使った永続化テストが含まれています。

### 1. MySQL を起動する

```bash
docker compose up -d
```

### 2. テストを実行する

```bash
go test ./...
```

### 3. 停止する

```bash
docker compose down
```

`testhelper/db.go` は、環境変数 `TEST_DB_DSN` が未設定の場合、以下の接続先を使用します。

```text
appuser:password@tcp(localhost:3307)/go_state_transition?parseTime=true&loc=Local
```

## 動作環境

- Go 1.21 以上
- Docker / Docker Compose
- MySQL 8.0
- GORM
