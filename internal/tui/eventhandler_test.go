package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// TestEventHandler はイベントハンドラーのテスト
func TestEventHandler(t *testing.T) {
	handler := NewEventHandler()
	model := NewModel()

	tests := []struct {
		name       string
		msg        tea.Msg
		wantUpdate bool
		wantCmd    bool
	}{
		{
			name: "キーボードメッセージ",
			msg: tea.KeyMsg{
				Type:  tea.KeyRunes,
				Runes: []rune("a"),
			},
			wantUpdate: true,
			wantCmd:    false,
		},
		{
			name: "ウィンドウサイズメッセージ",
			msg: tea.WindowSizeMsg{
				Width:  80,
				Height: 24,
			},
			wantUpdate: true,
			wantCmd:    false,
		},
		{
			name: "カスタムキーボードメッセージ",
			msg: KeyboardMessage{
				Key:  "enter",
				Type: tea.KeyEnter,
			},
			wantUpdate: true,
			wantCmd:    false,
		},
		{
			name: "プロセスメッセージ",
			msg: ProcessMessage{
				Type:   ProcessOutput,
				Output: "test output",
				PID:    1234,
			},
			wantUpdate: true,
			wantCmd:    false,
		},
		{
			name: "ウィンドウリサイズメッセージ",
			msg: WindowResizeMessage{
				Width:  100,
				Height: 30,
			},
			wantUpdate: true,
			wantCmd:    false,
		},
		{
			name: "エラーメッセージ",
			msg: ErrorMessage{
				Error:   ErrProcessFailed,
				Context: "test context",
			},
			wantUpdate: true,
			wantCmd:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedModel, cmd := handler.Handle(&model, tt.msg)

			if tt.wantUpdate && updatedModel == nil {
				t.Error("Handle() モデルが更新されるべきですが、nilが返されました")
			}

			if tt.wantCmd && cmd == nil {
				t.Error("Handle() コマンドが返されるべきですが、nilが返されました")
			}
		})
	}
}

// TestRegisterHandler はハンドラー登録のテスト
func TestRegisterHandler(t *testing.T) {
	handler := NewEventHandler()
	model := NewModel()

	// カスタムメッセージ型
	type CustomMessage struct {
		Data string
	}

	// カスタムハンドラーを登録
	handled := false
	customHandler := func(m *Model, msg tea.Msg) (*Model, tea.Cmd) {
		handled = true
		return m, nil
	}

	handler.RegisterHandler("CustomMessage", customHandler)

	// カスタムメッセージを処理
	customMsg := CustomMessage{Data: "test"}
	_, _ = handler.Handle(&model, customMsg)

	if !handled {
		t.Error("カスタムハンドラーが呼び出されませんでした")
	}
}

// TestHandleKeyboardMessage はキーボードメッセージ処理のテスト
func TestHandleKeyboardMessage(t *testing.T) {
	handler := NewEventHandler()
	model := NewModel()

	tests := []struct {
		name       string
		key        tea.KeyMsg
		wantQuit   bool
		wantUpdate bool
	}{
		{
			name: "Ctrl+C（終了）",
			key: tea.KeyMsg{
				Type:  tea.KeyCtrlC,
			},
			wantQuit:   true,
			wantUpdate: false,
		},
		{
			name: "q（終了）",
			key: tea.KeyMsg{
				Type:  tea.KeyRunes,
				Runes: []rune("q"),
			},
			wantQuit:   true,
			wantUpdate: false,
		},
		{
			name: "F1（ヘルプ）",
			key: tea.KeyMsg{
				Type: tea.KeyF1,
			},
			wantQuit:   false,
			wantUpdate: true,
		},
		{
			name: "通常の文字",
			key: tea.KeyMsg{
				Type:  tea.KeyRunes,
				Runes: []rune("a"),
			},
			wantQuit:   false,
			wantUpdate: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedModel, cmd := handler.HandleKeyboardMessage(&model, tt.key)

			if tt.wantQuit {
				if cmd == nil {
					t.Error("終了コマンドが返されるべきです")
				}
			}

			if tt.wantUpdate && updatedModel == nil {
				t.Error("モデルが更新されるべきです")
			}
		})
	}
}

// TestHandleProcessMessage はプロセスメッセージ処理のテスト
func TestHandleProcessMessage(t *testing.T) {
	handler := NewEventHandler()
	model := NewModel()

	tests := []struct {
		name    string
		msg     ProcessMessage
		wantLog bool
	}{
		{
			name: "プロセス出力",
			msg: ProcessMessage{
				Type:   ProcessOutput,
				Output: "test output",
				PID:    1234,
			},
			wantLog: true,
		},
		{
			name: "プロセス開始",
			msg: ProcessMessage{
				Type:    ProcessStarted,
				PID:     5678,
				Running: true,
			},
			wantLog: true,
		},
		{
			name: "プロセス停止",
			msg: ProcessMessage{
				Type:    ProcessStopped,
				PID:     9012,
				Running: false,
			},
			wantLog: true,
		},
		{
			name: "プロセスエラー",
			msg: ProcessMessage{
				Type:  ProcessError,
				Error: ErrProcessFailed,
				PID:   3456,
			},
			wantLog: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedModel, _ := handler.HandleProcessMessage(&model, tt.msg)

			if updatedModel == nil {
				t.Error("モデルが更新されるべきです")
			}
		})
	}
}

// TestHandleWindowResize はウィンドウリサイズ処理のテスト
func TestHandleWindowResize(t *testing.T) {
	handler := NewEventHandler()
	model := NewModel()

	tests := []struct {
		name       string
		msg        WindowResizeMessage
		wantWidth  int
		wantHeight int
	}{
		{
			name: "標準サイズ",
			msg: WindowResizeMessage{
				Width:  80,
				Height: 24,
			},
			wantWidth:  80,
			wantHeight: 24,
		},
		{
			name: "大きなサイズ",
			msg: WindowResizeMessage{
				Width:  200,
				Height: 60,
			},
			wantWidth:  200,
			wantHeight: 60,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedModel, _ := handler.HandleWindowResize(&model, tt.msg)

			if updatedModel == nil {
				t.Error("モデルが更新されるべきです")
			}

			if updatedModel != nil {
				if updatedModel.width != tt.wantWidth {
					t.Errorf("width = %v, want %v", updatedModel.width, tt.wantWidth)
				}
				if updatedModel.height != tt.wantHeight {
					t.Errorf("height = %v, want %v", updatedModel.height, tt.wantHeight)
				}
			}
		})
	}
}