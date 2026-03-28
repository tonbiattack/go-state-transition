package persistence

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/go-state-transition/go-state-transition/internal/domain"
)

// managementRecord は managements テーブルのレコード構造体
type managementRecord struct {
	ID         uint `gorm:"primaryKey;autoIncrement"`
	BankStatus int  `gorm:"column:bank_status;not null"`
}

// TableName はテーブル名を返す
func (managementRecord) TableName() string {
	return "managements"
}

// ManagementRepository は managements テーブルへの GORM 実装
type ManagementRepository struct {
	db *gorm.DB
}

// NewManagementRepository は ManagementRepository を生成する
func NewManagementRepository(db *gorm.DB) *ManagementRepository {
	return &ManagementRepository{db: db}
}

// FindByID は ID に対応する Management を返す
// DB から取得した bank_status は NewBankStatus で検証してからドメインオブジェクトに変換する
func (r *ManagementRepository) FindByID(id int) (*domain.Management, error) {
	var record managementRecord
	result := r.db.First(&record, id)
	if result.Error != nil {
		return nil, fmt.Errorf("ID %d の入金管理が見つかりません: %w", id, result.Error)
	}

	// DB の数値を BankStatus に変換する（不正値は error を返す）
	status, err := domain.NewBankStatus(record.BankStatus)
	if err != nil {
		return nil, fmt.Errorf("ステータス変換エラー: %w", err)
	}

	return &domain.Management{
		ID:         int(record.ID),
		BankStatus: status,
	}, nil
}

// Save は Management の状態を DB に保存する
func (r *ManagementRepository) Save(m *domain.Management) error {
	result := r.db.Model(&managementRecord{}).
		Where("id = ?", m.ID).
		Update("bank_status", int(m.BankStatus))
	if result.Error != nil {
		return fmt.Errorf("入金管理の保存に失敗しました: %w", result.Error)
	}
	return nil
}
