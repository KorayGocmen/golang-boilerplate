package context

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"github.com/koraygocmen/null"
)

func Map(ctx Ctx) map[ContextKey]interface{} {
	m := make(map[ContextKey]interface{})
	for _, key := range Keys {
		if val := ctx.Value(key); val != nil {
			switch valTyped := val.(type) {
			case null.String:
				val = valTyped.String
			case null.Bool:
				val = valTyped.Bool
			case null.Float:
				val = valTyped.Float64
			case null.Int:
				val = valTyped.Int64
			case net.IP:
				val = valTyped.String()
			}

			switch key {
			case KeyReqBody:
				marshalled, _ := json.Marshal(val)
				val = string(marshalled)
			}

			m[key] = val
		}
	}
	return m
}

func String(ctx Ctx, format string, v ...interface{}) string {
	m := Map(ctx)
	var prefixes []string
	for _, key := range Keys {
		if val, ok := m[key]; ok {
			prefixes = append(prefixes, fmt.Sprintf(`%s="%v"`, key, val))
		}
	}

	prefixes = append(prefixes, fmt.Sprintf(format, v...))
	return strings.Join(prefixes, " ")
}

func JSON(ctx Ctx, format string, v ...interface{}) string {
	m := Map(ctx)
	var prefixes []string
	for _, key := range Keys {
		if val, ok := m[key]; ok {
			prefixes = append(prefixes, fmt.Sprintf(`"%s": "%v"`, key, val))
		}
	}

	if len(v) > 0 {
		prefixes = append(prefixes, fmt.Sprintf(`"msg": "%s"`, fmt.Sprintf(format, v...)))
	}

	return fmt.Sprintf("{%s}", strings.Join(prefixes, ","))
}
