package tui

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatusBar_NewStatusBar(t *testing.T) {
	tests := []struct {
		name      string
		wantWidth int
	}{
		{
			name:      "デフォルトのStatusBarを作成",
			wantWidth: 80,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := NewStatusBar()
			assert.NotNil(t, sb)
			assert.Equal(t, tt.wantWidth, sb.width)
			assert.Equal(t, "", sb.activeTask)
			assert.Equal(t, Disconnected, sb.connectionStatus)
			assert.True(t, sb.showHelp)
		})
	}
}

func TestStatusBar_SetActiveTask(t *testing.T) {
	tests := []struct {
		name     string
		taskName string
		wantTask string
	}{
		{
			name:     "タスク名を設定",
			taskName: "新機能開発",
			wantTask: "新機能開発",
		},
		{
			name:     "空のタスク名を設定",
			taskName: "",
			wantTask: "",
		},
		{
			name:     "長いタスク名を設定",
			taskName: "非常に長いタスク名で表示領域を超える可能性があるもの",
			wantTask: "非常に長いタスク名で表示領域を超える可能性があるもの",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := NewStatusBar()
			sb.SetActiveTask(tt.taskName)
			assert.Equal(t, tt.wantTask, sb.activeTask)
		})
	}
}

func TestStatusBar_SetConnectionStatus(t *testing.T) {
	tests := []struct {
		name       string
		status     ConnectionStatus
		wantStatus ConnectionStatus
	}{
		{
			name:       "接続状態を設定",
			status:     Connected,
			wantStatus: Connected,
		},
		{
			name:       "切断状態を設定",
			status:     Disconnected,
			wantStatus: Disconnected,
		},
		{
			name:       "接続中状態を設定",
			status:     Connecting,
			wantStatus: Connecting,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := NewStatusBar()
			sb.SetConnectionStatus(tt.status)
			assert.Equal(t, tt.wantStatus, sb.connectionStatus)
		})
	}
}

func TestStatusBar_ToggleHelp(t *testing.T) {
	sb := NewStatusBar()

	// 初期状態はヘルプ表示ON
	assert.True(t, sb.showHelp)

	// トグルでOFF
	sb.ToggleHelp()
	assert.False(t, sb.showHelp)

	// 再度トグルでON
	sb.ToggleHelp()
	assert.True(t, sb.showHelp)
}

func TestStatusBar_ViewTask(t *testing.T) {
	tests := []struct {
		name         string
		activeTask   string
		wantContains string
	}{
		{
			name:         "アクティブタスクあり",
			activeTask:   "タスク1",
			wantContains: "タスク: タスク1",
		},
		{
			name:         "アクティブタスクなし",
			activeTask:   "",
			wantContains: "タスク: なし",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := &StatusBar{
				activeTask:       tt.activeTask,
				connectionStatus: Connected,
				showHelp:         true,
				width:            80,
			}
			view := sb.View()
			assert.Contains(t, view, tt.wantContains)
		})
	}
}

func TestStatusBar_ViewConnection(t *testing.T) {
	tests := []struct {
		name             string
		connectionStatus ConnectionStatus
		wantContains     string
	}{
		{
			name:             "接続済み",
			connectionStatus: Connected,
			wantContains:     "接続済み",
		},
		{
			name:             "切断",
			connectionStatus: Disconnected,
			wantContains:     "切断",
		},
		{
			name:             "接続中",
			connectionStatus: Connecting,
			wantContains:     "接続中...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := &StatusBar{
				activeTask:       "test",
				connectionStatus: tt.connectionStatus,
				showHelp:         true,
				width:            80,
			}
			view := sb.View()
			assert.Contains(t, view, tt.wantContains)
		})
	}
}

func TestStatusBar_ViewHelp(t *testing.T) {
	sbWithHelp := &StatusBar{
		activeTask:       "test",
		connectionStatus: Connected,
		showHelp:         true,
		width:            80,
	}

	sbWithoutHelp := &StatusBar{
		activeTask:       "test",
		connectionStatus: Connected,
		showHelp:         false,
		width:            80,
	}

	viewWithHelp := sbWithHelp.View()
	viewWithoutHelp := sbWithoutHelp.View()

	assert.Contains(t, viewWithHelp, "F1: ヘルプ")
	assert.NotContains(t, viewWithoutHelp, "F1: ヘルプ")
}

func TestStatusBar_SetWidth(t *testing.T) {
	tests := []struct {
		name      string
		newWidth  int
		wantWidth int
	}{
		{
			name:      "通常の幅を設定",
			newWidth:  100,
			wantWidth: 100,
		},
		{
			name:      "最小幅を設定",
			newWidth:  40,
			wantWidth: 40,
		},
		{
			name:      "最大幅を設定",
			newWidth:  200,
			wantWidth: 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := NewStatusBar()
			sb.SetWidth(tt.newWidth)
			assert.Equal(t, tt.wantWidth, sb.width)
		})
	}
}

func TestStatusBar_GetConnectionStatusText(t *testing.T) {
	tests := []struct {
		name      string
		status    ConnectionStatus
		wantText  string
		wantColor string // 色のテストは実装では行わないが、仕様として記載
	}{
		{
			name:     "接続済み状態",
			status:   Connected,
			wantText: "● 接続済み",
		},
		{
			name:     "切断状態",
			status:   Disconnected,
			wantText: "● 切断",
		},
		{
			name:     "接続中状態",
			status:   Connecting,
			wantText: "● 接続中...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := &StatusBar{
				connectionStatus: tt.status,
				width:            80,
			}

			view := sb.View()
			assert.Contains(t, view, tt.wantText)
		})
	}
}

func TestStatusBar_LongTaskNameTruncation(t *testing.T) {
	sb := &StatusBar{
		activeTask:       "非常に長いタスク名で表示領域を超える可能性があるものをテストするための文字列",
		connectionStatus: Connected,
		showHelp:         true,
		width:            50, // 狭い幅でテスト
	}

	view := sb.View()

	// ビューの長さが適切に制限されているかチェック
	lines := strings.Split(view, "\n")
	for _, line := range lines {
		// ANSI エスケープシーケンスを除いた実際の文字数をカウント
		// ここでは簡易的に長さチェックのみ行う
		assert.LessOrEqual(t, len(line), 200) // スタイル含めた最大長
	}
}
