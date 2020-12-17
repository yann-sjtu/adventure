package logger

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/types"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// only used in the period of excuting tx, not in other module
var (
	Logger *zap.Logger
)

func InitLogger(logLevel int8, logListenUrl string) {
	Logger = initLogger(logLevel, logListenUrl)
}

func initLogger(logLevel int8, logListenUrl string) *zap.Logger {
	alevel := zap.NewAtomicLevel()
	alevel.SetLevel(zapcore.Level(logLevel))
	http.HandleFunc("/handle/level", alevel.ServeHTTP)
	go func() {
		if err := http.ListenAndServe(logListenUrl, nil); err != nil {
			panic(err)
		}
	}()

	logcfg := zap.NewDevelopmentConfig()
	logcfg.Development = false
	logcfg.DisableStacktrace = true
	logcfg.Level = alevel
	logcfg.EncoderConfig = zapcore.EncoderConfig{
		TimeKey:       "T",
		LevelKey:      "L",
		NameKey:       "N",
		MessageKey:    "M",
		StacktraceKey: "S",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("[2006-01-02 15:04:05]"))
		},
		EncodeDuration: zapcore.StringDurationEncoder,
	}
	logger, err := logcfg.Build()
	if err != nil {
		log.Fatal("err", err)
	}
	return logger
}

func PrintQueryAccountError(err error, txtype string, info keys.Info) {
	Logger.Debug("failed to query account",
		zap.String("types", txtype),
		zap.String("address", info.GetAddress().String()),
		zap.String("period", "query"),
		zap.String("error", convertError(err)),
	)
}

func PrintQueryProductsError(err error, txtype string, info keys.Info) {
	if err != nil {
		Logger.Debug("failed to query token_pairs",
			zap.String("types", txtype),
			zap.String("address", info.GetAddress().String()),
			zap.String("period", "query"),
			zap.String("error", convertError(err)),
		)
	} else {
		Logger.Debug("there is no token_pair    ",
			zap.String("types", txtype),
			zap.String("address", info.GetAddress().String()),
			zap.String("period", "query"),
		)
	}
}

func PrintQueryTokensError(err error, txtype string, info keys.Info) {
	if err != nil {
		Logger.Debug("failed to query tokens    ",
			zap.String("types", txtype),
			zap.String("address", info.GetAddress().String()),
			zap.String("period", "query"),
			zap.String("error", convertError(err)),
		)
	} else {
		Logger.Debug("there is no tokens        ",
			zap.String("types", txtype),
			zap.String("address", info.GetAddress().String()),
			zap.String("period", "query"),
		)
	}
}

func PrintExecuteTxError(err error, txtype string, info keys.Info) {
	Logger.Debug("failed to execute tx          ",
		zap.String("types", txtype),
		zap.String("address", info.GetAddress().String()),
		zap.String("period", "execute"),
		zap.String("error", convertError(err)),
	)
}

func PrintExecuteTxSuccess(txtype string, info keys.Info) {
	Logger.Info("successfully execute tx       ",
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
