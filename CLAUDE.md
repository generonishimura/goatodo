# CLAUDE.md

## Project Overview

goatodo - 習慣化×爆速操作のデスクトップTODOアプリ。Wails v2 + Svelte + Go + SQLite。

## Build & Test Commands

```bash
# Go全テスト
go test ./...

# カバレッジ付き（ドメイン＋アプリケーション層）
go test ./domain/... ./application/... -cover

# SQLite結合テスト
go test ./infrastructure/persistence/sqlite/...

# フロントエンドビルド
cd frontend && npm run build

# Goビルド確認
go build -o /dev/null .

# Wails開発サーバー
wails dev

# Wailsビルド
wails build
```

## Architecture

DDD/Clean Architecture。依存方向: Frontend → Presenter → Application → Domain ← Infrastructure

- **domain/**: 外部依存ゼロ。Result[T]型でエラーハンドリング。例外スローしない
- **application/**: ユースケース。ドメインリポジトリIFに依存
- **infrastructure/**: SQLite実装 + テスト用インメモリ実装
- **presenter/**: Wailsバインディング。DTO変換 + ユースケース呼び出し
- **frontend/**: Svelte。Wails生成のJSバインディング経由でGoを呼ぶ

## Key Conventions

- Go module: `github.com/i-nishimura/goatodo`
- DB: `~/.goatodo/goatodo.db` (SQLite, modernc.org/sqlite Pure Go)
- フロントエンドバインディング: `frontend/wailsjs/go/presenter/TaskHandler.js`
- テスト: TDD (Red→Green→Refactor)。ドメイン層カバレッジ80%以上必須
