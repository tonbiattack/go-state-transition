package persistence

import (
	"testing"

	"github.com/go-state-transition/go-state-transition/internal/domain"
)

func TestNewManagementDataFromRecord(t *testing.T) {
	t.Run("DBレコードを中間データに詰め替えられる", func(t *testing.T) {
		record := managementRecord{
			ID:         10,
			BankStatus: 2,
		}

		data := newManagementDataFromRecord(record)

		if data.ID != 10 {
			t.Errorf("ID = %d, want %d", data.ID, 10)
		}
		if data.BankStatusCode != 2 {
			t.Errorf("BankStatusCode = %d, want %d", data.BankStatusCode, 2)
		}
	})
}

func TestManagementData_ToDomain(t *testing.T) {
	t.Run("中間データからドメインエンティティへ変換できる", func(t *testing.T) {
		data := managementData{
			ID:             11,
			BankStatusCode: 3,
		}

		m, err := data.ToDomain()
		if err != nil {
			t.Fatalf("ToDomain() error = %v", err)
		}
		if m.ID != 11 {
			t.Errorf("ID = %d, want %d", m.ID, 11)
		}
		if m.BankStatus != domain.BankStatusCompleted {
			t.Errorf("BankStatus = %v, want %v", m.BankStatus, domain.BankStatusCompleted)
		}
	})

	t.Run("不正なステータス値はエラーになる", func(t *testing.T) {
		data := managementData{
			ID:             12,
			BankStatusCode: 9,
		}

		_, err := data.ToDomain()
		if err == nil {
			t.Fatal("ToDomain() error = nil, want error")
		}
	})
}
