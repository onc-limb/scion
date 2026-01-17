# 共通ユーティリティ・依存関係TODO

## プロジェクト構造
- [ ] ディレクトリ構造の作成
  ```
  scion/
  ├── cmd/
  │   └── scion/
  │       └── main.go
  ├── internal/
  │   ├── config/
  │   ├── git/
  │   ├── worktree/
  │   └── ui/
  ├── pkg/
  │   └── utils/
  └── test/
  ```

## 外部依存関係
- [ ] CLIフレームワーク
  - [ ] cobra導入 (`github.com/spf13/cobra`)
  - [ ] viper検討 (`github.com/spf13/viper`)
- [ ] TOML処理
  - [ ] go-toml導入 (`github.com/pelletier/go-toml/v2`)
- [ ] カラー出力
  - [ ] fatih/color導入 (`github.com/fatih/color`)
- [ ] ロギング
  - [ ] logrus検討 (`github.com/sirupsen/logrus`)
  - [ ] zap検討 (`go.uber.org/zap`)

## Git操作ユーティリティ（internal/git）
- [ ] Gitコマンド実行ラッパー
  - [ ] `exec.Command`のラッパー
  - [ ] 出力のキャプチャ
  - [ ] エラーハンドリング
- [ ] リポジトリ情報取得
  - [ ] リポジトリルート取得
  - [ ] 現在のブランチ取得
  - [ ] リモート情報取得
- [ ] Worktree操作
  - [ ] worktree list解析
  - [ ] worktree add/remove実行
  - [ ] ステータスチェック

## 設定管理ユーティリティ（internal/config）
- [ ] ConfigManagerインターフェース
- [ ] ファイル読み書き処理
- [ ] 設定マージ処理
- [ ] 環境変数処理
- [ ] デフォルト値管理

## UI/UXユーティリティ（internal/ui）
- [ ] プロンプト機能
  - [ ] Yes/No確認
  - [ ] テキスト入力
- [ ] プログレス表示
  - [ ] スピナー
  - [ ] プログレスバー
- [ ] テーブル表示
- [ ] カラー管理
  - [ ] 成功（緑）
  - [ ] エラー（赤）
  - [ ] 警告（黄）
  - [ ] 情報（青）

## ファイルシステムユーティリティ（pkg/utils）
- [ ] パス操作
  - [ ] ホームディレクトリ展開
  - [ ] 相対パス/絶対パス変換
- [ ] ディレクトリ操作
  - [ ] 存在確認
  - [ ] 作成（親ディレクトリ含む）
  - [ ] 削除（再帰的）
- [ ] ファイル操作
  - [ ] 読み書き
  - [ ] 権限チェック

## エラー処理（internal/errors）
- [ ] カスタムエラー型定義
- [ ] エラーコード定義
- [ ] エラーラップ処理
- [ ] ユーザーフレンドリーメッセージ生成

## バリデーション（internal/validation）
- [ ] ブランチ名バリデーション
- [ ] パスバリデーション
- [ ] 設定値バリデーション
- [ ] コマンド引数バリデーション

## テストユーティリティ（test/）
- [ ] モック作成
  - [ ] Gitコマンドモック
  - [ ] ファイルシステムモック
- [ ] テストヘルパー
  - [ ] 一時ディレクトリ作成
  - [ ] テスト用リポジトリ作成
- [ ] フィクスチャ管理
- [ ] ベンチマーク用ユーティリティ

## ビルド・配布
- [ ] Makefile作成
  - [ ] build target
  - [ ] test target
  - [ ] install target
  - [ ] clean target
- [ ] バージョン管理
  - [ ] バージョン番号埋め込み
  - [ ] ビルド時刻埋め込み
- [ ] クロスコンパイル設定
  - [ ] Linux
  - [ ] macOS
  - [ ] Windows
- [ ] リリーススクリプト

## ドキュメンテーション
- [ ] GoDocコメント
- [ ] README.md更新
- [ ] CONTRIBUTING.md作成
- [ ] CHANGELOG.md管理

## CI/CD
- [ ] GitHub Actions設定
  - [ ] ビルドワークフロー
  - [ ] テストワークフロー
  - [ ] リリースワークフロー
- [ ] カバレッジ測定
- [ ] 静的解析
  - [ ] golangci-lint設定
  - [ ] go vet

## パフォーマンス最適化
- [ ] 並行処理の実装
- [ ] キャッシュ機構
- [ ] メモリ使用量最適化
- [ ] 起動時間最適化