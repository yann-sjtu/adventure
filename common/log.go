package common

import (
	"fmt"
	"log"
	"strings"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/types"
	"go.uber.org/zap"
)

func PrintQueryAccountError(err error, txtype string, info keys.Info) {
	log.Println("failed to query account",
		zap.String("types", txtype),
		zap.String("address", info.GetAddress().String()),
		zap.String("period", "query"),
		zap.String("error", convertError(err)),
	)
}

func PrintQueryProductsError(err error, txtype string, info keys.Info) {
	if err != nil {
		log.Println("failed to query token_pairs",
			zap.String("types", txtype),
			zap.String("address", info.GetAddress().String()),
			zap.String("period", "query"),
			zap.String("error", convertError(err)),
		)
	} else {
		log.Println("there is no token_pair    ",
			zap.String("types", txtype),
			zap.String("address", info.GetAddress().String()),
			zap.String("period", "query"),
		)
	}
}

func PrintQueryTokensError(err error, txtype string, info keys.Info) {
	if err != nil {
		log.Println("failed to query tokens    ",
			zap.String("types", txtype),
			zap.String("address", info.GetAddress().String()),
			zap.String("period", "query"),
			zap.String("error", convertError(err)),
		)
	} else {
		log.Println("there is no tokens        ",
			zap.String("types", txtype),
			zap.String("address", info.GetAddress().String()),
			zap.String("period", "query"),
		)
	}
}

func PrintExecuteTxError(err error, txtype string, info keys.Info) {
	log.Println("failed to execute tx          ",
		zap.String("types", txtype),
		zap.String("address", info.GetAddress().String()),
		zap.String("period", "execute"),
		zap.String("error", convertError(err)),
	)
}

func PrintExecuteTxSuccess(txtype string, info keys.Info) {
	log.Println("successfully execute tx       ",
		zap.String("types", txtype),
		zap.String("address", info.GetAddress().String()),
	)
}

var replacer = strings.NewReplacer("\\", "", "\"", "")

func convertError(err error) string {
	return replacer.Replace(err.Error())
}

func convertRes(r *types.TxResponse) {
	var sb strings.Builder
	if _, err := sb.WriteString("Response:\n"); err != nil {
		log.Println(err)
	}

	if r.Height > 0 {
		if _, err := sb.WriteString(fmt.Sprintf("Height: %d\n", r.Height)); err != nil {
			log.Println(err)
		}
	}

	if r.TxHash != "" {
		if _, err := sb.WriteString(fmt.Sprintf("  TxHash: %s\n", r.TxHash)); err != nil {
			log.Println(err)
		}
	}

	if r.Code > 0 {
		if _, err := sb.WriteString(fmt.Sprintf("  Code: %d\n", r.Code)); err != nil {
			log.Println(err)
		}
	}

	if r.Data != "" {
		if _, err := sb.WriteString(fmt.Sprintf("  Data: %s\n", r.Data)); err != nil {
			log.Println(err)
		}
	}

	if r.RawLog != "" {
		if _, err := sb.WriteString(fmt.Sprintf("  Raw Log: %s\n", r.RawLog)); err != nil {
			log.Println(err)
		}
	}

	if r.Logs != nil {
		if _, err := sb.WriteString(fmt.Sprintf("  Logs: %s\n", r.Logs)); err != nil {
			log.Println(err)
		}

	}

	if r.Info != "" {
		if _, err := sb.WriteString(fmt.Sprintf("  Info: %s\n", r.Info)); err != nil {
			log.Println(err)
		}
	}

	//if r.GasWanted != 0 {
	//	sb.WriteString(fmt.Sprintf("  GasWanted: %d\n", r.GasWanted))
	//}
	//
	//if r.GasUsed != 0 {
	//	sb.WriteString(fmt.Sprintf("  GasUsed: %d\n", r.GasUsed))
	//}

	if r.Codespace != "" {
		if _, err := sb.WriteString(fmt.Sprintf("  Codespace: %s\n", r.Codespace)); err != nil {
			log.Println(err)
		}
	}

	if r.Timestamp != "" {
		if _, err := sb.WriteString(fmt.Sprintf("  Timestamp: %s\n", r.Timestamp)); err != nil {
			log.Println(err)
		}
	}

	if len(r.Logs) > 0 {
		if _, err := sb.WriteString(fmt.Sprintf("  Events: \n%s\n", r.Logs.String())); err != nil {
			log.Println(err)
		}
	}

}
