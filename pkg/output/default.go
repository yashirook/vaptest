package output

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/yashirook/vaptest/pkg/validator"
)

type DefaultFormatter struct {
}

func NewDefaultFormatter() *DefaultFormatter {
	return &DefaultFormatter{}
}

func (d *DefaultFormatter) Format(results []validator.ValidationResult) error {
	writer := tabwriter.NewWriter(
		// 標準出力を指定
		os.Stdout,
		// タブ幅、余白、パディング、パディング文字を指定
		0, 0, 2, ' ', 0,
	)

	if len(results) == 0 {
		fmt.Println("all validation success!")
		return nil
	}

	// ヘッダー行の出力
	fmt.Fprintln(writer, "API_VERSION\tKIND\tRESOURCE_NAME\tNAMESPACE\tVALIDATION_POLICY\tERRORS")

	// 各検証結果をテーブルに追加
	for _, result := range results {
		if result.Success {
			continue
		}

		// オブジェクト識別情報をまとめる
		obj := result.Target
		pol := result.Policy

		// オブジェクト識別情報をまとめる
		apiVersion := fmt.Sprintf("%s/%s",
			obj.APIGroup, obj.APIVersion,
		)

		// エラー内容をまとめる
		var errorDetails []string
		for _, err := range result.ValidationErrors {
			errorDetails = append(errorDetails, fmt.Sprintf("%s (Expression: %s)", err.Message, err.CELExpr))
		}
		errors := strings.Join(errorDetails, ", ")

		// エラーがない場合は空白を設定
		if result.Success {
			errors = "-"
		}

		// テーブル行を作成
		fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%v\t%s\n",
			apiVersion,
			obj.Kind,
			obj.ResourceName,
			obj.Namespace,
			pol.PolicyName,
			errors,
		)
	}

	// タブライターをフラッシュして出力
	writer.Flush()

	return nil
}
