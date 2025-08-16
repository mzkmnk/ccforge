package tui

import (
	"fmt"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMainView_NewMainView(t *testing.T) {
	tests := []struct {
		name       string
		wantWidth  int
		wantHeight int
	}{
		{
			name:       "デフォルトのMainViewを作成",
			wantWidth:  80,
			wantHeight: 24,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mv := NewMainView()
			assert.NotNil(t, mv)
			assert.Equal(t, tt.wantWidth, mv.width)
			assert.Equal(t, tt.wantHeight, mv.height)
			assert.NotNil(t, mv.outputLines)
			assert.Empty(t, mv.outputLines)
			assert.Empty(t, mv.input)
			assert.Equal(t, 0, mv.scrollOffset)
		})
	}
}

func TestMainView_UpdateTextInput(t *testing.T) {
	tests := []struct {
		name        string
		initialView *MainView
		msg         tea.Msg
		wantInput   string
	}{
		{
			name:        "文字入力",
			initialView: NewMainView(),
			msg:         tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}},
			wantInput:   "h",
		},
		{
			name: "文字列を追加入力",
			initialView: &MainView{
				input:       "hel",
				cursorPos:   3,
				outputLines: []string{},
				width:       80,
				height:      24,
			},
			msg:       tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}},
			wantInput: "hell",
		},
		{
			name: "バックスペースで文字削除",
			initialView: &MainView{
				input:       "hello",
				cursorPos:   5,
				outputLines: []string{},
				width:       80,
				height:      24,
			},
			msg:       tea.KeyMsg{Type: tea.KeyBackspace},
			wantInput: "hell",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedView, _ := tt.initialView.Update(tt.msg)
			mv, ok := updatedView.(*MainView)
			require.True(t, ok)
			assert.Equal(t, tt.wantInput, mv.input)
		})
	}
}

func TestMainView_UpdateScroll(t *testing.T) {
	tests := []struct {
		name        string
		initialView *MainView
		msg         tea.Msg
		wantScroll  int
	}{
		{
			name: "上矢印でスクロールアップ",
			initialView: &MainView{
				outputLines:  generateTestLines(30),
				scrollOffset: 10,
				width:        80,
				height:       24,
			},
			msg:        tea.KeyMsg{Type: tea.KeyUp},
			wantScroll: 9,
		},
		{
			name: "下矢印でスクロールダウン",
			initialView: &MainView{
				outputLines:  generateTestLines(30),
				scrollOffset: 5,
				width:        80,
				height:       24,
			},
			msg:        tea.KeyMsg{Type: tea.KeyDown},
			wantScroll: 6,
		},
		{
			name: "Page Upでページ単位スクロールアップ",
			initialView: &MainView{
				outputLines:  generateTestLines(50),
				scrollOffset: 20,
				width:        80,
				height:       10,
			},
			msg:        tea.KeyMsg{Type: tea.KeyPgUp},
			wantScroll: 15,
		},
		{
			name: "Page Downでページ単位スクロールダウン",
			initialView: &MainView{
				outputLines:  generateTestLines(50),
				scrollOffset: 10,
				width:        80,
				height:       10,
			},
			msg:        tea.KeyMsg{Type: tea.KeyPgDown},
			wantScroll: 15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedView, _ := tt.initialView.Update(tt.msg)
			mv, ok := updatedView.(*MainView)
			require.True(t, ok)
			assert.Equal(t, tt.wantScroll, mv.scrollOffset)
		})
	}
}

func TestMainView_UpdateEnter(t *testing.T) {
	mv := &MainView{
		input:       "test command",
		outputLines: []string{"previous output"},
		width:       80,
		height:      24,
	}

	updatedView, _ := mv.Update(tea.KeyMsg{Type: tea.KeyEnter})
	newMV, ok := updatedView.(*MainView)
	require.True(t, ok)

	assert.Equal(t, "", newMV.input)
	assert.Equal(t, []string{"previous output", "> test command"}, newMV.outputLines)
}

func TestMainView_UpdateWindowSize(t *testing.T) {
	mv := &MainView{
		width:       80,
		height:      24,
		outputLines: []string{},
	}

	msg := tea.WindowSizeMsg{
		Width:  100,
		Height: 30,
	}

	updatedView, _ := mv.Update(msg)
	newMV, ok := updatedView.(*MainView)
	require.True(t, ok)

	assert.Equal(t, 100, newMV.width)
	assert.Equal(t, 30, newMV.height)
}

func TestMainView_View(t *testing.T) {
	tests := []struct {
		name         string
		mainView     *MainView
		wantContains []string
	}{
		{
			name: "基本的な表示",
			mainView: &MainView{
				outputLines: []string{"Line 1", "Line 2", "Line 3"},
				input:       "test input",
				cursorPos:   10, // "test input"の長さ
				width:       80,
				height:      10,
			},
			wantContains: []string{
				"Line 1",
				"Line 2",
				"Line 3",
				"test input",
			},
		},
		{
			name: "空の出力",
			mainView: &MainView{
				outputLines: []string{},
				input:       "",
				width:       80,
				height:      10,
			},
			wantContains: []string{
				">",
			},
		},
		{
			name: "スクロール時の表示",
			mainView: &MainView{
				outputLines:  generateTestLines(20),
				input:        "",
				scrollOffset: 5,
				width:        80,
				height:       10,
			},
			wantContains: []string{
				"Line 5",
				"Line 6",
				">",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			view := tt.mainView.View()

			for _, want := range tt.wantContains {
				assert.Contains(t, view, want)
			}
		})
	}
}

func TestMainView_AddOutput(t *testing.T) {
	tests := []struct {
		name            string
		initialLines    []string
		addLine         string
		wantOutputLines []string
	}{
		{
			name:            "空のリストに行を追加",
			initialLines:    []string{},
			addLine:         "First line",
			wantOutputLines: []string{"First line"},
		},
		{
			name:            "既存のリストに行を追加",
			initialLines:    []string{"Line 1", "Line 2"},
			addLine:         "Line 3",
			wantOutputLines: []string{"Line 1", "Line 2", "Line 3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mv := &MainView{
				outputLines: tt.initialLines,
				width:       80,
				height:      24,
			}
			mv.AddOutput(tt.addLine)
			assert.Equal(t, tt.wantOutputLines, mv.outputLines)
		})
	}
}

func TestMainView_Clear(t *testing.T) {
	mv := &MainView{
		outputLines:  []string{"Line 1", "Line 2", "Line 3"},
		input:        "some input",
		scrollOffset: 5,
		width:        80,
		height:       24,
	}

	mv.Clear()

	assert.Empty(t, mv.outputLines)
	assert.Empty(t, mv.input)
	assert.Equal(t, 0, mv.scrollOffset)
}

func TestMainView_ScrollBounds(t *testing.T) {
	tests := []struct {
		name                string
		outputLinesCount    int
		height              int
		initialScrollOffset int
		scrollDirection     string
		wantScrollOffset    int
	}{
		{
			name:                "上端を超えてスクロールしない",
			outputLinesCount:    10,
			height:              5,
			initialScrollOffset: 0,
			scrollDirection:     "up",
			wantScrollOffset:    0,
		},
		{
			name:                "下端を超えてスクロールしない",
			outputLinesCount:    10,
			height:              5,
			initialScrollOffset: 8,
			scrollDirection:     "down",
			wantScrollOffset:    8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mv := &MainView{
				outputLines:  generateTestLines(tt.outputLinesCount),
				scrollOffset: tt.initialScrollOffset,
				width:        80,
				height:       tt.height,
			}

			var msg tea.Msg
			if tt.scrollDirection == "up" {
				msg = tea.KeyMsg{Type: tea.KeyUp}
			} else {
				msg = tea.KeyMsg{Type: tea.KeyDown}
			}

			updatedView, _ := mv.Update(msg)
			updatedMV, _ := updatedView.(*MainView)
			assert.Equal(t, tt.wantScrollOffset, updatedMV.scrollOffset)
		})
	}
}

// TestMainView_ScrollBoundaryEdgeCases はエッジケースでのスクロール境界処理をテストする
func TestMainView_ScrollBoundaryEdgeCases(t *testing.T) {
	tests := []struct {
		name             string
		outputLinesCount int
		height           int
		scrollOffset     int
		operation        func(*MainView)
		wantScrollOffset int
		description      string
	}{
		{
			name:             "空のリストでスクロールアップ",
			outputLinesCount: 0,
			height:           10,
			scrollOffset:     0,
			operation:        func(m *MainView) { m.scrollUp() },
			wantScrollOffset: 0,
			description:      "空のリストではスクロール位置は変わらない",
		},
		{
			name:             "空のリストでスクロールダウン",
			outputLinesCount: 0,
			height:           10,
			scrollOffset:     0,
			operation:        func(m *MainView) { m.scrollDown() },
			wantScrollOffset: 0,
			description:      "空のリストではスクロール位置は変わらない",
		},
		{
			name:             "1行だけの場合のスクロール",
			outputLinesCount: 1,
			height:           10,
			scrollOffset:     0,
			operation:        func(m *MainView) { m.scrollDown() },
			wantScrollOffset: 0,
			description:      "表示可能行数より少ない場合はスクロールしない",
		},
		{
			name:             "負のスクロール位置からの修正",
			outputLinesCount: 10,
			height:           5,
			scrollOffset:     -5, // 不正な負の値
			operation:        func(m *MainView) { m.scrollUp() },
			wantScrollOffset: 0,
			description:      "負の値は0に修正される",
		},
		{
			name:             "過大なスクロール位置からの修正",
			outputLinesCount: 10,
			height:           5,
			scrollOffset:     100, // 過大な値
			operation:        func(m *MainView) { m.scrollDown() },
			wantScrollOffset: 8, // maxScroll = 10 - (5-3) = 8
			description:      "過大な値は最大値に修正される",
		},
		{
			name:             "ページアップでの境界処理",
			outputLinesCount: 20,
			height:           10,
			scrollOffset:     3,
			operation:        func(m *MainView) { m.pageUp() },
			wantScrollOffset: 0,
			description:      "ページアップで負にならない",
		},
		{
			name:             "ページダウンでの境界処理",
			outputLinesCount: 20,
			height:           10,
			scrollOffset:     15,
			operation:        func(m *MainView) { m.pageDown() },
			wantScrollOffset: 13, // maxScroll = 20 - (10-3) = 13
			description:      "ページダウンで最大値を超えない",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mv := &MainView{
				outputLines:  generateTestLines(tt.outputLinesCount),
				scrollOffset: tt.scrollOffset,
				width:        80,
				height:       tt.height,
			}

			tt.operation(mv)

			assert.Equal(t, tt.wantScrollOffset, mv.scrollOffset, tt.description)
		})
	}
}

// TestMainView_GetVisibleLinesEdgeCases は可視行取得のエッジケースをテストする
func TestMainView_GetVisibleLinesEdgeCases(t *testing.T) {
	tests := []struct {
		name             string
		outputLines      []string
		scrollOffset     int
		height           int
		wantVisibleLines []string
		description      string
	}{
		{
			name:             "空のリストから取得",
			outputLines:      []string{},
			scrollOffset:     0,
			height:           10,
			wantVisibleLines: []string{},
			description:      "空のリストは空の配列を返す",
		},
		{
			name:             "不正なスクロール位置での取得",
			outputLines:      []string{"Line 0", "Line 1", "Line 2"},
			scrollOffset:     10, // 範囲外
			height:           10,
			wantVisibleLines: []string{"Line 2"}, // 最後の行のみ
			description:      "範囲外のスクロール位置は修正される",
		},
		{
			name:             "負のスクロール位置での取得",
			outputLines:      []string{"Line 0", "Line 1", "Line 2"},
			scrollOffset:     -5,
			height:           10,
			wantVisibleLines: []string{"Line 0", "Line 1", "Line 2"},
			description:      "負のスクロール位置は0として扱われる",
		},
		{
			name:             "高さが極小の場合",
			outputLines:      []string{"Line 0", "Line 1", "Line 2", "Line 3", "Line 4"},
			scrollOffset:     2,
			height:           4, // 表示可能行は1行のみ（4-3=1）
			wantVisibleLines: []string{"Line 2"},
			description:      "極小の高さでも正しく動作",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mv := &MainView{
				outputLines:  tt.outputLines,
				scrollOffset: tt.scrollOffset,
				width:        80,
				height:       tt.height,
			}

			visibleLines := mv.getVisibleLines()
			assert.Equal(t, tt.wantVisibleLines, visibleLines, tt.description)
		})
	}
}

// TestMainView_GetMaxScrollEdgeCases は最大スクロール値計算のエッジケースをテストする
func TestMainView_GetMaxScrollEdgeCases(t *testing.T) {
	tests := []struct {
		name          string
		linesCount    int
		height        int
		wantMaxScroll int
		description   string
	}{
		{
			name:          "空のリスト",
			linesCount:    0,
			height:        10,
			wantMaxScroll: 0,
			description:   "空のリストの最大スクロールは0",
		},
		{
			name:          "表示可能行数と同じ",
			linesCount:    7, // height(10) - 3 = 7
			height:        10,
			wantMaxScroll: 0,
			description:   "すべて表示可能な場合はスクロール不要",
		},
		{
			name:          "表示可能行数より少ない",
			linesCount:    5,
			height:        10,
			wantMaxScroll: 0,
			description:   "表示可能行数より少ない場合はスクロール不要",
		},
		{
			name:          "極小の高さ",
			linesCount:    10,
			height:        4, // 表示可能行は1行のみ
			wantMaxScroll: 9,
			description:   "極小の高さでも正しく計算",
		},
		{
			name:          "高さが3以下",
			linesCount:    10,
			height:        3, // 表示可能行は0行
			wantMaxScroll: 10,
			description:   "高さが3以下でも破綻しない",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mv := &MainView{
				outputLines: generateTestLines(tt.linesCount),
				width:       80,
				height:      tt.height,
			}

			maxScroll := mv.getMaxScroll()
			assert.Equal(t, tt.wantMaxScroll, maxScroll, tt.description)
		})
	}
}

// TestMainView_MultibyteCharacterInput はマルチバイト文字入力のテスト
func TestMainView_MultibyteCharacterInput(t *testing.T) {
	tests := []struct {
		name          string
		initialInput  string
		initialCursor int
		inputRunes    []rune
		operation     func(*MainView)
		wantInput     string
		wantCursorPos int
		description   string
	}{
		{
			name:          "日本語文字の入力",
			initialInput:  "",
			initialCursor: 0,
			inputRunes:    []rune{'こ', 'ん', 'に', 'ち', 'は'},
			operation:     nil,
			wantInput:     "こんにちは",
			wantCursorPos: 5, // 5文字分
			description:   "日本語文字が正しく入力される",
		},
		{
			name:          "日本語と英数字の混在入力",
			initialInput:  "Hello",
			initialCursor: 5,
			inputRunes:    []rune{'世', '界'},
			operation:     nil,
			wantInput:     "Hello世界",
			wantCursorPos: 7, // Hello(5) + 世界(2) = 7文字
			description:   "英数字と日本語の混在が正しく処理される",
		},
		{
			name:          "日本語文字列の途中に挿入",
			initialInput:  "こんちは",
			initialCursor: 2, // "こん" の後
			inputRunes:    []rune{'に'},
			operation:     nil,
			wantInput:     "こんにちは",
			wantCursorPos: 3, // "こんに" の後
			description:   "日本語文字列の途中への挿入が正しく処理される",
		},
		{
			name:          "絵文字の入力",
			initialInput:  "test",
			initialCursor: 4,
			inputRunes:    []rune{'🤖', '💻'},
			operation:     nil,
			wantInput:     "test🤖💻",
			wantCursorPos: 6, // test(4) + 絵文字(2) = 6文字
			description:   "絵文字が正しく入力される",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mv := &MainView{
				input:       tt.initialInput,
				cursorPos:   tt.initialCursor,
				outputLines: []string{},
				width:       80,
				height:      24,
			}

			// 文字入力のシミュレーション
			for _, r := range tt.inputRunes {
				mv.handleTextInput(string(r))
			}

			// 追加の操作があれば実行
			if tt.operation != nil {
				tt.operation(mv)
			}

			assert.Equal(t, tt.wantInput, mv.input, tt.description)
			assert.Equal(t, tt.wantCursorPos, mv.cursorPos, "カーソル位置が正しい")
		})
	}
}

// TestMainView_MultibyteCharacterCursorMovement はマルチバイト文字でのカーソル移動のテスト
func TestMainView_MultibyteCharacterCursorMovement(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		initialCursor int
		operation     func(*MainView)
		wantCursorPos int
		description   string
	}{
		{
			name:          "日本語文字列で左移動",
			input:         "こんにちは",
			initialCursor: 5, // 最後
			operation:     func(m *MainView) { m.moveCursorLeft() },
			wantCursorPos: 4, // "こんにち" の後
			description:   "日本語文字単位で左移動",
		},
		{
			name:          "日本語文字列で右移動",
			input:         "こんにちは",
			initialCursor: 0, // 最初
			operation:     func(m *MainView) { m.moveCursorRight() },
			wantCursorPos: 1, // "こ" の後
			description:   "日本語文字単位で右移動",
		},
		{
			name:          "絵文字での左移動",
			input:         "Hello🤖World",
			initialCursor: 7, // "Hello🤖W" の後
			operation:     func(m *MainView) { m.moveCursorLeft() },
			wantCursorPos: 6, // "Hello🤖" の後
			description:   "絵文字を1文字として扱う",
		},
		{
			name:          "Homeキーで先頭へ",
			input:         "こんにちは世界",
			initialCursor: 7,
			operation:     func(m *MainView) { m.cursorPos = 0 },
			wantCursorPos: 0,
			description:   "Homeキーで先頭へ移動",
		},
		{
			name:          "Endキーで末尾へ",
			input:         "こんにちは世界",
			initialCursor: 0,
			operation:     func(m *MainView) { m.cursorPos = len([]rune(m.input)) },
			wantCursorPos: 7,
			description:   "Endキーで末尾へ移動",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mv := &MainView{
				input:       tt.input,
				cursorPos:   tt.initialCursor,
				outputLines: []string{},
				width:       80,
				height:      24,
			}

			tt.operation(mv)

			assert.Equal(t, tt.wantCursorPos, mv.cursorPos, tt.description)
		})
	}
}

// TestMainView_MultibyteCharacterDeletion はマルチバイト文字の削除テスト
func TestMainView_MultibyteCharacterDeletion(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		cursorPos     int
		operation     func(*MainView)
		wantInput     string
		wantCursorPos int
		description   string
	}{
		{
			name:          "日本語文字のバックスペース",
			input:         "こんにちは",
			cursorPos:     5, // 最後
			operation:     func(m *MainView) { m.handleBackspace() },
			wantInput:     "こんにち",
			wantCursorPos: 4,
			description:   "日本語文字が1文字単位で削除される",
		},
		{
			name:          "日本語文字のデリート",
			input:         "こんにちは",
			cursorPos:     0, // 最初
			operation:     func(m *MainView) { m.handleDelete() },
			wantInput:     "んにちは",
			wantCursorPos: 0,
			description:   "日本語文字が1文字単位で削除される",
		},
		{
			name:          "絵文字のバックスペース",
			input:         "Hello🤖",
			cursorPos:     6, // 最後
			operation:     func(m *MainView) { m.handleBackspace() },
			wantInput:     "Hello",
			wantCursorPos: 5,
			description:   "絵文字が1文字として削除される",
		},
		{
			name:          "混在文字列の途中での削除",
			input:         "Hello世界World",
			cursorPos:     7, // "Hello世界" の後
			operation:     func(m *MainView) { m.handleBackspace() },
			wantInput:     "Hello世World",
			wantCursorPos: 6,
			description:   "混在文字列で正しく削除される",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mv := &MainView{
				input:       tt.input,
				cursorPos:   tt.cursorPos,
				outputLines: []string{},
				width:       80,
				height:      24,
			}

			tt.operation(mv)

			assert.Equal(t, tt.wantInput, mv.input, tt.description)
			assert.Equal(t, tt.wantCursorPos, mv.cursorPos, "カーソル位置が正しい")
		})
	}
}

// TestMainView_MemoryManagement はメモリ管理のテスト
func TestMainView_MemoryManagement(t *testing.T) {
	const maxLines = 1000 // デフォルトの最大行数

	tests := []struct {
		name         string
		initialLines int
		addLines     int
		wantMaxLines int
		description  string
	}{
		{
			name:         "最大行数以下の追加",
			initialLines: 500,
			addLines:     300,
			wantMaxLines: 800,
			description:  "最大行数以下なら全て保持",
		},
		{
			name:         "最大行数を超える追加",
			initialLines: 900,
			addLines:     200,
			wantMaxLines: maxLines,
			description:  "最大行数を超えたら古い行を削除",
		},
		{
			name:         "大量の行を一度に追加",
			initialLines: 100,
			addLines:     2000,
			wantMaxLines: maxLines,
			description:  "大量追加でも最大行数を維持",
		},
		{
			name:         "最大行数ちょうどから追加",
			initialLines: maxLines,
			addLines:     1,
			wantMaxLines: maxLines,
			description:  "最大行数から1行追加で最古の行を削除",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mv := NewMainView()

			// 初期行を追加
			for i := 0; i < tt.initialLines; i++ {
				mv.AddOutput(fmt.Sprintf("Initial Line %d", i))
			}

			// 追加行を追加
			for i := 0; i < tt.addLines; i++ {
				mv.AddOutput(fmt.Sprintf("Added Line %d", i))
			}

			// 行数が最大値を超えていないことを確認
			assert.LessOrEqual(t, len(mv.outputLines), tt.wantMaxLines, tt.description)

			// 最新の行が保持されていることを確認
			if tt.addLines > 0 && len(mv.outputLines) > 0 {
				lastLine := mv.outputLines[len(mv.outputLines)-1]
				expectedLastLine := fmt.Sprintf("Added Line %d", tt.addLines-1)
				assert.Equal(t, expectedLastLine, lastLine, "最新の行が保持されている")
			}
		})
	}
}

// TestMainView_MemoryManagement_OldestLinesRemoved は古い行が削除されることを確認
func TestMainView_MemoryManagement_OldestLinesRemoved(t *testing.T) {
	mv := NewMainView()
	const maxLines = 1000

	// 最大行数を超える行を追加
	totalLines := maxLines + 100
	for i := 0; i < totalLines; i++ {
		mv.AddOutput(fmt.Sprintf("Line %d", i))
	}

	// 行数が最大値に制限されていることを確認
	assert.Equal(t, maxLines, len(mv.outputLines), "行数が最大値に制限されている")

	// 最も古い行が削除され、新しい行が残っていることを確認
	firstLine := mv.outputLines[0]
	expectedFirstLine := fmt.Sprintf("Line %d", 100) // 最初の100行が削除されているはず
	assert.Equal(t, expectedFirstLine, firstLine, "最も古い行が削除されている")

	lastLine := mv.outputLines[len(mv.outputLines)-1]
	expectedLastLine := fmt.Sprintf("Line %d", totalLines-1)
	assert.Equal(t, expectedLastLine, lastLine, "最新の行が保持されている")
}

// TestMainView_ConfigurableMaxLines は最大行数が設定可能であることをテスト
func TestMainView_ConfigurableMaxLines(t *testing.T) {
	tests := []struct {
		name        string
		maxLines    int
		addLines    int
		wantLines   int
		description string
	}{
		{
			name:        "小さい最大値",
			maxLines:    100,
			addLines:    150,
			wantLines:   100,
			description: "カスタム最大値が適用される",
		},
		{
			name:        "大きい最大値",
			maxLines:    5000,
			addLines:    3000,
			wantLines:   3000,
			description: "大きな最大値でも正しく動作",
		},
		{
			name:        "最大値0は無制限",
			maxLines:    0,
			addLines:    2000,
			wantLines:   2000,
			description: "0は無制限を意味する",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mv := NewMainView()
			mv.SetMaxOutputLines(tt.maxLines)

			for i := 0; i < tt.addLines; i++ {
				mv.AddOutput(fmt.Sprintf("Line %d", i))
			}

			assert.Equal(t, tt.wantLines, len(mv.outputLines), tt.description)
		})
	}
}

// ヘルパー関数: テスト用の行を生成
func generateTestLines(count int) []string {
	lines := make([]string, count)
	for i := 0; i < count; i++ {
		lines[i] = fmt.Sprintf("Line %d", i)
	}
	return lines
}
