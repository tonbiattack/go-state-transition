package testhelper

import (
	"os"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SetupTestDB はテスト用のDB接続を返す
// 環境変数 TEST_DB_DSN が設定されている場合はそれを使用する
// 未設定の場合は docker-compose のデフォルト値を使用する
func SetupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		// docker-compose のデフォルト接続先
		dsn = "appuser:password@tcp(localhost:3307)/go_state_transition?parseTime=true&loc=Local"
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// テスト中は不要なログを抑制する
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("テスト用DBへの接続に失敗しました: %v", err)
	}

	return db
}

// CleanupManagements はテスト後に managements テーブルのデータを削除する
func CleanupManagements(t *testing.T, db *gorm.DB) {
	t.Helper()
	if err := db.Exec("DELETE FROM managements").Error; err != nil {
		t.Fatalf("テストデータの削除に失敗しました: %v", err)
	}
	// AUTO_INCREMENT をリセットする
	if err := db.Exec("ALTER TABLE managements AUTO_INCREMENT = 1").Error; err != nil {
		// リセット失敗はテストを止めるほどではないため警告にとどめる
		t.Logf("AUTO_INCREMENTのリセットに失敗しました（無視します）: %v", err)
	}
}

// MustInsertManagement はテスト用データを挿入して ID を返す
func MustInsertManagement(t *testing.T, db *gorm.DB, bankStatus int) int {
	t.Helper()
	result := db.Exec(
		"INSERT INTO managements (bank_status) VALUES (?)",
		bankStatus,
	)
	if result.Error != nil {
		t.Fatalf("テストデータの挿入に失敗しました: %v", result.Error)
	}
	var id int
	if err := db.Raw("SELECT LAST_INSERT_ID()").Scan(&id).Error; err != nil {
		t.Fatalf("挿入IDの取得に失敗しました: %v", err)
	}
	if id == 0 {
		t.Fatal("挿入IDが0です")
	}
	return id
}
