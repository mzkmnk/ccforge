package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ConnectionStatus は接続状態を表す型
type ConnectionStatus int

const (
	// Disconnected は切断状態
	Disconnected ConnectionStatus = iota
	// Connecting は接続中状態
	Connecting
	// Connected は接続済み状態
	Connected
)

// StatusBar はステータスバーコンポーネント
type StatusBar struct {
	activeTask       string           // アクティブなタスク名
	connectionStatus ConnectionStatus // 接続状態
	showHelp         bool             // ヘルプ表示フラグ
	width            int              // ステータスバーの幅
}

// NewStatusBar は新しいStatusBarを作成する
func NewStatusBar() *StatusBar {
	return &StatusBar{
		activeTask:       "",
		connectionStatus: Disconnected,
		showHelp:         true,
		width:            80,
	}
}

// SetActiveTask はアクティブなタスクを設定する
func (s *StatusBar) SetActiveTask(taskName string) {
	s.activeTask = taskName
}

// SetConnectionStatus は接続状態を設定する
func (s *StatusBar) SetConnectionStatus(status ConnectionStatus) {
	s.connectionStatus = status
}

// ToggleHelp はヘルプ表示を切り替える
func (s *StatusBar) ToggleHelp() {
	s.showHelp = !s.showHelp
}

// SetWidth はステータスバーの幅を設定する
func (s *StatusBar) SetWidth(width int) {
	s.width = width
}

// View は現在の状態を文字列として描画する
func (s *StatusBar) View() string {
	// スタイルの定義
	baseStyle := lipgloss.NewStyle().
		Width(s.width).
		Background(lipgloss.Color("235")).
		Foreground(lipgloss.Color("252"))

	// 左側: タスク情報
	taskText := s.getTaskText()
	taskStyle := baseStyle.Copy().
		Padding(0, 1).
		Align(lipgloss.Left)

	// 中央: 接続状態
	connectionText := s.getConnectionStatusText()
	connectionStyle := baseStyle.Copy().
		Padding(0, 1).
		Align(lipgloss.Center)

	// 右側: ヘルプ
	helpText := s.getHelpText()
	helpStyle := baseStyle.Copy().
		Padding(0, 1).
		Align(lipgloss.Right)

	// レイアウトの構築
	// 各セクションの幅を計算
	taskWidth := s.width / 3
	connectionWidth := s.width / 3
	helpWidth := s.width - taskWidth - connectionWidth

	// スタイルを適用
	taskSection := taskStyle.Width(taskWidth).Render(taskText)
	connectionSection := connectionStyle.Width(connectionWidth).Render(connectionText)
	helpSection := helpStyle.Width(helpWidth).Render(helpText)

	// 横に並べる
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		taskSection,
		connectionSection,
		helpSection,
	)
}

// getTaskText はタスク表示テキストを取得する
func (s *StatusBar) getTaskText() string {
	if s.activeTask == "" {
		return "タスク: なし"
	}

	taskText := fmt.Sprintf("タスク: %s", s.activeTask)

	// 長すぎる場合は省略
	maxLength := s.width/3 - 4
	if len(taskText) > maxLength && maxLength > 3 {
		taskText = taskText[:maxLength-3] + "..."
	}

	return taskText
}

// getConnectionStatusText は接続状態の表示テキストを取得する
func (s *StatusBar) getConnectionStatusText() string {
	var statusIcon string
	var statusText string
	var statusColor lipgloss.Color

	switch s.connectionStatus {
	case Connected:
		statusIcon = "●"
		statusText = "接続済み"
		statusColor = lipgloss.Color("42") // 緑
	case Connecting:
		statusIcon = "●"
		statusText = "接続中..."
		statusColor = lipgloss.Color("226") // 黄色
	case Disconnected:
		statusIcon = "●"
		statusText = "切断"
		statusColor = lipgloss.Color("196") // 赤
	default:
		statusIcon = "●"
		statusText = "不明"
		statusColor = lipgloss.Color("245") // グレー
	}

	// アイコンに色を適用
	iconStyle := lipgloss.NewStyle().Foreground(statusColor)
	coloredIcon := iconStyle.Render(statusIcon)

	return fmt.Sprintf("%s %s", coloredIcon, statusText)
}

// getHelpText はヘルプ表示テキストを取得する
func (s *StatusBar) getHelpText() string {
	if !s.showHelp {
		return ""
	}

	helpItems := []string{
		"F1: ヘルプ",
		"Ctrl+C: 終了",
	}

	return strings.Join(helpItems, " | ")
}

// GetActiveTask はアクティブなタスク名を取得する
func (s *StatusBar) GetActiveTask() string {
	return s.activeTask
}

// GetConnectionStatus は接続状態を取得する
func (s *StatusBar) GetConnectionStatus() ConnectionStatus {
	return s.connectionStatus
}

// IsHelpVisible はヘルプが表示されているかを取得する
func (s *StatusBar) IsHelpVisible() bool {
	return s.showHelp
}
