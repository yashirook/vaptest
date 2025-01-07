package output

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/yashirook/vaptest/pkg/validator"
)

type TableFormatter struct {
}

func NewTableFormatter() *TableFormatter {
	return &TableFormatter{}
}

func (d *TableFormatter) Output(results validator.ValidationResultList) error {
	writer := tabwriter.NewWriter(
		os.Stdout,
		0, 0, 2, ' ', 0,
	)

	if len(results.FailedResults()) == 0 {
		fmt.Println("all validation success!")
		return nil
	}

	fmt.Fprintln(writer, "POLICY\tEVALUATED_RESOURCE\tRESULT\tERRORS")

	for _, result := range results {
		if result.Success {
			continue
		}

		obj := result.Target
		pol := result.Policy

		resource := fmt.Sprintf("%s/%s",
			obj.Resource, obj.ResourceName,
		)

		var errorDetails []string
		for _, err := range result.ValidationErrors {
			errorDetails = append(errorDetails, fmt.Sprintf("%s (Expression: %s)", err.Message, err.CELExpr))
		}
		errors := strings.Join(errorDetails, ", ")

		if result.Success {
			errors = "-"
		}

		var res string
		if result.Success {
			res = "Pass"
		} else {
			res = "Fail"
		}

		fmt.Fprintf(writer, "%s\t%s\t%s\t%s\n",
			pol.PolicyName,
			resource,
			res,
			errors,
		)
	}

	writer.Flush()

	return nil
}
