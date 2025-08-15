# プロジェクトセットアップ

## 初期設定
- [ ] Go モジュールの初期化 (`go mod init`)
- [ ] 必要な依存関係の追加
  - [ ] Bubble Tea (`github.com/charmbracelet/bubbletea`)
  - [ ] Bubbles (`github.com/charmbracelet/bubbles`)
  - [ ] Lipgloss (`github.com/charmbracelet/lipgloss`)
  - [ ] PTY (`github.com/creack/pty`)
  - [ ] go-git (`github.com/go-git/go-git/v5`)
  - [ ] SQLite3 (`github.com/mattn/go-sqlite3`)
  - [ ] fsnotify (`github.com/fsnotify/fsnotify`)
  - [ ] Viper (`github.com/spf13/viper`)

## プロジェクト構造の作成
- [ ] 基本ディレクトリ構造の作成
  - [ ] `internal/` ディレクトリ
  - [ ] `pkg/` ディレクトリ
  - [ ] `cmd/` ディレクトリ
- [ ] 各パッケージディレクトリの作成
  - [ ] `internal/tui/`
  - [ ] `internal/claude/`
  - [ ] `internal/tasks/`
  - [ ] `internal/config/`

## 開発環境の整備
- [ ] Makefile の作成
  - [ ] ビルドタスク
  - [ ] テストタスク
  - [ ] リントタスク
  - [ ] インストールタスク
- [ ] GitHub Actions の設定
  - [ ] CI/CD パイプライン
  - [ ] テスト自動実行
  - [ ] リリース自動化
- [ ] pre-commit フックの設定
  - [ ] gofmt 実行
  - [ ] golangci-lint 実行
  - [ ] テスト実行