# goatodo

俺が考えた最強TODOアプリ。「習慣化 × 爆速操作」で最強を目指す、自分専用のデスクトップTODOアプリ。

## 何において最強か

1. **ゼロフリクション**: グローバルホットキーで即キャプチャ。アプリを「開く」必要すらない
2. **習慣エンジン**: デイリーレビュー＋ストリーク＋Today表示で、使い続ける仕組みが内蔵
3. **爆速**: ネイティブWebView、SQLite、Go。起動もCRUDもサブ秒
4. **キーボードファースト**: マウス不要。全操作にショートカット

## Tech Stack

| レイヤー | 技術 | 理由 |
|---------|------|------|
| デスクトップ | Wails v2 | Go製、ネイティブWebView（Electronと違いChromium不要） |
| フロントエンド | Svelte | 軽量・高速・ボイラープレート少 |
| バックエンド | Go | 型安全、高速、クロスコンパイル容易 |
| DB | SQLite (modernc.org/sqlite) | ローカルファースト、ゼロ設定、Pure Go（CGO不要） |

## Architecture (DDD/Clean Architecture)

```
goatodo/
├── main.go / app.go              # Wails起動・DI
├── domain/                        # ドメイン層（外部依存ゼロ）
│   ├── task/                      # Task集約
│   │   ├── task.go                # エンティティ
│   │   ├── status.go              # Value Object (todo/doing/done)
│   │   ├── priority.go            # Value Object (none/low/medium/high)
│   │   ├── repository.go          # リポジトリIF
│   │   └── errors.go              # ドメインエラー
│   └── shared/
│       ├── result.go              # Result[T]型
│       └── id.go                  # UUID生成
├── application/                   # アプリケーション層（ユースケース）
│   └── task/
│       ├── create_task.go
│       ├── complete_task.go
│       ├── list_tasks.go
│       ├── update_task.go
│       └── delete_task.go
├── infrastructure/                # インフラ層
│   └── persistence/
│       ├── sqlite/                # SQLite実装
│       └── memory/                # テスト用インメモリ実装
├── presenter/                     # プレゼンテーション層（Wailsバインディング）
│   ├── task_handler.go
│   └── dto/
└── frontend/                      # Svelte
    └── src/
        ├── App.svelte
        └── lib/
            ├── components/        # TaskList, TaskItem, AddTask, StatusBar
            └── stores/
```

**依存の方向**: Frontend → Presenter → Application → Domain ← Infrastructure

## キーボードショートカット

| キー | 操作 |
|------|------|
| `j` / `k` | タスク間を上下移動 |
| `a` | 新規タスク追加 |
| `e` / `Enter` | タスク名を編集 |
| `x` | ステータスを次に進める (todo → doing → done) |
| `d` | タスクを削除 |
| `0` - `3` | 優先度設定 (0=none, 1=low, 2=medium, 3=high) |
| `g` / `G` | 先頭 / 末尾に移動 |
| `Esc` | 入力キャンセル |

## 開発

### 前提条件

- Go 1.23+
- Node.js 20+
- [Wails CLI v2](https://wails.io/docs/gettingstarted/installation)

### セットアップ

```bash
# Wails CLIのインストール
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# フロントエンド依存のインストール
cd frontend && npm install && cd ..

# 開発サーバー起動
wails dev

# ビルド
wails build
```

### テスト

```bash
# 全テスト実行
go test ./...

# カバレッジ付き
go test ./domain/... ./application/... -cover
```

### データ保存先

`~/.goatodo/goatodo.db` (SQLite)

## Roadmap

- [x] Phase 1: MVP（コアCRUD + 基本UI + キーボード操作 + SQLite永続化）
- [ ] Phase 2: 習慣エンジン（グローバルホットキー、デイリーレビュー、ストリーク、システムトレイ）
- [ ] Phase 3: パワー機能（タグ、繰り返しタスク、期限通知、全文検索）
