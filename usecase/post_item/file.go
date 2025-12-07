package post_item

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

// saveUploadedFile: アップロードされたファイルを保存し、保存パスを返す
func saveUploadedFile(file io.Reader, fileHeader *multipart.FileHeader) (string, error) {

	// ファイルがアップロードされていない場合はスキップ
	if file == nil {
		return "", nil
	}

	// 💡 注意：セキュリティリスク軽減のため、ファイル名のサニタイズ（安全な処理）が必要です。
	// 今回は例として前回のロジックをほぼ踏襲しますが、本番環境では必ず対策をしてください。

	defer file.(multipart.File).Close() // multipart.File は io.Closerでもあるため、型アサーションが必要な場合がある

	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("could not create directory: %w", err)
	}

	// 簡易的なファイル名生成 (セキュリティ上の注意点あり)
	// 【セキュリティ修正の提案】ファイル名からディレクトリ情報を除去
	safeFilename := filepath.Base(fileHeader.Filename)
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), safeFilename)

	filePath := filepath.Join(uploadDir, filename)

	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("could not create file on disk: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("could not save file content: %w", err)
	}

	return filePath, nil
}
