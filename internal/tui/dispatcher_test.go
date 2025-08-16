package tui

import (
	"testing"
)

// TestCommandDispatcher はコマンドディスパッチャーのテスト
func TestCommandDispatcher(t *testing.T) {
	tests := []struct {
		name        string
		command     string
		args        []string
		wantHandled bool
		wantError   bool
	}{
		{
			name:        "clearコマンド",
			command:     "clear",
			args:        []string{},
			wantHandled: true,
			wantError:   false,
		},
		{
			name:        "helpコマンド",
			command:     "help",
			args:        []string{},
			wantHandled: true,
			wantError:   false,
		},
		{
			name:        "exitコマンド",
			command:     "exit",
			args:        []string{},
			wantHandled: true,
			wantError:   false,
		},
		{
			name:        "taskコマンド_引数あり",
			command:     "task",
			args:        []string{"new-task"},
			wantHandled: true,
			wantError:   false,
		},
		{
			name:        "taskコマンド_引数なし",
			command:     "task",
			args:        []string{},
			wantHandled: true,
			wantError:   false,
		},
		{
			name:        "未知のコマンド",
			command:     "unknown",
			args:        []string{},
			wantHandled: false,
			wantError:   true,
		},
		{
			name:        "空のコマンド",
			command:     "",
			args:        []string{},
			wantHandled: false,
			wantError:   true,
		},
	}

	dispatcher := NewCommandDispatcher()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handled, err := dispatcher.Dispatch(tt.command, tt.args)

			if handled != tt.wantHandled {
				t.Errorf("Dispatch() handled = %v, want %v", handled, tt.wantHandled)
			}

			if (err != nil) != tt.wantError {
				t.Errorf("Dispatch() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

// TestRegisterCommand はコマンド登録のテスト
func TestRegisterCommand(t *testing.T) {
	dispatcher := NewCommandDispatcher()

	// カスタムコマンドを登録
	customHandled := false
	customHandler := func(args []string) error {
		customHandled = true
		return nil
	}

	dispatcher.Register("custom", customHandler)

	// 登録したコマンドが実行できることを確認
	handled, err := dispatcher.Dispatch("custom", []string{})
	if !handled {
		t.Error("カスタムコマンドがハンドルされませんでした")
	}
	if err != nil {
		t.Errorf("カスタムコマンドでエラーが発生しました: %v", err)
	}
	if !customHandled {
		t.Error("カスタムハンドラーが呼び出されませんでした")
	}
}

// TestParseCommand はコマンドパースのテスト
func TestParseCommand(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantCommand string
		wantArgs    []string
	}{
		{
			name:        "コマンドのみ",
			input:       "help",
			wantCommand: "help",
			wantArgs:    []string{},
		},
		{
			name:        "コマンドと引数",
			input:       "task new-feature",
			wantCommand: "task",
			wantArgs:    []string{"new-feature"},
		},
		{
			name:        "複数の引数",
			input:       "exec command arg1 arg2",
			wantCommand: "exec",
			wantArgs:    []string{"command", "arg1", "arg2"},
		},
		{
			name:        "前後の空白",
			input:       "  clear  ",
			wantCommand: "clear",
			wantArgs:    []string{},
		},
		{
			name:        "空文字列",
			input:       "",
			wantCommand: "",
			wantArgs:    []string{},
		},
		{
			name:        "スペースのみ",
			input:       "   ",
			wantCommand: "",
			wantArgs:    []string{},
		},
	}

	dispatcher := NewCommandDispatcher()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			command, args := dispatcher.ParseCommand(tt.input)

			if command != tt.wantCommand {
				t.Errorf("ParseCommand() command = %v, want %v", command, tt.wantCommand)
			}

			if len(args) != len(tt.wantArgs) {
				t.Errorf("ParseCommand() args length = %v, want %v", len(args), len(tt.wantArgs))
			} else {
				for i, arg := range args {
					if arg != tt.wantArgs[i] {
						t.Errorf("ParseCommand() args[%d] = %v, want %v", i, arg, tt.wantArgs[i])
					}
				}
			}
		})
	}
}

// TestDispatchToModel はModelへのディスパッチテスト
func TestDispatchToModel(t *testing.T) {
	model := NewModel()
	dispatcher := NewCommandDispatcher()

	// clearコマンドのテスト
	cmd := dispatcher.DispatchToModel(&model, "clear", []string{})
	if cmd != nil {
		t.Error("clearコマンドは即座に実行されるべきでCmdを返すべきではありません")
	}

	// exitコマンドのテスト
	cmd = dispatcher.DispatchToModel(&model, "exit", []string{})
	if cmd == nil {
		t.Error("exitコマンドはtea.Quitを返すべきです")
	}

	// 無効なコマンドのテスト
	cmd = dispatcher.DispatchToModel(&model, "invalid", []string{})
	if cmd != nil {
		t.Error("無効なコマンドはnilを返すべきです")
	}
}