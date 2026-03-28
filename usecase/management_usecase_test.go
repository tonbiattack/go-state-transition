package usecase

import (
	"fmt"
	"testing"

	"github.com/go-state-transition/go-state-transition/internal/domain"
)

// inMemoryManagementRepository は ManagementRepository のインメモリ実装
type inMemoryManagementRepository struct {
	// データをIDをキーとして保持するマップ
	data map[int]*domain.Management
}

// newInMemoryManagementRepository はインメモリリポジトリを生成する
func newInMemoryManagementRepository(initial []*domain.Management) *inMemoryManagementRepository {
	data := make(map[int]*domain.Management)
	for _, m := range initial {
		// ポインタのコピーを避けるため値をコピーして保存する
		copied := *m
		data[m.ID] = &copied
	}
	return &inMemoryManagementRepository{data: data}
}

// FindByID はIDに対応するエンティティを返す
func (r *inMemoryManagementRepository) FindByID(id int) (*domain.Management, error) {
	m, ok := r.data[id]
	if !ok {
		return nil, fmt.Errorf("ID %d の入金管理が見つかりません", id)
	}
	// 呼び出し元が変更しても影響しないよう値をコピーして返す
	copied := *m
	return &copied, nil
}

// Save はエンティティを保存する
func (r *inMemoryManagementRepository) Save(m *domain.Management) error {
	copied := *m
	r.data[m.ID] = &copied
	return nil
}

// TestManagementUsecase_StartProcessing は StartProcessing の正常・異常ケースを確認する
func TestManagementUsecase_StartProcessing(t *testing.T) {
	tests := []struct {
		name           string
		initialStatus  domain.BankStatus
		wantErr        bool
		wantStatus     domain.BankStatus
	}{
		{
			name:          "受付から処理中への遷移は成功する",
			initialStatus: domain.BankStatusAccepted,
			wantErr:       false,
			wantStatus:    domain.BankStatusProcessing,
		},
		{
			name:          "処理中から処理中への遷移はエラーになる",
			initialStatus: domain.BankStatusProcessing,
			wantErr:       true,
			wantStatus:    domain.BankStatusProcessing,
		},
		{
			name:          "完了から処理中への遷移はエラーになる",
			initialStatus: domain.BankStatusCompleted,
			wantErr:       true,
			wantStatus:    domain.BankStatusCompleted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newInMemoryManagementRepository([]*domain.Management{
				{ID: 1, BankStatus: tt.initialStatus},
			})
			uc := NewManagementUsecase(repo)

			err := uc.StartProcessing(1)
			if (err != nil) != tt.wantErr {
				t.Errorf("StartProcessing() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// エラーがなければステータスが更新されていることを確認する
			if !tt.wantErr {
				saved, _ := repo.FindByID(1)
				if saved.BankStatus != tt.wantStatus {
					t.Errorf("保存後のステータス = %v, want %v", saved.BankStatus, tt.wantStatus)
				}
			}
		})
	}
}

// TestManagementUsecase_Complete は Complete の正常・異常ケースを確認する
func TestManagementUsecase_Complete(t *testing.T) {
	tests := []struct {
		name          string
		initialStatus domain.BankStatus
		wantErr       bool
		wantStatus    domain.BankStatus
	}{
		{
			name:          "処理中から完了への遷移は成功する",
			initialStatus: domain.BankStatusProcessing,
			wantErr:       false,
			wantStatus:    domain.BankStatusCompleted,
		},
		{
			name:          "受付から完了への遷移はエラーになる",
			initialStatus: domain.BankStatusAccepted,
			wantErr:       true,
			wantStatus:    domain.BankStatusAccepted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newInMemoryManagementRepository([]*domain.Management{
				{ID: 1, BankStatus: tt.initialStatus},
			})
			uc := NewManagementUsecase(repo)

			err := uc.Complete(1)
			if (err != nil) != tt.wantErr {
				t.Errorf("Complete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				saved, _ := repo.FindByID(1)
				if saved.BankStatus != tt.wantStatus {
					t.Errorf("保存後のステータス = %v, want %v", saved.BankStatus, tt.wantStatus)
				}
			}
		})
	}
}

// TestManagementUsecase_Suspend は Suspend の正常・異常ケースを確認する
func TestManagementUsecase_Suspend(t *testing.T) {
	tests := []struct {
		name          string
		initialStatus domain.BankStatus
		wantErr       bool
		wantStatus    domain.BankStatus
	}{
		{
			name:          "処理中から停止への遷移は成功する",
			initialStatus: domain.BankStatusProcessing,
			wantErr:       false,
			wantStatus:    domain.BankStatusSuspended,
		},
		{
			name:          "受付から停止への遷移はエラーになる",
			initialStatus: domain.BankStatusAccepted,
			wantErr:       true,
			wantStatus:    domain.BankStatusAccepted,
		},
		{
			name:          "完了から停止への遷移はエラーになる",
			initialStatus: domain.BankStatusCompleted,
			wantErr:       true,
			wantStatus:    domain.BankStatusCompleted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newInMemoryManagementRepository([]*domain.Management{
				{ID: 1, BankStatus: tt.initialStatus},
			})
			uc := NewManagementUsecase(repo)

			err := uc.Suspend(1)
			if (err != nil) != tt.wantErr {
				t.Errorf("Suspend() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				saved, _ := repo.FindByID(1)
				if saved.BankStatus != tt.wantStatus {
					t.Errorf("保存後のステータス = %v, want %v", saved.BankStatus, tt.wantStatus)
				}
			}
		})
	}
}
