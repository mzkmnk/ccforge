package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// TestKeyboardMessage はキーボードメッセージのテスト
func TestKeyboardMessage(t *testing.T) {
	tests := []struct {
		name     string
		msg      KeyboardMessage
		wantKey  string
		wantCtrl bool
		wantAlt  bool
		wantType tea.KeyType
	}{
		{
			name: "通常の文字キー",
			msg: KeyboardMessage{
				Key:  "a",
				Type: tea.KeyRunes,
			},
			wantKey:  "a",
			wantType: tea.KeyRunes,
		},
		{
			name: "Enterキー",
			msg: KeyboardMessage{
				Key:  "enter",
				Type: tea.KeyEnter,
			},
			wantKey:  "enter",
			wantType: tea.KeyEnter,
		},
		{
			name: "Ctrl+Cキー",
			msg: KeyboardMessage{
				Key:  "c",
				Type: tea.KeyRunes,
				Ctrl: true,
			},
			wantKey:  "c",
			wantType: tea.KeyRunes,
			wantCtrl: true,
		},
		{
			name: "F1キー",
			msg: KeyboardMessage{
				Key:  "f1",
				Type: tea.KeyF1,
			},
			wantKey:  "f1",
			wantType: tea.KeyF1,
		},
		{
			name: "矢印キー上",
			msg: KeyboardMessage{
				Key:  "up",
				Type: tea.KeyUp,
			},
			wantKey:  "up",
			wantType: tea.KeyUp,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.msg.Key != tt.wantKey {
				t.Errorf("Key = %v, want %v", tt.msg.Key, tt.wantKey)
			}
			if tt.msg.Type != tt.wantType {
				t.Errorf("Type = %v, want %v", tt.msg.Type, tt.wantType)
			}
			if tt.msg.Ctrl != tt.wantCtrl {
				t.Errorf("Ctrl = %v, want %v", tt.msg.Ctrl, tt.wantCtrl)
			}
			if tt.msg.Alt != tt.wantAlt {
				t.Errorf("Alt = %v, want %v", tt.msg.Alt, tt.wantAlt)
			}
		})
	}
}

// TestProcessMessage はプロセスメッセージのテスト
func TestProcessMessage(t *testing.T) {
	tests := []struct {
		name        string
		msg         ProcessMessage
		wantOutput  string
		wantError   error
		wantPID     int
		wantRunning bool
	}{
		{
			name: "プロセス出力メッセージ",
			msg: ProcessMessage{
				Type:   ProcessOutput,
				Output: "Hello from process",
				PID:    1234,
			},
			wantOutput: "Hello from process",
			wantPID:    1234,
		},
		{
			name: "プロセス開始メッセージ",
			msg: ProcessMessage{
				Type:    ProcessStarted,
				PID:     5678,
				Running: true,
			},
			wantPID:     5678,
			wantRunning: true,
		},
		{
			name: "プロセス停止メッセージ",
			msg: ProcessMessage{
				Type:    ProcessStopped,
				PID:     9012,
				Running: false,
			},
			wantPID:     9012,
			wantRunning: false,
		},
		{
			name: "プロセスエラーメッセージ",
			msg: ProcessMessage{
				Type:  ProcessError,
				Error: ErrProcessFailed,
				PID:   3456,
			},
			wantError: ErrProcessFailed,
			wantPID:   3456,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.msg.Output != tt.wantOutput {
				t.Errorf("Output = %v, want %v", tt.msg.Output, tt.wantOutput)
			}
			if tt.msg.Error != tt.wantError {
				t.Errorf("Error = %v, want %v", tt.msg.Error, tt.wantError)
			}
			if tt.msg.PID != tt.wantPID {
				t.Errorf("PID = %v, want %v", tt.msg.PID, tt.wantPID)
			}
			if tt.msg.Running != tt.wantRunning {
				t.Errorf("Running = %v, want %v", tt.msg.Running, tt.wantRunning)
			}
		})
	}
}

// TestWindowResizeMessage はウィンドウリサイズメッセージのテスト
func TestWindowResizeMessage(t *testing.T) {
	tests := []struct {
		name       string
		msg        WindowResizeMessage
		wantWidth  int
		wantHeight int
	}{
		{
			name: "通常のリサイズ",
			msg: WindowResizeMessage{
				Width:  80,
				Height: 24,
			},
			wantWidth:  80,
			wantHeight: 24,
		},
		{
			name: "大きなウィンドウ",
			msg: WindowResizeMessage{
				Width:  200,
				Height: 60,
			},
			wantWidth:  200,
			wantHeight: 60,
		},
		{
			name: "小さなウィンドウ",
			msg: WindowResizeMessage{
				Width:  40,
				Height: 10,
			},
			wantWidth:  40,
			wantHeight: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.msg.Width != tt.wantWidth {
				t.Errorf("Width = %v, want %v", tt.msg.Width, tt.wantWidth)
			}
			if tt.msg.Height != tt.wantHeight {
				t.Errorf("Height = %v, want %v", tt.msg.Height, tt.wantHeight)
			}
		})
	}
}

// TestConvertToTeaKeyMsg はKeyboardMessageからtea.KeyMsgへの変換をテスト
func TestConvertToTeaKeyMsg(t *testing.T) {
	tests := []struct {
		name string
		msg  KeyboardMessage
		want tea.KeyMsg
	}{
		{
			name: "通常の文字",
			msg: KeyboardMessage{
				Key:  "a",
				Type: tea.KeyRunes,
			},
			want: tea.KeyMsg{
				Type:  tea.KeyRunes,
				Runes: []rune{'a'},
			},
		},
		{
			name: "Ctrl+C",
			msg: KeyboardMessage{
				Key:  "c",
				Type: tea.KeyRunes,
				Ctrl: true,
			},
			want: tea.KeyMsg{
				Type:  tea.KeyRunes,
				Runes: []rune{'c'},
				Alt:   false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.msg.ToTeaKeyMsg()
			if got.Type != tt.want.Type {
				t.Errorf("ToTeaKeyMsg().Type = %v, want %v", got.Type, tt.want.Type)
			}
			// Runesの比較
			if len(got.Runes) > 0 && len(tt.want.Runes) > 0 {
				if got.Runes[0] != tt.want.Runes[0] {
					t.Errorf("ToTeaKeyMsg().Runes[0] = %v, want %v", got.Runes[0], tt.want.Runes[0])
				}
			}
		})
	}
}

// TestConvertToTeaWindowSizeMsg はWindowResizeMessageからtea.WindowSizeMsgへの変換をテスト
func TestConvertToTeaWindowSizeMsg(t *testing.T) {
	tests := []struct {
		name string
		msg  WindowResizeMessage
		want tea.WindowSizeMsg
	}{
		{
			name: "標準サイズ",
			msg: WindowResizeMessage{
				Width:  80,
				Height: 24,
			},
			want: tea.WindowSizeMsg{
				Width:  80,
				Height: 24,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.msg.ToTeaWindowSizeMsg()
			if got.Width != tt.want.Width {
				t.Errorf("ToTeaWindowSizeMsg().Width = %v, want %v", got.Width, tt.want.Width)
			}
			if got.Height != tt.want.Height {
				t.Errorf("ToTeaWindowSizeMsg().Height = %v, want %v", got.Height, tt.want.Height)
			}
		})
	}
}
