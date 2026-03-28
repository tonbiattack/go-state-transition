package persistence

import (
	"testing"

	"github.com/go-state-transition/go-state-transition/internal/domain"
	"github.com/go-state-transition/go-state-transition/testhelper"
)

func TestManagementRepository_FindByID(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	defer testhelper.CleanupManagements(t, db)

	repo := NewManagementRepository(db)

	t.Run("存在するIDは正しいドメインオブジェクトに変換されて返る", func(t *testing.T) {
		tests := []struct {
			name           string
			bankStatusCode int
			wantStatus     domain.BankStatus
		}{
			{"受付（1）を取得できる", 1, domain.BankStatusAccepted},
			{"処理中（2）を取得できる", 2, domain.BankStatusProcessing},
			{"完了（3）を取得できる", 3, domain.BankStatusCompleted},
			{"停止（4）を取得できる", 4, domain.BankStatusSuspended},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				testhelper.CleanupManagements(t, db)
				id := testhelper.MustInsertManagement(t, db, tt.bankStatusCode)

				m, err := repo.FindByID(id)
				if err != nil {
					t.Fatalf("FindByID() error = %v", err)
				}
				if m.ID != id {
					t.Errorf("ID = %d, want %d", m.ID, id)
				}
				if m.BankStatus != tt.wantStatus {
					t.Errorf("BankStatus = %v, want %v", m.BankStatus, tt.wantStatus)
				}
			})
		}
	})

	t.Run("存在しないIDはエラーを返す", func(t *testing.T) {
		testhelper.CleanupManagements(t, db)
		_, err := repo.FindByID(99999)
		if err == nil {
			t.Error("エラーが返ることを期待しましたが nil でした")
		}
	})

	t.Run("不正なステータス値が入っていた場合はエラーを返す", func(t *testing.T) {
		tests := []struct {
			name           string
			bankStatusCode int
		}{
			{"ステータス0はエラー", 0},
			{"ステータス9はエラー", 9},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				testhelper.CleanupManagements(t, db)
				id := testhelper.MustInsertManagement(t, db, tt.bankStatusCode)
				_, err := repo.FindByID(id)
				if err == nil {
					t.Errorf("不正なステータス %d に対してエラーが返ることを期待しましたが nil でした", tt.bankStatusCode)
				}
			})
		}
	})
}

func TestManagementRepository_Save(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	defer testhelper.CleanupManagements(t, db)

	repo := NewManagementRepository(db)

	tests := []struct {
		name          string
		initialStatus int
		nextStatus    domain.BankStatus
	}{
		{"受付から処理中に更新できる", 1, domain.BankStatusProcessing},
		{"処理中から完了に更新できる", 2, domain.BankStatusCompleted},
		{"処理中から停止に更新できる", 2, domain.BankStatusSuspended},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testhelper.CleanupManagements(t, db)
			id := testhelper.MustInsertManagement(t, db, tt.initialStatus)

			m := &domain.Management{ID: id, BankStatus: tt.nextStatus}
			if err := repo.Save(m); err != nil {
				t.Fatalf("Save() error = %v", err)
			}

			// 保存後に再取得してステータスが更新されていることを確認する
			saved, err := repo.FindByID(id)
			if err != nil {
				t.Fatalf("FindByID() after Save() error = %v", err)
			}
			if saved.BankStatus != tt.nextStatus {
				t.Errorf("保存後のステータス = %v, want %v", saved.BankStatus, tt.nextStatus)
			}
		})
	}
}
