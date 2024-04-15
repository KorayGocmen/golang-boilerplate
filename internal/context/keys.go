package context

type ContextKey string

const (
	// Top level context keys.
	KeyDate     ContextKey = "date"
	KeyLogLevel ContextKey = "log_level"
	KeyLang     ContextKey = "lang"

	// Server context keys.
	KeyMethod     ContextKey = "method"
	KeyPath       ContextKey = "path"
	KeyRequestID  ContextKey = "request_id"
	KeyQuery      ContextKey = "query"
	KeyReqBody    ContextKey = "req_body"
	KeyStatus     ContextKey = "status"
	KeyResBody    ContextKey = "res_body"
	KeyRemoteIP   ContextKey = "remote_ip"
	KeyError      ContextKey = "error"
	KeyErrorAPI   ContextKey = "error_api"
	KeyErrorStack ContextKey = "error_stack"

	// Database keys.
	KeyDBQuery   ContextKey = "db_query"
	KeyDBElapsed ContextKey = "db_elapsed"
	KeyDBRows    ContextKey = "db_rows"
	KeyDBFile    ContextKey = "db_file"

	// User context keys.
	KeyUserID        ContextKey = "user_id"
	KeyUserSessionID ContextKey = "user_session_id"
)

var (
	Keys = []ContextKey{
		// Top level context keys.
		KeyDate,
		KeyLogLevel,
		KeyLang,

		// Server context keys.
		KeyMethod,
		KeyPath,
		KeyRequestID,
		KeyQuery,
		KeyReqBody,
		KeyStatus,
		KeyResBody,
		KeyRemoteIP,
		KeyError,
		KeyErrorAPI,
		KeyErrorStack,

		// Database keys.
		KeyDBQuery,
		KeyDBElapsed,
		KeyDBRows,
		KeyDBFile,

		// User context keys.
		KeyUserID,
		KeyUserSessionID,
	}
)
