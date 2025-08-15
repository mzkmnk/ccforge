package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	// InitializingMessage は初期化中に表示するメッセージ
	InitializingMessage = "初期化中..."
)

// Model はTUIアプリケーションの状態を管理する構造体
type Model struct {
	width  int    // ターミナル幅
	height int    // ターミナル高さ
	ready  bool   // 初期化完了フラグ
	err    error  // エラー状態
	output string // 出力内容
}

// NewModel は新しいModelを作成する
func NewModel() Model {
	return Model{
		output: "ccforge - Claude Code TUIアプリケーション\n準備完了",
	}
}

// Init はBubble Teaの初期化処理
func (m Model) Init() tea.Cmd {
	// 初期化時に実行するコマンドはなし
	return nil
}

// Update はメッセージを受け取って状態を更新する
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// キー入力の処理
		switch msg.String() {
		case "ctrl+c", "q":
			// 終了
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		// ウィンドウサイズ変更の処理
		m.width = msg.Width
		m.height = msg.Height
		if !m.ready {
			m.ready = true
		}

	case error:
		// エラーメッセージの処理
		m.err = msg
		return m, nil
	}

	return m, nil
}

// View は現在の状態を文字列として描画する
func (m Model) View() string {
	if !m.ready {
		return InitializingMessage
	}

	if m.err != nil {
		return fmt.Sprintf("エラー: %v\n\nqキーまたはCtrl+Cで終了", m.err)
	}

	// 基本的な表示
	view := m.output + "\n\n"
	view += "キーバインド:\n"
	view += "  q, Ctrl+C: 終了\n"

	return view
}
