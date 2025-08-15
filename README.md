# ccforge - Claude Code TUI with Design-Driven Development

**ccforge** は Claude Code をラップした設計駆動開発（Design-Driven Development）を実現するターミナルユーザーインターフェース（TUI）アプリケーション。specs（仕様書）管理、セッション管理、リアルタイムな差分表示などの機能により、効率的な開発ワークフローを提供する。

## ✨ 主要機能

### 🎯 設計駆動開発
- **タスクベース管理**: タスクごとに要件定義・設計書・TODOを管理
- **Specs管理**: Markdownファイルで仕様を体系的に記述
- **柔軟な粒度**: 画面単位から機能単位まで、開発者が自由に設計粒度を選択可能
- **Claude Code統合**: specsを基にClaude Codeへシームレスに指示を送信

### ⌨️ コマンドシステム
- **独自コマンド**: `Ctrl+Shift` でアプリ固有のコマンドを実行
- **Claude Codeパススルー**: `/` コマンドの透過的な実行
- **カスタマイズ可能**: ユーザー定義コマンドの追加・編集

### 📊 サイドバー機能
- **タスク一覧**: 現在のプロジェクトのタスク階層表示
- **Specs表示**: 選択中タスクの要件定義・設計書表示
- **タスク管理**: 各タスクのTODOリストと進捗管理
- **差分表示**: Git連携によるリアルタイムな変更ファイル表示
- **ファイルツリー**: プロジェクト構造の視覚的な表示

### 🔄 セッション管理
- **タスク別セッション**: タスクごとに独立したClaude Codeセッション
- **マルチセッション**: 複数のタスクセッションを同時管理
- **コンテキスト保持**: タスク単位でのコンテキスト永続化
- **迅速な切り替え**: ホットキーによる高速タスク切り替え

## 🛠 技術スタック

### コア技術
| カテゴリ | 技術 | 説明 |
|---------|------|------|
| **言語** | Go 1.21+ | 高性能な並行処理とシンプルな開発 |
| **TUIフレームワーク** | Bubble Tea | Elmアーキテクチャベースの優れた状態管理 |
| **UIコンポーネント** | Bubbles | Charmエコシステムの豊富なコンポーネント |
| **スタイリング** | Lipgloss | 柔軟で美しいターミナルスタイリング |

### 統合ライブラリ
| 用途 | ライブラリ | 説明 |
|------|-----------|------|
| **PTY管理** | creack/pty | Claude Codeプロセスの擬似端末制御 |
| **Git連携** | go-git/go-git | 差分検出とバージョン管理 |
| **データベース** | mattn/go-sqlite3 | セッション情報の永続化 |
| **ファイル監視** | fsnotify/fsnotify | specsファイルの自動検出 |
| **設定管理** | spf13/viper | 柔軟な設定ファイル管理 |

## 🏗 アーキテクチャ

### システム構成
```
┌─────────────────────────────────────────────────┐
│                  TUI Interface                   │
│  ┌──────────┬─────────────────┬─────────────┐  │
│  │ Sidebar  │   Main View     │   Status    │  │
│  │          │                 │             │  │
│  │ • Tasks  │  Claude Code    │  Active     │  │
│  │ • Specs  │    Output       │   Task      │  │
│  │ • Files  │                 │             │  │
│  └──────────┴─────────────────┴─────────────┘  │
└─────────────────────────────────────────────────┘
                         │
                         ▼
             ┌─────────────────────┐
             │   Task Manager       │
             │  ┌───────────────┐  │
             │  │ Session       │  │
             │  │   Manager     │  │
             │  └───────────────┘  │
             └─────────────────────┘
                         │
                         ▼
             ┌─────────────────────┐
             │   PTY Process       │
             │     Manager         │
             └─────────────────────┘
                         │
                         ▼
              ┌────────────────────┐
              │   Claude Code CLI   │
              └────────────────────┘
```

### データ構造
```
<projectRoot>/
└── ccforge/
    └── {task-name}/           # タスクごとのディレクトリ
        ├── requirements.md    # 要件定義
        ├── design.md          # 設計書
        └── tasks.md           # タスク管理
```

例：
```
my-project/
└── ccforge/
    ├── dashboard-feature/
    │   ├── requirements.md
    │   ├── design.md
    │   └── tasks.md
    └── auth-refactor/
        ├── requirements.md
        ├── design.md
        └── tasks.md
```

## 📦 インストール

### Homebrewを使用（macOS/Linux）
```bash
brew tap yourusername/ccforge
brew install ccforge
```

### Go installを使用
```bash
go install github.com/yourusername/ccforge@latest
```

### ソースからビルド
```bash
git clone https://github.com/yourusername/ccforge.git
cd ccforge
make build
sudo make install
```

## 🚀 使い方

### 基本的な使用方法
```bash
# プロジェクトディレクトリで起動
cd your-project
ccforge

# 新しいタスクを作成
ccforge new dashboard-feature

# 特定のタスクでセッションを開始
ccforge start auth-refactor

# タスク一覧を表示
ccforge list
```

### キーボードショートカット
| ショートカット | 機能 |
|---------------|------|
| `Ctrl+Shift+S` | タスク切り替え |
| `Ctrl+Shift+N` | 新規タスク作成 |
| `Ctrl+Shift+P` | Specs表示/非表示 |
| `Ctrl+Shift+D` | 差分表示 |
| `Ctrl+Shift+T` | タスク管理 |
| `/` | Claude Codeコマンド入力 |
| `Esc` | メニューを閉じる |

## ⚙️ 設定

### 設定ファイル
- **グローバル設定**: `~/.config/ccforge/config.toml`
- **プロジェクト設定**: `<projectRoot>/ccforge/config.toml`

設定例（config.toml）:
```toml
[general]
theme = "dark"
auto_save = true
session_timeout = 3600

[claude]
model = "claude-3-opus-20240229"
max_tokens = 4096

[keybindings]
switch_session = "ctrl+shift+s"
show_specs = "ctrl+shift+p"
show_diff = "ctrl+shift+d"

[ui]
sidebar_width = 30
show_line_numbers = true
```

## 🔧 開発

### 必要要件
- Go 1.21以上
- Claude Code CLI
- Git

### セットアップ
```bash
# リポジトリのクローン
git clone https://github.com/yourusername/ccforge.git
cd ccforge

# 依存関係のインストール
go mod download

# 開発モードで実行
go run main.go

# テストの実行
go test ./...

# ビルド
make build
```

### プロジェクト構造
```
ccforge/
├── main.go                 # エントリーポイント
├── internal/
│   ├── tui/               # UIコンポーネント
│   │   ├── app.go         # メインアプリケーション
│   │   ├── sidebar.go     # サイドバー
│   │   └── session.go     # セッションビュー
│   ├── claude/            # Claude Code統合
│   │   ├── pty.go         # PTY管理
│   │   └── process.go     # プロセス管理
│   ├── tasks/             # タスク管理
│   │   ├── manager.go     # タスクマネージャー
│   │   └── specs.go       # Specs管理
│   └── config/            # 設定管理
├── pkg/                   # 公開パッケージ
└── cmd/                   # CLIコマンド
```