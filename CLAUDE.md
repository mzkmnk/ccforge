# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## プロジェクト概要

**ccforge**は、Claude Code TUIアプリケーションです。設計駆動開発（Design-Driven Development）を実現するために、Claude CodeをラップしたGoアプリケーションです。

## 言語設定

**重要**: このプロジェクトでは、すべてのコード、コメント、ドキュメント、コミットメッセージを**日本語**で記載してください。

## 技術スタック

### コア技術
- **言語**: Go 1.21+
- **TUIフレームワーク**: Bubble Tea (Elmアーキテクチャベース)
- **UIコンポーネント**: Bubbles (Charmエコシステム)
- **スタイリング**: Lipgloss

### 主要ライブラリ
- **PTY管理**: github.com/creack/pty
- **Git連携**: github.com/go-git/go-git
- **データベース**: github.com/mattn/go-sqlite3
- **ファイル監視**: github.com/fsnotify/fsnotify
- **設定管理**: github.com/spf13/viper

## プロジェクト構造

```
ccforge/
├── main.go                 # エントリーポイント
├── internal/              # 内部パッケージ
│   ├── tui/              # UIコンポーネント
│   │   ├── app.go        # メインアプリケーション
│   │   ├── sidebar.go    # サイドバー実装
│   │   └── session.go    # セッションビュー
│   ├── claude/           # Claude Code統合
│   │   ├── pty.go        # PTY管理
│   │   └── process.go    # プロセス管理
│   ├── tasks/            # タスク管理
│   │   ├── manager.go    # タスクマネージャー
│   │   └── specs.go      # Specs管理
│   └── config/           # 設定管理
├── pkg/                  # 公開パッケージ
└── cmd/                  # CLIコマンド
```

## 開発コマンド

### プロジェクトの初期化
```bash
# Go モジュールの初期化
go mod init github.com/yourusername/ccforge

# 依存関係のダウンロード
go mod download

# 必要な依存関係の追加
go get github.com/charmbracelet/bubbletea
go get github.com/charmbracelet/bubbles
go get github.com/charmbracelet/lipgloss
go get github.com/creack/pty
go get github.com/go-git/go-git/v5
go get github.com/mattn/go-sqlite3
go get github.com/fsnotify/fsnotify
go get github.com/spf13/viper
```

### ビルドとテスト
```bash
# アプリケーションのビルド
go build -o ccforge main.go

# Makefileが作成された場合のビルド
make build

# テストの実行
go test ./...

# カバレッジ付きテスト
go test -cover ./...

# 特定のテストの実行
go test -run TestName ./internal/tui

# ベンチマークテスト
go test -bench=. ./...
```

### 開発とデバッグ
```bash
# 開発モードで実行
go run main.go

# デバッグ情報付きビルド
go build -gcflags="all=-N -l" -o ccforge main.go

# 競合状態の検出
go run -race main.go

# リンターの実行（golangci-lintが必要）
golangci-lint run

# フォーマット
go fmt ./...

# インポートの整理
goimports -w .
```

## アーキテクチャの要点

### Bubble Teaアーキテクチャ
このプロジェクトはElmアーキテクチャパターンを採用しています：
- **Model**: アプリケーションの状態を保持
- **Update**: メッセージに基づいて状態を更新
- **View**: 状態に基づいてUIをレンダリング

### 主要コンポーネント

1. **TUIインターフェース** (`internal/tui/`)
   - メインビュー、サイドバー、ステータスバーの管理
   - キーボードイベントとコマンド処理

2. **タスクマネージャー** (`internal/tasks/`)
   - タスクごとのセッション管理
   - specs（要件定義・設計書）の読み込みと管理

3. **PTYプロセス管理** (`internal/claude/`)
   - Claude Code CLIプロセスの起動と制御
   - 入出力のストリーミング処理

4. **設定管理** (`internal/config/`)
   - グローバル設定とプロジェクト設定の読み込み
   - キーバインディングのカスタマイズ

## データ構造

### タスク管理ディレクトリ
プロジェクトルートに`ccforge/`ディレクトリを作成し、タスクごとに以下の構造で管理：
```
<projectRoot>/
└── ccforge/
    └── {task-name}/
        ├── requirements.md  # 要件定義
        ├── design.md       # 設計書
        └── tasks.md        # TODOリスト
```

## 設定ファイル

### 設定ファイルの場所
- グローバル: `~/.config/ccforge/config.toml`
- プロジェクト: `<projectRoot>/ccforge/config.toml`

### 設定ファイル形式
TOML形式を使用し、以下の構造を持つ：
```toml
[general]
theme = "dark"
auto_save = true

[claude]
model = "claude-3-opus-20240229"

[keybindings]
switch_session = "ctrl+shift+s"

[ui]
sidebar_width = 30
```

## 開発原則

### 基本原則の厳守

このプロジェクトでは以下の開発原則を**必ず遵守**してください：

#### 1. YAGNI原則 (You Aren't Gonna Need It)
- 現時点で必要な機能のみを実装する
- 将来の拡張を過度に考慮した設計を避ける
- シンプルで理解しやすいコードを保つ

#### 2. DRY原則 (Don't Repeat Yourself)
- 同じロジックの重複を避ける
- 共通処理は適切に抽出して再利用する
- ただし、過度な抽象化は避ける（YAGNIとのバランス）

#### 3. SOLID原則
- **S**ingle Responsibility: 一つの構造体/関数は一つの責任のみ
- **O**pen/Closed: 拡張に開き、修正に閉じた設計
- **L**iskov Substitution: インターフェースの一貫性を保つ
- **I**nterface Segregation: 小さく特化したインターフェース
- **D**ependency Inversion: 具象ではなく抽象に依存

#### 4. TDD (Test-Driven Development)
**t-wada氏が推奨するTDDサイクルを厳守**：

1. **Red**: 失敗するテストを先に書く
2. **Green**: テストを通す最小限のコードを書く
3. **Refactor**: コードをリファクタリングする

```go
// テストファースト例
// 1. まずテストを書く (_test.go)
func TestTaskManager_CreateTask(t *testing.T) {
    manager := NewTaskManager()
    task, err := manager.CreateTask("新機能開発")
    
    assert.NoError(t, err)
    assert.Equal(t, "新機能開発", task.Name)
    assert.NotEmpty(t, task.ID)
}

// 2. その後実装を書く
func (m *TaskManager) CreateTask(name string) (*Task, error) {
    // 最小限の実装
}

// 3. リファクタリング
```

### テスト作成の指針
- テストは仕様書として機能する
- テストケース名は日本語で記述可能
- テーブル駆動テストを活用
- モックは最小限に留める

```go
func TestTaskManager(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {
            name:    "正常系_タスク作成",
            input:   "新規タスク",
            want:    "新規タスク",
            wantErr: false,
        },
        {
            name:    "異常系_空文字列",
            input:   "",
            want:    "",
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // テスト実装
        })
    }
}
```

## コーディング規約

### Go標準規約の遵守
- Go公式のコーディングスタイルガイドに従う
- エラーハンドリングは明示的に行う
- goroutineの適切な管理とリソースクリーンアップ

### コメントとドキュメント
- すべてのパブリック関数・型には日本語でGoDocコメントを記載
- 複雑なロジックには日本語で説明コメントを追加

### エラーハンドリング
```go
// エラーは必ず処理する
if err != nil {
    return fmt.Errorf("処理名: %w", err)
}
```

## 実装の優先順位

1. **基本構造の実装**
   - main.goとTUIアプリケーションの骨組み
   - Bubble Teaを使用した基本的なUI

2. **Claude Code統合**
   - PTYプロセス管理
   - 入出力のハンドリング

3. **タスク管理機能**
   - タスクの作成と切り替え
   - specsファイルの読み込み

4. **拡張機能**
   - Git連携による差分表示
   - セッションの永続化
   - カスタムコマンド

## 注意事項

- Claude Code CLIがインストールされていることを前提とする
- ターミナルのサイズ変更に対応する必要がある
- 並行処理では適切なチャネルとミューテックスを使用する
- ユーザー入力の検証とサニタイゼーションを行う