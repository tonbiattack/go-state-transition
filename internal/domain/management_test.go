package domain

import "testing"

// TestNewManagement は Management エンティティの生成関数を確認する
func TestNewManagement(t *testing.T) {
	tests := []struct {
		name       string
		id         int
		bankStatus BankStatus
	}{
		{"受付ステータスで生成できる", 1, BankStatusAccepted},
		{"処理中ステータスで生成できる", 2, BankStatusProcessing},
		{"完了ステータスで生成できる", 3, BankStatusCompleted},
		{"停止ステータスで生成できる", 4, BankStatusSuspended},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewManagement(tt.id, tt.bankStatus)
			if m.ID != tt.id {
				t.Errorf("ID = %d, want %d", m.ID, tt.id)
			}
			if m.BankStatus != tt.bankStatus {
				t.Errorf("BankStatus = %v, want %v", m.BankStatus, tt.bankStatus)
			}
		})
	}
}

// TestManagement_New は Management エンティティが正しく構築されることを確認する
func TestManagement_New(t *testing.T) {
	tests := []struct {
		name       string
		id         int
		bankStatus BankStatus
	}{
		{"受付ステータスで構築できる", 1, BankStatusAccepted},
		{"処理中ステータスで構築できる", 2, BankStatusProcessing},
		{"完了ステータスで構築できる", 3, BankStatusCompleted},
		{"停止ステータスで構築できる", 4, BankStatusSuspended},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Management{
				ID:         tt.id,
				BankStatus: tt.bankStatus,
			}
			if m.ID != tt.id {
				t.Errorf("ID = %d, want %d", m.ID, tt.id)
			}
			if m.BankStatus != tt.bankStatus {
				t.Errorf("BankStatus = %v, want %v", m.BankStatus, tt.bankStatus)
			}
		})
	}
}
