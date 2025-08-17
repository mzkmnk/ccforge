package tui

import (
	"errors"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// TestNewModel tests Model構造体の作成
func TestNewModel(t *testing.T) {
	tests := []struct {
		name          string
		wantWidth     int
		wantHeight    int
		wantReady     bool
		wantErrNil    bool
		wantMainView  bool
		wantStatusBar bool
	}{
		{
			name:          "正常系_初期化時の状態",
			wantWidth:     0,
			wantHeight:    0,
			wantReady:     false,
			wantErrNil:    true,
			wantMainView:  true,
			wantStatusBar: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewModel()

			if m.mainView == nil && tt.wantMainView {
				t.Error("NewModel().mainView should not be nil")
			}
			if m.statusBar == nil && tt.wantStatusBar {
				t.Error("NewModel().statusBar should not be nil")
			}
			if m.width != tt.wantWidth {
				t.Errorf("NewModel().width = %v, want %v", m.width, tt.wantWidth)
			}
			if m.height != tt.wantHeight {
				t.Errorf("NewModel().height = %v, want %v", m.height, tt.wantHeight)
			}
			if m.ready != tt.wantReady {
				t.Errorf("NewModel().ready = %v, want %v", m.ready, tt.wantReady)
			}
			if (m.err == nil) != tt.wantErrNil {
				t.Errorf("NewModel().err = %v, want nil = %v", m.err, tt.wantErrNil)
			}
		})
	}
}

// TestModel_Init tests Initメソッド
func TestModel_Init(t *testing.T) {
	tests := []struct {
		name    string
		model   Model
		wantNil bool
	}{
		{
			name:    "正常系_初期化コマンドなし",
			model:   NewModel(),
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.model.Init()

			if (cmd == nil) != tt.wantNil {
				t.Errorf("Model.Init() = %v, want nil = %v", cmd, tt.wantNil)
			}
		})
	}
}

// TestModel_Update_KeyMessages tests キー入力メッセージの処理
func TestModel_Update_KeyMessages(t *testing.T) {
	tests := []struct {
		name      string
		model     Model
		msg       tea.Msg
		checkQuit bool
	}{
		{
			name:      "正常系_Ctrl+C終了",
			model:     NewModel(),
			msg:       tea.KeyMsg{Type: tea.KeyCtrlC},
			checkQuit: true,
		},
		{
			name:      "正常系_qキー終了",
			model:     NewModel(),
			msg:       tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}},
			checkQuit: true,
		},
		{
			name:      "正常系_その他のキー入力",
			model:     NewModel(),
			msg:       tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}},
			checkQuit: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, cmd := tt.model.Update(tt.msg)

			if tt.checkQuit {
				if cmd == nil {
					t.Error("Model.Update() expected quit command, got nil")
				}
			} else {
				if cmd != nil {
					t.Errorf("Model.Update() cmd = %v, want nil", cmd)
				}
			}
		})
	}
}

// TestModel_Update_WindowResize tests ウィンドウリサイズメッセージの処理
func TestModel_Update_WindowResize(t *testing.T) {
	tests := []struct {
		name       string
		model      Model
		msg        tea.WindowSizeMsg
		wantWidth  int
		wantHeight int
		wantReady  bool
	}{
		{
			name:       "正常系_ウィンドウサイズ変更",
			model:      NewModel(),
			msg:        tea.WindowSizeMsg{Width: 100, Height: 50},
			wantWidth:  100,
			wantHeight: 50,
			wantReady:  true,
		},
		{
			name: "正常系_ウィンドウサイズ変更_既にready",
			model: func() Model {
				m := NewModel()
				m.ready = true
				return m
			}(),
			msg:        tea.WindowSizeMsg{Width: 80, Height: 24},
			wantWidth:  80,
			wantHeight: 24,
			wantReady:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedModel, _ := tt.model.Update(tt.msg)
			m := updatedModel.(Model)

			if m.width != tt.wantWidth {
				t.Errorf("Model.Update() width = %v, want %v", m.width, tt.wantWidth)
			}
			if m.height != tt.wantHeight {
				t.Errorf("Model.Update() height = %v, want %v", m.height, tt.wantHeight)
			}
			if m.ready != tt.wantReady {
				t.Errorf("Model.Update() ready = %v, want %v", m.ready, tt.wantReady)
			}
		})
	}
}

// TestModel_Update_ErrorMessage tests エラーメッセージの処理
func TestModel_Update_ErrorMessage(t *testing.T) {
	model := NewModel()
	// ウィンドウサイズを設定してreadyにする
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	model = updatedModel.(Model)

	testErr := errors.New("テストエラー")

	// mainViewのポインタアドレスを確認（デバッグ用）
	originalMainView := model.mainView
	t.Logf("Original mainView pointer: %p", originalMainView)

	updatedModel, cmd := model.Update(testErr)
	m := updatedModel.(Model)

	t.Logf("Updated mainView pointer: %p", m.mainView)
	t.Logf("Are they the same? %v", originalMainView == m.mainView)

	// 値レシーバーのため、m.errの変更は反映されない
	// エラーメッセージがmainViewに追加されることを確認するのが適切

	// mainViewがnilでないことを確認
	if m.mainView == nil {
		t.Fatal("mainView should not be nil")
	}

	if cmd != nil {
		t.Errorf("Model.Update() cmd = %v, want nil", cmd)
	}

	// mainViewに直接エラーメッセージが追加されているか確認
	// （初期メッセージ8行 + エラーメッセージ1行 = 9行）
	if len(m.mainView.outputLines) < 9 {
		t.Errorf("Expected at least 9 lines in outputLines, got %d", len(m.mainView.outputLines))
		for i, line := range m.mainView.outputLines {
			t.Logf("Line %d: %s", i, line)
		}
	}

	// Viewにエラーメッセージが表示されることを確認
	view := m.View()
	if !strings.Contains(view, "エラー: テストエラー") {
		// デバッグ出力
		t.Logf("View content:\n%s", view)
		t.Logf("scrollOffset: %d", m.mainView.scrollOffset)
		t.Errorf("View should contain error message 'エラー: テストエラー'")
	}
}

// TestModel_View_NotReady tests View未初期化状態の表示
func TestModel_View_NotReady(t *testing.T) {
	model := Model{
		ready: false,
	}

	got := model.View()
	want := InitializingMessage

	if got != want {
		t.Errorf("Model.View() = %v, want %v", got, want)
	}
}

// TestModel_View_Error tests Viewエラー状態の表示
func TestModel_View_Error(t *testing.T) {
	model := Model{
		ready: true,
		err:   errors.New("テストエラー"),
	}

	got := model.View()
	expectedContains := []string{
		"エラー: テストエラー",
		"qキーまたはCtrl+Cで終了",
	}

	for _, substr := range expectedContains {
		if !strings.Contains(got, substr) {
			t.Errorf("Model.View() does not contain %q, got %v", substr, got)
		}
	}
}

// TestModel_View_Normal tests View通常状態の表示
func TestModel_View_Normal(t *testing.T) {
	tests := []struct {
		name     string
		model    Model
		contains []string
	}{
		{
			name: "コンポーネント初期化済み",
			model: func() Model {
				m := NewModel()
				m.ready = true
				return m
			}(),
			contains: []string{
				"ccforge",
				"準備完了",
			},
		},
		{
			name: "デフォルト出力",
			model: func() Model {
				m := NewModel()
				m.ready = true
				return m
			}(),
			contains: []string{
				"ccforge - Claude Code TUIアプリケーション",
				"準備完了",
				"使い方:",
				"Ctrl+Cまたはqで終了",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.model.View()

			for _, substr := range tt.contains {
				if !strings.Contains(got, substr) {
					t.Errorf("Model.View() does not contain %q, got %v", substr, got)
				}
			}
		})
	}
}

// TestKeyMsg_String tests キーメッセージの文字列変換
func TestKeyMsg_String(t *testing.T) {
	tests := []struct {
		name   string
		keyMsg tea.KeyMsg
		want   string
	}{
		{
			name:   "正常系_Ctrl+C",
			keyMsg: tea.KeyMsg{Type: tea.KeyCtrlC},
			want:   "ctrl+c",
		},
		{
			name:   "正常系_qキー",
			keyMsg: tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}},
			want:   "q",
		},
		{
			name:   "正常系_Enterキー",
			keyMsg: tea.KeyMsg{Type: tea.KeyEnter},
			want:   "enter",
		},
		{
			name:   "正常系_Tabキー",
			keyMsg: tea.KeyMsg{Type: tea.KeyTab},
			want:   "tab",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.keyMsg.String()
			if got != tt.want {
				t.Errorf("KeyMsg.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestModel_Integration tests 統合テスト
func TestModel_Integration(t *testing.T) {
	t.Run("統合テスト_初期化から終了まで", func(t *testing.T) {
		// 1. モデルの作成
		m := NewModel()

		// 2. 初期化
		cmd := m.Init()
		if cmd != nil {
			t.Error("Init() should return nil")
		}

		// 3. 初期表示の確認
		view := m.View()
		if view != InitializingMessage {
			t.Errorf("Initial view = %v, want %v", view, InitializingMessage)
		}

		// 4. ウィンドウサイズの設定
		updatedModel, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
		m = updatedModel.(Model)

		if !m.ready {
			t.Error("Model should be ready after window size message")
		}

		// 5. 通常表示の確認
		view = m.View()
		if !strings.Contains(view, "ccforge") {
			t.Errorf("View should contain 'ccforge', got %v", view)
		}

		// 6. エラーの発生
		updatedModel, _ = m.Update(errors.New("統合テストエラー"))
		m = updatedModel.(Model)

		// エラーメッセージが mainView に追加されていることを確認
		// mainViewに "エラー: 統合テストエラー" という形式で追加される
		view = m.View()
		if !strings.Contains(view, "エラー: 統合テストエラー") {
			t.Errorf("View should contain 'エラー: 統合テストエラー', got %v", view)
		}

		// 7. 終了
		_, cmd = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
		if cmd == nil {
			t.Error("Update() should return quit command for 'q' key")
		}
	})
}
