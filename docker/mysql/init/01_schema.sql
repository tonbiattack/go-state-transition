-- 入金管理テーブル
-- bank_status は BankStatus 型の数値（1:受付 2:処理中 3:完了 4:停止）を格納する
CREATE TABLE IF NOT EXISTS managements (
  id          INT UNSIGNED NOT NULL AUTO_INCREMENT,
  bank_status TINYINT      NOT NULL COMMENT 'ステータス（1:受付 2:処理中 3:完了 4:停止）',
  created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
