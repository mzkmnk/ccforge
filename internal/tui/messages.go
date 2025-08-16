package tui

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
)

// ProcessMessageType はプロセスメッセージの種類を表す
type ProcessMessageType int

const (
	// ProcessOutput はプロセスからの出力
	ProcessOutput ProcessMessageType = iota
	// ProcessStarted はプロセスが開始された
	ProcessStarted
	// ProcessStopped はプロセスが停止した
	ProcessStopped
	// ProcessError はプロセスでエラーが発生した
	ProcessError
)

// エラー定義
var (
	// ErrProcessFailed はプロセスが失敗したことを示すエラー
	ErrProcessFailed = errors.New("プロセスが失敗しました")
)

// KeyboardMessage はキーボード入力を表すメッセージ
type KeyboardMessage struct {
	// Key は押されたキーの文字列表現
	Key string
	// Type はキーのタイプ
	Type tea.KeyType
	// Alt は Alt キーが押されているか
	Alt bool
	// Ctrl は Ctrl キーが押されているか
	Ctrl bool
}

// ToTeaKeyMsg はKeyboardMessageをtea.KeyMsgに変換する
func (k KeyboardMessage) ToTeaKeyMsg() tea.KeyMsg {
	msg := tea.KeyMsg{
		Type: k.Type,
	}

	// Runesの設定
	if k.Type == tea.KeyRunes && len(k.Key) > 0 {
		msg.Runes = []rune(k.Key)
	}

	// 修飾キーの設定
	msg.Alt = k.Alt

	return msg
}

// ProcessMessage はプロセス関連のメッセージ
type ProcessMessage struct {
	// Type はメッセージの種類
	Type ProcessMessageType
	// Output はプロセスの出力内容
	Output string
	// Error はエラー情報
	Error error
	// PID はプロセスID
	PID int
	// Running はプロセスが実行中かどうか
	Running bool
}

// WindowResizeMessage はウィンドウサイズ変更のメッセージ
type WindowResizeMessage struct {
	// Width は新しい幅
	Width int
	// Height は新しい高さ
	Height int
}

// ToTeaWindowSizeMsg はWindowResizeMessageをtea.WindowSizeMsgに変換する
func (w WindowResizeMessage) ToTeaWindowSizeMsg() tea.WindowSizeMsg {
	return tea.WindowSizeMsg{
		Width:  w.Width,
		Height: w.Height,
	}
}

// CommandMessage はコマンド実行要求のメッセージ
type CommandMessage struct {
	// Command は実行するコマンド
	Command string
	// Args はコマンドの引数
	Args []string
}

// ErrorMessage はエラー情報を含むメッセージ
type ErrorMessage struct {
	// Error はエラー内容
	Error error
	// Context はエラーが発生したコンテキスト
	Context string
}
