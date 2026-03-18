# セッション引き継ぎ

**日付**: 2026-03-18
**ブランチ**: `feat/project-scaffolding`
**PR**: https://github.com/generonishimura/goatodo/pull/1

## 完了した作業

- Wails v2 + Svelte テンプレートからプロジェクトスキャフォールド、DDD構造にリストラクチャ
- `domain/shared`: Result[T]ジェネリック型、UUID ID生成（カバレッジ100%）
- `domain/task`: Task集約エンティティ、Status遷移ルール、Priority Value Object、Repositoryインターフェース（カバレッジ95.8%）
- `application/task`: CreateTask / CompleteTask / ListTasks / UpdateTask / DeleteTask ユースケース（カバレッジ87.5%）
- `infrastructure/persistence/sqlite`: modernc.org/sqlite（Pure Go）によるTaskRepository + マイグレーション + 結合テスト
- `infrastructure/persistence/memory`: テスト用インメモリTaskRepository
- `presenter`: TaskHandler（Wailsバインディング）、DTO変換、DI配線（app.go）
- `frontend/src`: Svelte製キーボードファーストUI（TaskList, TaskItem, AddTask, StatusBar）
- README.md更新（アーキテクチャ図、ショートカット一覧、セットアップ手順、Roadmap）
- CLAUDE.md作成（ビルド/テストコマンド、規約）
- `wails dev` で起動確認済み、動作OK
- PR #1 作成済み（リモートにプッシュ済み）

## 未コミットの変更

- `wails dev` 実行時にWailsが自動生成したバインディングファイルの更新:
  - `frontend/wailsjs/go/presenter/TaskHandler.d.ts` — 型定義が正式な型に更新
  - `frontend/wailsjs/go/presenter/TaskHandler.js` — 同上
  - `frontend/wailsjs/go/models.ts` — 新規生成（DTOの型定義）
  - `frontend/wailsjs/go/main/App.d.ts` / `App.js` — 削除（旧テンプレートの残骸）
  - `frontend/wailsjs/runtime/runtime.d.ts` / `runtime.js` — 更新
  - `frontend/package.json.md5` — 新規生成
- これらは `wails dev` 時に自動再生成されるので、コミットしてもしなくてもOK

## 次にやること

1. **未コミット変更の整理**: Wails自動生成ファイルをコミットする（`chore: update wails-generated bindings`）
2. **PR #1 のマージ**: Phase 1 MVPをmainに取り込む
3. **Phase 2: 習慣エンジン** — 実装計画の Step 7, 8
   - `feat/quick-capture`: システムトレイ常駐 + グローバルホットキー（Cmd+Shift+T）
   - `feat/habit-engine`: DailyReview / Streak ドメインモデル + UI
4. **ライトモード切替**: ダーク/ライトモード対応（Phase 1の残タスク）

## 決定事項・コンテキスト

- **Go module名**: `github.com/i-nishimura/goatodo`（GitHub上は `generonishimura/goatodo`）
- **DB保存先**: `~/.goatodo/goatodo.db`
- **SQLite**: modernc.org/sqlite（Pure Go、CGO不要）を採用
- **ステータス遷移ルール**: todo→doing→done は自由、done からの逆戻りは不可
- **フロントエンドバインディング**: `wails dev` 実行時に `frontend/wailsjs/` 以下が自動再生成される
- **TDD**: 全ドメイン/アプリケーション層はRed→Green→Refactorで実装済み

## ブロッカー・注意点

- なし
