package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mzkmnk/ccforge/internal/tui"
)

// Application はアプリケーションの状態を保持する構造体
type Application struct {
	program    *tea.Program
	testMode   bool // テストモード用フラグ
	forceError bool // テスト用エラー強制フラグ
}

// parseCLIArgs はコマンドライン引数を解析する
func parseCLIArgs(args []string) (help bool, err error) {
	fs := flag.NewFlagSet("ccforge", flag.ContinueOnError)
	fs.BoolVar(&help, "h", false, "ヘルプを表示")
	fs.BoolVar(&help, "help", false, "ヘルプを表示")
	
	// 引数をパース
	err = fs.Parse(args)
	if err != nil {
		return false, err
	}
	
	return help, nil
}

// initializeApp はアプリケーションを初期化する
func initializeApp() (*Application, error) {
	// テスト用エラー処理
	if os.Getenv("CCFORGE_INIT_ERROR") == "true" {
		return nil, fmt.Errorf("初期化エラー（テスト用）")
	}
	
	// テストモードの確認
	testMode := os.Getenv("CCFORGE_TEST_MODE") == "true"
	
	app := &Application{
		testMode: testMode,
	}
	
	// テストモードでない場合はBubble Teaプログラムを初期化
	if !testMode {
		// TUIモデルの作成
		model := tui.NewModel()
		
		// Bubble Teaプログラムの作成
		p := tea.NewProgram(model, tea.WithAltScreen())
		app.program = p
	}
	
	return app, nil
}

// runApp はアプリケーションを実行する
func runApp(app *Application) error {
	if app == nil {
		return fmt.Errorf("アプリケーションがnilです")
	}
	
	// テスト用エラー処理
	if app.forceError {
		return fmt.Errorf("実行時エラー（テスト用）")
	}
	
	// テストモードの場合は即座に終了
	if app.testMode {
		return nil
	}
	
	// Bubble Teaプログラムの実行
	if app.program != nil {
		if _, err := app.program.Run(); err != nil {
			return fmt.Errorf("プログラム実行エラー: %w", err)
		}
	}
	
	return nil
}

// mainFlow はmain関数のロジックを分離した関数（テスト用）
func mainFlow(args []string) error {
	// CLIコマンドのパース
	help, err := parseCLIArgs(args)
	if err != nil {
		return fmt.Errorf("引数パースエラー: %w", err)
	}
	
	// ヘルプ表示
	if help {
		showHelp()
		return nil
	}
	
	// アプリケーションの初期化
	app, err := initializeApp()
	if err != nil {
		return fmt.Errorf("初期化エラー: %w", err)
	}
	
	// アプリケーションの実行
	if err := runApp(app); err != nil {
		return fmt.Errorf("実行エラー: %w", err)
	}
	
	return nil
}

// showHelp はヘルプメッセージを表示する
func showHelp() {
	help := `
ccforge - Claude Code TUIアプリケーション

使用法:
  ccforge [オプション]

オプション:
  -h, --help    このヘルプメッセージを表示

説明:
  ccforgeは、Claude CodeをラップしたTUIアプリケーションです。
  設計駆動開発（Design-Driven Development）を実現するための
  タスク管理とセッション管理機能を提供します。

キーバインド:
  Ctrl+C        アプリケーションを終了
  Tab           フォーカスを切り替え
  ↑/↓          項目を選択
  Enter         選択した項目を実行
`
	fmt.Print(help)
}

func main() {
	// os.Argsから実行ファイル名を除いた引数を取得
	args := os.Args[1:]
	
	// メインフローの実行
	if err := mainFlow(args); err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}
}
