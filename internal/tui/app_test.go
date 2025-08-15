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
		name         string
		wantOutput   string
		wantWidth    int
		wantHeight   int
		wantReady    bool
		wantErrNil   bool
	}{
		{
			name:       "正常系_初期化時の状態",
			wantOutput: "ccforge - Claude Code TUIアプリケーション\n準備完了",
			wantWidth:  0,
			wantHeight: 0,
			wantReady:  false,
			wantErrNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewModel()

			if m.output != tt.wantOutput {
				t.Errorf("NewModel().output = %v, want %v", m.output, tt.wantOutput)
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

// TestModel_Update tests Updateメソッド（各メッセージタイプ）
func TestModel_Update(t *testing.T) {
	tests := []struct {
		name      string
		model     Model
		msg       tea.Msg
		wantCmd   tea.Cmd
		wantWidth int
		wantHeight int
		wantReady bool
		wantErr   error
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
			name:       "正常系_ウィンドウサイズ変更",
			model:      NewModel(),
			msg:        tea.WindowSizeMsg{Width: 100, Height: 50},
			wantWidth:  100,
			wantHeight: 50,
			wantReady:  true,
		},
		{
			name: "正常系_ウィンドウサイズ変更_既にready",
			model: Model{
				ready: true,
			},
			msg:        tea.WindowSizeMsg{Width: 80, Height: 24},
			wantWidth:  80,
			wantHeight: 24,
			wantReady:  true,
		},
		{
			name:    "正常系_エラーメッセージ処理",
			model:   NewModel(),
			msg:     errors.New("テストエラー"),
			wantErr: errors.New("テストエラー"),
		},
		{
			name:    "正常系_その他のキー入力",
			model:   NewModel(),
			msg:     tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}},
			wantCmd: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedModel, cmd := tt.model.Update(tt.msg)
			m := updatedModel.(Model)

			// 終了コマンドのチェック
			if tt.checkQuit {
				if cmd == nil {
					t.Error("Model.Update() expected quit command, got nil")
				}
				// tea.Quitは関数なので、nilでないことだけ確認
				return
			}

			// ウィンドウサイズのチェック
			if _, ok := tt.msg.(tea.WindowSizeMsg); ok {
				if m.width != tt.wantWidth {
					t.Errorf("Model.Update() width = %v, want %v", m.width, tt.wantWidth)
				}
				if m.height != tt.wantHeight {
					t.Errorf("Model.Update() height = %v, want %v", m.height, tt.wantHeight)
				}
				if m.ready != tt.wantReady {
					t.Errorf("Model.Update() ready = %v, want %v", m.ready, tt.wantReady)
				}
			}

			// エラーのチェック
			if _, ok := tt.msg.(error); ok {
				if m.err == nil || m.err.Error() != tt.wantErr.Error() {
					t.Errorf("Model.Update() err = %v, want %v", m.err, tt.wantErr)
				}
			}

			// その他のコマンドチェック
			if !tt.checkQuit && tt.wantCmd == nil && cmd != nil {
				t.Errorf("Model.Update() cmd = %v, want nil", cmd)
			}
		})
	}
}

// TestModel_View tests Viewメソッドのレンダリング
func TestModel_View(t *testing.T) {
	tests := []struct {
		name     string
		model    Model
		want     string
		contains []string
	}{
		{
			name: "正常系_未初期化状態",
			model: Model{
				ready: false,
			},
			want: "初期化中...",
		},
		{
			name: "正常系_エラー状態",
			model: Model{
				ready: true,
				err:   errors.New("テストエラー"),
			},
			contains: []string{"エラー: テストエラー", "qキーまたはCtrl+Cで終了"},
		},
		{
			name: "正常系_通常表示",
			model: Model{
				ready:  true,
				output: "テスト出力",
			},
			contains: []string{
				"テスト出力",
				"キーバインド:",
				"q, Ctrl+C: 終了",
			},
		},
		{
			name: "正常系_デフォルト表示",
			model: func() Model {
				m := NewModel()
				m.ready = true
				return m
			}(),
			contains: []string{
				"ccforge - Claude Code TUIアプリケーション",
				"準備完了",
				"キーバインド:",
				"q, Ctrl+C: 終了",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.model.View()

			// 完全一致のチェック
			if tt.want != "" && got != tt.want {
				t.Errorf("Model.View() = %v, want %v", got, tt.want)
			}

			// 部分文字列のチェック
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
		name    string
		keyMsg  tea.KeyMsg
		want    string
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
		if view != "初期化中..." {
			t.Errorf("Initial view = %v, want '初期化中...'", view)
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
		
		view = m.View()
		if !strings.Contains(view, "統合テストエラー") {
			t.Errorf("View should contain error message, got %v", view)
		}
		
		// 7. 終了
		_, cmd = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
		if cmd == nil {
			t.Error("Update() should return quit command for 'q' key")
		}
	})
}