package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	// InitializingMessage は初期化中に表示するメッセージ
	InitializingMessage = "初期化中..."
)

// Model はTUIアプリケーションの状態を管理する構造体
type Model struct {
	width        int           // ターミナル幅
	height       int           // ターミナル高さ
	ready        bool          // 初期化完了フラグ
	err          error         // エラー状態
	mainView     *MainView     // メインビューコンポーネント
	statusBar    *StatusBar    // ステータスバーコンポーネント
	eventHandler *EventHandler // イベントハンドラー
}

// NewModel は新しいModelを作成する
func NewModel() Model {
	mainView := NewMainView()
	statusBar := NewStatusBar()
	eventHandler := NewEventHandler()

	// 初期メッセージを追加
	mainView.AddOutput("ccforge - Claude Code TUIアプリケーション")
	mainView.AddOutput("準備完了")
	mainView.AddOutput("")
	mainView.AddOutput("使い方:")
	mainView.AddOutput("  - テキストを入力してEnterキーで送信")
	mainView.AddOutput("  - ↑/↓キーでスクロール")
	mainView.AddOutput("  - F1キーでヘルプ表示切り替え")
	mainView.AddOutput("  - Ctrl+Cまたはqで終了")

	return Model{
		mainView:     mainView,
		statusBar:    statusBar,
		eventHandler: eventHandler,
	}
}

// Init はBubble Teaの初期化処理
func (m Model) Init() tea.Cmd {
	// 初期化時に実行するコマンドはなし
	return nil
}

// Update はメッセージを受け取って状態を更新する
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// イベントハンドラーがある場合はそちらを使用
	if m.eventHandler != nil {
		updatedModel, cmd := m.eventHandler.Handle(&m, msg)
		if updatedModel != nil {
			return *updatedModel, cmd
		}
	}

	// フォールバック処理（互換性のため）
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// グローバルキーバインドの処理
		switch msg.String() {
		case "ctrl+c", "q":
			// 終了
			return m, tea.Quit
		case "f1":
			// ヘルプ表示の切り替え
			m.statusBar.ToggleHelp()
			return m, nil
		case "ctrl+l":
			// 画面クリア
			m.mainView.Clear()
			m.mainView.AddOutput("画面をクリアしました")
			return m, nil
		default:
			// メインビューにキーイベントを渡す
			_, cmd = m.mainView.Update(msg)
			cmds = append(cmds, cmd)
		}

	case tea.WindowSizeMsg:
		// ウィンドウサイズ変更の処理
		m.width = msg.Width
		m.height = msg.Height
		if !m.ready {
			m.ready = true
		}

		// コンポーネントのサイズを更新
		m.mainView.width = msg.Width
		m.mainView.height = msg.Height - 1 // ステータスバーの分を引く
		m.statusBar.SetWidth(msg.Width)

	case error:
		// エラーメッセージの処理
		m.err = msg
		if m.mainView != nil {
			m.mainView.AddOutput(fmt.Sprintf("エラー: %v", msg))
		}
		return m, nil
	}

	return m, tea.Batch(cmds...)
}

// View は現在の状態を文字列として描画する
func (m Model) View() string {
	if !m.ready {
		return InitializingMessage
	}

	if m.err != nil && m.mainView == nil {
		return fmt.Sprintf("エラー: %v\n\nqキーまたはCtrl+Cで終了", m.err)
	}

	// コンポーネントが初期化されていない場合
	if m.mainView == nil || m.statusBar == nil {
		return "コンポーネントを初期化中..."
	}

	// メインビューとステータスバーを結合
	mainContent := m.mainView.View()
	statusContent := m.statusBar.View()

	// 垂直に結合
	return lipgloss.JoinVertical(
		lipgloss.Top,
		mainContent,
		statusContent,
	)
}
