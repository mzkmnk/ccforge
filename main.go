package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("ccforge - Claude Code TUIアプリケーション")
	fmt.Println("初期化中...")
	
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	fmt.Println("準備完了")
	return nil
}