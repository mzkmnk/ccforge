# プロジェクトセットアップ

## 初期設定
- [x] Go モジュールの初期化 (`go mod init`)
- [x] 必要な依存関係の追加
  - [x] Bubble Tea (`github.com/charmbracelet/bubbletea`)
  - [x] Bubbles (`github.com/charmbracelet/bubbles`)
  - [x] Lipgloss (`github.com/charmbracelet/lipgloss`)
  - [x] PTY (`github.com/creack/pty`)
  - [x] go-git (`github.com/go-git/go-git/v5`)
  - [x] SQLite3 (`github.com/mattn/go-sqlite3`)
  - [x] fsnotify (`github.com/fsnotify/fsnotify`)
  - [x] Viper (`github.com/spf13/viper`)

## プロジェクト構造の作成
- [x] 基本ディレクトリ構造の作成
  - [x] `internal/` ディレクトリ
  - [x] `pkg/` ディレクトリ
  - [x] `cmd/` ディレクトリ
- [x] 各パッケージディレクトリの作成
  - [x] `internal/tui/`
  - [x] `internal/claude/`
  - [x] `internal/tasks/`
  - [x] `internal/config/`

## 開発環境の整備
- [x] GitHub Actions の設定
  - [x] CI/CD パイプライン
  - [x] テスト自動実行
  - [x] リリース自動化
- [x] pre-commit フックの設定
  - [x] gofmt 実行
  - [x] golangci-lint 実行
  - [x] テスト実行