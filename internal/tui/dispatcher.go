package tui

import (
	"errors"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// CommandHandler はコマンドを処理する関数型
type CommandHandler func(args []string) error

// CommandDispatcher はコマンドのディスパッチを管理する
type CommandDispatcher struct {
	// handlers はコマンド名からハンドラーへのマップ
	handlers map[string]CommandHandler
}

// NewCommandDispatcher は新しいCommandDispatcherを作成する
func NewCommandDispatcher() *CommandDispatcher {
	dispatcher := &CommandDispatcher{
		handlers: make(map[string]CommandHandler),
	}

	// 組み込みコマンドを登録
	dispatcher.registerBuiltinCommands()

	return dispatcher
}

// registerBuiltinCommands は組み込みコマンドを登録する
func (d *CommandDispatcher) registerBuiltinCommands() {
	// clearコマンド
	d.handlers["clear"] = func(args []string) error {
		return nil
	}

	// helpコマンド
	d.handlers["help"] = func(args []string) error {
		return nil
	}

	// exitコマンド
	d.handlers["exit"] = func(args []string) error {
		return nil
	}

	// taskコマンド
	d.handlers["task"] = func(args []string) error {
		return nil
	}
}

// Register はカスタムコマンドハンドラーを登録する
func (d *CommandDispatcher) Register(command string, handler CommandHandler) {
	d.handlers[command] = handler
}

// Dispatch はコマンドを実行する
func (d *CommandDispatcher) Dispatch(command string, args []string) (bool, error) {
	// 空のコマンドはエラー
	if command == "" {
		return false, errors.New("コマンドが指定されていません")
	}

	// ハンドラーを探す
	handler, exists := d.handlers[command]
	if !exists {
		return false, fmt.Errorf("未知のコマンド: %s", command)
	}

	// ハンドラーを実行
	err := handler(args)
	return true, err
}

// ParseCommand は入力文字列をコマンドと引数に分解する
func (d *CommandDispatcher) ParseCommand(input string) (string, []string) {
	// 前後の空白を削除
	input = strings.TrimSpace(input)

	// 空文字列の場合
	if input == "" {
		return "", []string{}
	}

	// スペースで分割
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return "", []string{}
	}

	// 最初の要素がコマンド、残りが引数
	command := parts[0]
	args := []string{}
	if len(parts) > 1 {
		args = parts[1:]
	}

	return command, args
}

// DispatchToModel はModelに対してコマンドを実行する
func (d *CommandDispatcher) DispatchToModel(model *Model, command string, args []string) tea.Cmd {
	switch command {
	case "clear":
		// 画面クリア
		if model.mainView != nil {
			model.mainView.Clear()
			model.mainView.AddOutput("画面をクリアしました")
		}
		return nil

	case "exit":
		// アプリケーション終了
		return tea.Quit

	case "help":
		// ヘルプ表示
		if model.mainView != nil {
			model.mainView.AddOutput("利用可能なコマンド:")
			model.mainView.AddOutput("  clear - 画面をクリア")
			model.mainView.AddOutput("  help  - このヘルプを表示")
			model.mainView.AddOutput("  task  - タスクを切り替え")
			model.mainView.AddOutput("  exit  - アプリケーションを終了")
		}
		return nil

	case "task":
		// タスク切り替え
		if model.mainView != nil {
			if len(args) > 0 {
				model.mainView.AddOutput(fmt.Sprintf("タスク '%s' に切り替えました", args[0]))
			} else {
				model.mainView.AddOutput("現在のタスク: デフォルト")
			}
		}
		return nil

	default:
		// 未知のコマンド
		if model.mainView != nil {
			model.mainView.AddOutput(fmt.Sprintf("未知のコマンド: %s", command))
			model.mainView.AddOutput("'help' でコマンドの一覧を表示")
		}
		return nil
	}
}
