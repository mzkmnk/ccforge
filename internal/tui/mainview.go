package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// MainView はメインビューコンポーネント
type MainView struct {
	width        int      // ビューの幅
	height       int      // ビューの高さ
	outputLines  []string // 出力行のリスト
	input        string   // 現在の入力
	scrollOffset int      // スクロールオフセット
	cursorPos    int      // カーソル位置
}

// NewMainView は新しいMainViewを作成する
func NewMainView() *MainView {
	return &MainView{
		width:        80,
		height:       24,
		outputLines:  []string{},
		input:        "",
		scrollOffset: 0,
		cursorPos:    0,
	}
}

// Update はメッセージを処理して状態を更新する
func (m *MainView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyRunes:
			// 文字入力
			runes := string(msg.Runes)
			m.input = m.input[:m.cursorPos] + runes + m.input[m.cursorPos:]
			m.cursorPos += len(runes)

		case tea.KeyBackspace:
			// バックスペース
			if m.cursorPos > 0 {
				m.input = m.input[:m.cursorPos-1] + m.input[m.cursorPos:]
				m.cursorPos--
			}

		case tea.KeyDelete:
			// デリートキー
			if m.cursorPos < len(m.input) {
				m.input = m.input[:m.cursorPos] + m.input[m.cursorPos+1:]
			}

		case tea.KeyEnter:
			// エンターキー - コマンドを実行
			if m.input != "" {
				m.outputLines = append(m.outputLines, "> "+m.input)
				m.input = ""
				m.cursorPos = 0
				// 自動スクロール
				m.autoScroll()
			}

		case tea.KeyUp:
			// 上矢印 - スクロールアップ
			if m.scrollOffset > 0 {
				m.scrollOffset--
			}

		case tea.KeyDown:
			// 下矢印 - スクロールダウン
			maxScroll := m.getMaxScroll()
			if m.scrollOffset < maxScroll {
				m.scrollOffset++
			}

		case tea.KeyPgUp:
			// Page Up - ページ単位でスクロールアップ
			scrollAmount := m.height - 5
			m.scrollOffset -= scrollAmount
			if m.scrollOffset < 0 {
				m.scrollOffset = 0
			}

		case tea.KeyPgDown:
			// Page Down - ページ単位でスクロールダウン
			scrollAmount := m.height - 5
			m.scrollOffset += scrollAmount
			maxScroll := m.getMaxScroll()
			if m.scrollOffset > maxScroll {
				m.scrollOffset = maxScroll
			}

		case tea.KeyLeft:
			// 左矢印 - カーソルを左に移動
			if m.cursorPos > 0 {
				m.cursorPos--
			}

		case tea.KeyRight:
			// 右矢印 - カーソルを右に移動
			if m.cursorPos < len(m.input) {
				m.cursorPos++
			}

		case tea.KeyHome:
			// Home - カーソルを行頭に移動
			m.cursorPos = 0

		case tea.KeyEnd:
			// End - カーソルを行末に移動
			m.cursorPos = len(m.input)
		}

	case tea.WindowSizeMsg:
		// ウィンドウサイズ変更
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

// View は現在の状態を文字列として描画する
func (m *MainView) View() string {
	// 出力エリアのスタイル
	outputStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height - 3) // 入力エリアとボーダー分を引く

	// 入力エリアのスタイル
	inputStyle := lipgloss.NewStyle().
		Width(m.width).
		BorderStyle(lipgloss.NormalBorder()).
		BorderTop(true).
		BorderBottom(false).
		BorderLeft(false).
		BorderRight(false)

	// 出力内容の構築
	visibleLines := m.getVisibleLines()
	outputContent := strings.Join(visibleLines, "\n")

	// スクロールインジケーター
	if m.needsScrollIndicator() {
		scrollInfo := fmt.Sprintf(" [%d/%d]", m.scrollOffset+1, len(m.outputLines))
		outputContent = outputContent + scrollInfo
	}

	// 入力行の構築
	inputLine := "> " + m.input
	if m.cursorPos == len(m.input) {
		inputLine += "█" // カーソル表示
	} else {
		// カーソル位置に応じた表示
		inputLine = "> " + m.input[:m.cursorPos] + "█" + m.input[m.cursorPos:]
	}

	// 最終的なビューの構築
	output := outputStyle.Render(outputContent)
	input := inputStyle.Render(inputLine)

	return output + "\n" + input
}

// AddOutput は出力に新しい行を追加する
func (m *MainView) AddOutput(line string) {
	m.outputLines = append(m.outputLines, line)
	m.autoScroll()
}

// Clear は出力をクリアする
func (m *MainView) Clear() {
	m.outputLines = []string{}
	m.input = ""
	m.scrollOffset = 0
	m.cursorPos = 0
}

// getVisibleLines は現在表示されるべき行を取得する
func (m *MainView) getVisibleLines() []string {
	if len(m.outputLines) == 0 {
		return []string{}
	}

	// 表示可能な行数
	visibleHeight := m.height - 3

	// スクロール位置に基づいて表示する行を決定
	start := m.scrollOffset
	end := start + visibleHeight

	if end > len(m.outputLines) {
		end = len(m.outputLines)
	}

	if start >= len(m.outputLines) {
		start = len(m.outputLines) - 1
		if start < 0 {
			start = 0
		}
	}

	return m.outputLines[start:end]
}

// getMaxScroll は最大スクロール位置を取得する
func (m *MainView) getMaxScroll() int {
	visibleHeight := m.height - 3
	maxScroll := len(m.outputLines) - visibleHeight
	if maxScroll < 0 {
		maxScroll = 0
	}
	return maxScroll
}

// needsScrollIndicator はスクロールインジケーターが必要かどうかを判定する
func (m *MainView) needsScrollIndicator() bool {
	visibleHeight := m.height - 3
	return len(m.outputLines) > visibleHeight
}

// autoScroll は新しい出力が追加されたときに自動的にスクロールする
func (m *MainView) autoScroll() {
	m.scrollOffset = m.getMaxScroll()
}

// Init はBubble Teaの初期化処理（tea.Modelインターフェースの実装）
func (m *MainView) Init() tea.Cmd {
	return nil
}
