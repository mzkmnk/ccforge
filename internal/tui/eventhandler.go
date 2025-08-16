package tui

import (
	"fmt"
	"reflect"

	tea "github.com/charmbracelet/bubbletea"
)

// MessageHandler はメッセージを処理する関数型
type MessageHandler func(model *Model, msg tea.Msg) (*Model, tea.Cmd)

// EventHandler はイベント処理を管理する
type EventHandler struct {
	// handlers はメッセージ型からハンドラーへのマップ
	handlers map[string]MessageHandler
	// dispatcher はコマンドディスパッチャー
	dispatcher *CommandDispatcher
}

// NewEventHandler は新しいEventHandlerを作成する
func NewEventHandler() *EventHandler {
	handler := &EventHandler{
		handlers:   make(map[string]MessageHandler),
		dispatcher: NewCommandDispatcher(),
	}

	// 組み込みハンドラーを登録
	handler.registerBuiltinHandlers()

	return handler
}

// registerBuiltinHandlers は組み込みハンドラーを登録する
func (e *EventHandler) registerBuiltinHandlers() {
	// tea.KeyMsg のハンドラー
	e.handlers["KeyMsg"] = func(model *Model, msg tea.Msg) (*Model, tea.Cmd) {
		keyMsg := msg.(tea.KeyMsg)
		return e.HandleKeyboardMessage(model, keyMsg)
	}

	// tea.WindowSizeMsg のハンドラー
	e.handlers["WindowSizeMsg"] = func(model *Model, msg tea.Msg) (*Model, tea.Cmd) {
		sizeMsg := msg.(tea.WindowSizeMsg)
		resizeMsg := WindowResizeMessage{
			Width:  sizeMsg.Width,
			Height: sizeMsg.Height,
		}
		return e.HandleWindowResize(model, resizeMsg)
	}

	// KeyboardMessage のハンドラー
	e.handlers["KeyboardMessage"] = func(model *Model, msg tea.Msg) (*Model, tea.Cmd) {
		keyMsg := msg.(KeyboardMessage)
		teaKeyMsg := keyMsg.ToTeaKeyMsg()
		return e.HandleKeyboardMessage(model, teaKeyMsg)
	}

	// ProcessMessage のハンドラー
	e.handlers["ProcessMessage"] = func(model *Model, msg tea.Msg) (*Model, tea.Cmd) {
		procMsg := msg.(ProcessMessage)
		return e.HandleProcessMessage(model, procMsg)
	}

	// WindowResizeMessage のハンドラー
	e.handlers["WindowResizeMessage"] = func(model *Model, msg tea.Msg) (*Model, tea.Cmd) {
		resizeMsg := msg.(WindowResizeMessage)
		return e.HandleWindowResize(model, resizeMsg)
	}

	// ErrorMessage のハンドラー
	e.handlers["ErrorMessage"] = func(model *Model, msg tea.Msg) (*Model, tea.Cmd) {
		errMsg := msg.(ErrorMessage)
		return e.HandleErrorMessage(model, errMsg)
	}

	// CommandMessage のハンドラー
	e.handlers["CommandMessage"] = func(model *Model, msg tea.Msg) (*Model, tea.Cmd) {
		cmdMsg := msg.(CommandMessage)
		return e.HandleCommandMessage(model, cmdMsg)
	}
}

// RegisterHandler はカスタムハンドラーを登録する
func (e *EventHandler) RegisterHandler(msgType string, handler MessageHandler) {
	e.handlers[msgType] = handler
}

// Handle はメッセージを処理する
func (e *EventHandler) Handle(model *Model, msg tea.Msg) (*Model, tea.Cmd) {
	// メッセージの型名を取得
	msgType := reflect.TypeOf(msg).Name()

	// ハンドラーを探す
	if handler, exists := e.handlers[msgType]; exists {
		return handler(model, msg)
	}

	// デフォルト処理
	return model, nil
}

// HandleKeyboardMessage はキーボードメッセージを処理する
func (e *EventHandler) HandleKeyboardMessage(model *Model, msg tea.KeyMsg) (*Model, tea.Cmd) {
	// グローバルキーバインドの処理
	switch msg.String() {
	case "ctrl+c", "q":
		// 終了
		return model, tea.Quit
	
	case "f1":
		// ヘルプ表示の切り替え
		if model.statusBar != nil {
			model.statusBar.ToggleHelp()
		}
		return model, nil
	
	case "ctrl+l":
		// 画面クリア
		if model.mainView != nil {
			model.mainView.Clear()
			model.mainView.AddOutput("画面をクリアしました")
		}
		return model, nil
	
	default:
		// メインビューにキーイベントを渡す
		if model.mainView != nil {
			_, cmd := model.mainView.Update(msg)
			return model, cmd
		}
		return model, nil
	}
}

// HandleProcessMessage はプロセスメッセージを処理する
func (e *EventHandler) HandleProcessMessage(model *Model, msg ProcessMessage) (*Model, tea.Cmd) {
	if model.mainView == nil {
		return model, nil
	}

	switch msg.Type {
	case ProcessOutput:
		// プロセス出力を表示
		model.mainView.AddOutput(msg.Output)
	
	case ProcessStarted:
		// プロセス開始を表示
		model.mainView.AddOutput(fmt.Sprintf("プロセス開始 (PID: %d)", msg.PID))
		if model.statusBar != nil {
			model.statusBar.SetStatus(fmt.Sprintf("プロセス実行中 (PID: %d)", msg.PID))
		}
	
	case ProcessStopped:
		// プロセス停止を表示
		model.mainView.AddOutput(fmt.Sprintf("プロセス停止 (PID: %d)", msg.PID))
		if model.statusBar != nil {
			model.statusBar.SetStatus("待機中")
		}
	
	case ProcessError:
		// プロセスエラーを表示
		model.mainView.AddOutput(fmt.Sprintf("プロセスエラー (PID: %d): %v", msg.PID, msg.Error))
		if model.statusBar != nil {
			model.statusBar.SetStatus("エラー")
		}
	}

	return model, nil
}

// HandleWindowResize はウィンドウリサイズを処理する
func (e *EventHandler) HandleWindowResize(model *Model, msg WindowResizeMessage) (*Model, tea.Cmd) {
	// ウィンドウサイズを更新
	model.width = msg.Width
	model.height = msg.Height

	// 初期化完了フラグを設定
	if !model.ready {
		model.ready = true
	}

	// コンポーネントのサイズを更新
	if model.mainView != nil {
		model.mainView.width = msg.Width
		model.mainView.height = msg.Height - 1 // ステータスバーの分を引く
	}

	if model.statusBar != nil {
		model.statusBar.SetWidth(msg.Width)
	}

	return model, nil
}

// HandleErrorMessage はエラーメッセージを処理する
func (e *EventHandler) HandleErrorMessage(model *Model, msg ErrorMessage) (*Model, tea.Cmd) {
	// エラーを記録
	model.err = msg.Error

	// エラーを表示
	if model.mainView != nil {
		model.mainView.AddOutput(fmt.Sprintf("エラー (%s): %v", msg.Context, msg.Error))
	}

	// ステータスバーにエラーを表示
	if model.statusBar != nil {
		model.statusBar.SetStatus(fmt.Sprintf("エラー: %s", msg.Context))
	}

	return model, nil
}

// HandleCommandMessage はコマンドメッセージを処理する
func (e *EventHandler) HandleCommandMessage(model *Model, msg CommandMessage) (*Model, tea.Cmd) {
	// コマンドをディスパッチ
	cmd := e.dispatcher.DispatchToModel(model, msg.Command, msg.Args)
	return model, cmd
}