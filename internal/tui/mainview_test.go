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
				width:       80,
				height:      10,
			},
			wantContains: []string{
				"Line 1",
				"Line 2",
				"Line 3",
				"> test input",
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
			initialScrollOffset: 6,
			scrollDirection:     "down",
			wantScrollOffset:    6,
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

// ヘルパー関数: テスト用の行を生成
func generateTestLines(count int) []string {
	lines := make([]string, count)
	for i := 0; i < count; i++ {
		lines[i] = fmt.Sprintf("Line %d", i)
	}
	return lines
}
