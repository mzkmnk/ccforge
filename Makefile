# ccforge Makefile
# 開発タスクの自動化

# 変数定義
APP_NAME := ccforge
GO := go
GOLANGCI_LINT := golangci-lint
BUILD_DIR := ./build
MAIN_FILE := main.go

# デフォルトタスク
.PHONY: help
help: ## ヘルプを表示
	@echo "利用可能なタスク:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ビルドタスク
.PHONY: build
build: ## アプリケーションをビルド
	@echo "ビルド中..."
	@mkdir -p $(BUILD_DIR)
	@$(GO) build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_FILE)
	@echo "ビルド完了: $(BUILD_DIR)/$(APP_NAME)"

.PHONY: build-debug
build-debug: ## デバッグ情報付きでビルド
	@echo "デバッグビルド中..."
	@mkdir -p $(BUILD_DIR)
	@$(GO) build -gcflags="all=-N -l" -o $(BUILD_DIR)/$(APP_NAME)-debug $(MAIN_FILE)
	@echo "デバッグビルド完了: $(BUILD_DIR)/$(APP_NAME)-debug"

# テストタスク
.PHONY: test
test: ## すべてのテストを実行
	@echo "テスト実行中..."
	@$(GO) test -v ./...

.PHONY: test-cover
test-cover: ## カバレッジ付きでテストを実行
	@echo "カバレッジ測定中..."
	@$(GO) test -cover -coverprofile=coverage.out ./...
	@$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "カバレッジレポート: coverage.html"

.PHONY: test-race
test-race: ## 競合状態の検出付きでテストを実行
	@echo "競合状態検出中..."
	@$(GO) test -race ./...

.PHONY: bench
bench: ## ベンチマークテストを実行
	@echo "ベンチマーク実行中..."
	@$(GO) test -bench=. -benchmem ./...

# リントタスク
.PHONY: lint
lint: ## リンターを実行
	@echo "リント実行中..."
	@if command -v $(GOLANGCI_LINT) > /dev/null; then \
		$(GOLANGCI_LINT) run ./...; \
	else \
		echo "golangci-lintがインストールされていません。"; \
		echo "インストール: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

.PHONY: fmt
fmt: ## コードをフォーマット
	@echo "フォーマット中..."
	@$(GO) fmt ./...
	@goimports -w .

# インストールタスク
.PHONY: install
install: build ## ビルドしてシステムにインストール
	@echo "インストール中..."
	@cp $(BUILD_DIR)/$(APP_NAME) $(GOPATH)/bin/$(APP_NAME)
	@echo "インストール完了: $(GOPATH)/bin/$(APP_NAME)"

.PHONY: install-tools
install-tools: ## 開発ツールをインストール
	@echo "開発ツールをインストール中..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@echo "開発ツールのインストール完了"

# 実行タスク
.PHONY: run
run: ## アプリケーションを実行
	@echo "実行中..."
	@$(GO) run $(MAIN_FILE)

.PHONY: run-debug
run-debug: ## デバッグモードで実行
	@echo "デバッグモードで実行中..."
	@$(GO) run -race $(MAIN_FILE)

# クリーンアップタスク
.PHONY: clean
clean: ## ビルド成果物をクリーンアップ
	@echo "クリーンアップ中..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "クリーンアップ完了"

# 依存関係管理
.PHONY: deps
deps: ## 依存関係を更新
	@echo "依存関係を更新中..."
	@$(GO) mod download
	@$(GO) mod tidy
	@echo "依存関係の更新完了"

.PHONY: deps-upgrade
deps-upgrade: ## 依存関係を最新版にアップグレード
	@echo "依存関係をアップグレード中..."
	@$(GO) get -u ./...
	@$(GO) mod tidy
	@echo "依存関係のアップグレード完了"

# CI/CD用タスク
.PHONY: ci
ci: deps fmt lint test ## CI環境用の完全チェック
	@echo "CI チェック完了"

# 初期セットアップ
.PHONY: setup
setup: deps install-tools ## 開発環境をセットアップ
	@echo "セットアップ完了"

# デフォルトタスク
.DEFAULT_GOAL := help