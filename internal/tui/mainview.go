package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// MainView はメインビューコンポーネント
type MainView struct {
	width          int      // ビューの幅
	height         int      // ビューの高さ
	outputLines    []string // 出力行のリスト
	input          string   // 現在の入力
	scrollOffset   int      // スクロールオフセット
	cursorPos      int      // カーソル位置
	maxOutputLines int      // 最大出力行数 (0 = 無制限)
}

// デフォルトの最大出力行数
const defaultMaxOutputLines = 1000

// NewMainView は新しいMainViewを作成する
func NewMainView() *MainView {
	return &MainView{
		width:          80,
		height:         24,
		outputLines:    []string{},
		input:          "",
		scrollOffset:   0,
		cursorPos:      0,
		maxOutputLines: defaultMaxOutputLines,
	}
}

// Update はメッセージを処理して状態を更新する
func (m *MainView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.handleKeyMsg(msg)
	case tea.WindowSizeMsg:
		// ウィンドウサイズ変更
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

// handleKeyMsg はキーボード入力を処理する
func (m *MainView) handleKeyMsg(msg tea.KeyMsg) {
	switch msg.Type {
	case tea.KeyRunes:
		m.handleTextInput(string(msg.Runes))
	case tea.KeyBackspace:
		m.handleBackspace()
	case tea.KeyDelete:
		m.handleDelete()
	case tea.KeyEnter:
		m.handleEnter()
	case tea.KeyUp:
		m.scrollUp()
	case tea.KeyDown:
		m.scrollDown()
	case tea.KeyPgUp:
		m.pageUp()
	case tea.KeyPgDown:
		m.pageDown()
	case tea.KeyLeft:
		m.moveCursorLeft()
	case tea.KeyRight:
		m.moveCursorRight()
	case tea.KeyHome:
		m.cursorPos = 0
	case tea.KeyEnd:
		// rune数で位置を設定
		m.cursorPos = len([]rune(m.input))
	}
}

// handleTextInput は文字入力を処理する
func (m *MainView) handleTextInput(text string) {
	// runeスライスに変換して処理
	inputRunes := []rune(m.input)
	textRunes := []rune(text)

	// カーソル位置に文字を挿入
	newInput := make([]rune, 0, len(inputRunes)+len(textRunes))
	newInput = append(newInput, inputRunes[:m.cursorPos]...)
	newInput = append(newInput, textRunes...)
	newInput = append(newInput, inputRunes[m.cursorPos:]...)

	m.input = string(newInput)
	m.cursorPos += len(textRunes) // 文字数分カーソルを進める
}

// handleBackspace はバックスペースキーを処理する
func (m *MainView) handleBackspace() {
	if m.cursorPos > 0 {
		// runeスライスに変換して処理
		inputRunes := []rune(m.input)

		// カーソル位置の前の文字を削除
		newInput := make([]rune, 0, len(inputRunes)-1)
		newInput = append(newInput, inputRunes[:m.cursorPos-1]...)
		newInput = append(newInput, inputRunes[m.cursorPos:]...)

		m.input = string(newInput)
		m.cursorPos--
	}
}

// handleDelete はデリートキーを処理する
func (m *MainView) handleDelete() {
	// runeスライスに変換して処理
	inputRunes := []rune(m.input)

	if m.cursorPos < len(inputRunes) {
		// カーソル位置の文字を削除
		newInput := make([]rune, 0, len(inputRunes)-1)
		newInput = append(newInput, inputRunes[:m.cursorPos]...)
		newInput = append(newInput, inputRunes[m.cursorPos+1:]...)

		m.input = string(newInput)
	}
}

// handleEnter はエンターキーを処理する
func (m *MainView) handleEnter() {
	if m.input != "" {
		m.outputLines = append(m.outputLines, "> "+m.input)
		m.input = ""
		m.cursorPos = 0
		m.autoScroll()
	}
}

// scrollUp は上方向へのスクロールを処理する
func (m *MainView) scrollUp() {
	// 負の値の場合は0に修正
	if m.scrollOffset < 0 {
		m.scrollOffset = 0
		return
	}
	if m.scrollOffset > 0 {
		m.scrollOffset--
	}
}

// scrollDown は下方向へのスクロールを処理する
func (m *MainView) scrollDown() {
	maxScroll := m.getMaxScroll()
	// 過大な値の場合は最大値に修正
	if m.scrollOffset > maxScroll {
		m.scrollOffset = maxScroll
		return
	}
	// 負の値の場合は0に修正
	if m.scrollOffset < 0 {
		m.scrollOffset = 0
		return
	}
	if m.scrollOffset < maxScroll {
		m.scrollOffset++
	}
}

// pageUp はページアップを処理する
func (m *MainView) pageUp() {
	scrollAmount := m.height - 5
	m.scrollOffset -= scrollAmount
	if m.scrollOffset < 0 {
		m.scrollOffset = 0
	}
}

// pageDown はページダウンを処理する
func (m *MainView) pageDown() {
	scrollAmount := m.height - 5
	m.scrollOffset += scrollAmount
	maxScroll := m.getMaxScroll()
	if m.scrollOffset > maxScroll {
		m.scrollOffset = maxScroll
	}
}

// moveCursorLeft はカーソルを左に移動する
func (m *MainView) moveCursorLeft() {
	if m.cursorPos > 0 {
		m.cursorPos--
	}
}

// moveCursorRight はカーソルを右に移動する
func (m *MainView) moveCursorRight() {
	// rune数でチェック
	inputRunes := []rune(m.input)
	if m.cursorPos < len(inputRunes) {
		m.cursorPos++
	}
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
		outputContent += scrollInfo
	}

	// 入力行の構築 (runeベースで処理)
	inputRunes := []rune(m.input)
	var inputLine string

	if m.cursorPos >= len(inputRunes) {
		// カーソルが最後にある場合
		inputLine = "> " + m.input + "█"
	} else {
		// カーソルが途中にある場合
		beforeCursor := string(inputRunes[:m.cursorPos])
		afterCursor := string(inputRunes[m.cursorPos:])
		inputLine = "> " + beforeCursor + "█" + afterCursor
	}

	// 最終的なビューの構築
	output := outputStyle.Render(outputContent)
	input := inputStyle.Render(inputLine)

	return output + "\n" + input
}

// AddOutput は出力に新しい行を追加する
func (m *MainView) AddOutput(line string) {
	m.outputLines = append(m.outputLines, line)

	// 最大行数を超えた場合、古い行を削除
	if m.maxOutputLines > 0 && len(m.outputLines) > m.maxOutputLines {
		// 削除する行数を計算
		excessLines := len(m.outputLines) - m.maxOutputLines
		// 古い行を削除して新しいスライスを作成
		m.outputLines = m.outputLines[excessLines:]

		// スクロール位置を調整
		if m.scrollOffset >= excessLines {
			m.scrollOffset -= excessLines
		} else {
			m.scrollOffset = 0
		}
	}

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
	if visibleHeight <= 0 {
		visibleHeight = 0
	}

	// スクロール位置の正規化
	start := m.scrollOffset
	if start < 0 {
		start = 0
	}
	if start >= len(m.outputLines) {
		start = len(m.outputLines) - 1
		if start < 0 {
			start = 0
		}
	}

	// 終了位置の計算
	end := start + visibleHeight
	if end > len(m.outputLines) {
		end = len(m.outputLines)
	}

	// startがendより大きい場合の処理
	if start >= end {
		// 最後の行のみを返す
		if len(m.outputLines) > 0 && start < len(m.outputLines) {
			return m.outputLines[start : start+1]
		}
		return []string{}
	}

	return m.outputLines[start:end]
}

// getMaxScroll は最大スクロール位置を取得する
func (m *MainView) getMaxScroll() int {
	visibleHeight := m.height - 3
	// 高さが極小の場合の処理
	if visibleHeight <= 0 {
		// 表示可能行が0以下の場合、全行数が最大スクロール
		return len(m.outputLines)
	}
	// 表示可能な行数よりも出力行が多い場合のみスクロール可能
	if len(m.outputLines) <= visibleHeight {
		return 0
	}
	// 最後の行まで表示できる最大スクロール位置
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

// SetMaxOutputLines は最大出力行数を設定する
// 0 を設定すると無制限になる
func (m *MainView) SetMaxOutputLines(max int) {
	m.maxOutputLines = max

	// 既存の行数が新しい最大値を超えている場合は削除
	if max > 0 && len(m.outputLines) > max {
		excessLines := len(m.outputLines) - max
		m.outputLines = m.outputLines[excessLines:]

		// スクロール位置を調整
		if m.scrollOffset >= excessLines {
			m.scrollOffset -= excessLines
		} else {
			m.scrollOffset = 0
		}
	}
}

// GetMaxOutputLines は現在の最大出力行数を取得する
func (m *MainView) GetMaxOutputLines() int {
	return m.maxOutputLines
}

// Init はBubble Teaの初期化処理（tea.Modelインターフェースの実装）
func (m *MainView) Init() tea.Cmd {
	return nil
}
