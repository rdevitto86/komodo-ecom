package context

type ctxKey string

const (
	START_TIME_KEY 				ctxKey = "start_time"
	END_TIME_KEY   				ctxKey = "end_time"
	DURATION_KEY  				ctxKey = "duration"
	VERSION_KEY 					ctxKey = "version"
	URI_KEY       				ctxKey = "uri"
	PATH_PARAMS_KEY 			ctxKey = "path_params"
	QUERY_PARAMS_KEY 			ctxKey = "query_params"
	VALIDATION_RULE_KEY 	ctxKey = "validation_rule"
	REQUEST_ID_KEY     		ctxKey = "request_id"
	REQUEST_TIMEOUT_KEY   ctxKey = "request_timeout"
	CLIENT_IP_KEY       	ctxKey = "client_ip"
	CLIENT_TYPE_KEY			 	ctxKey = "client_type"
	USER_AGENT_KEY     		ctxKey = "user_agent"
	METHOD_KEY         		ctxKey = "method"
	AUTH_VALID_KEY			 	ctxKey = "auth_valid"
	SESSION_ID_KEY       	ctxKey = "session_id"
	SESSION_VALID_KEY   	ctxKey = "session_valid"
	USER_ID_KEY           ctxKey = "user_id"
	REQUEST_TYPE_KEY      ctxKey = "request_type"
	SCOPES_KEY            ctxKey = "scopes"
	IS_ADMIN_KEY          ctxKey = "is_admin"
	IDEMPOTENCY_KEY       ctxKey = "idempotency_key"
	IDEMPOTENCY_VALID_KEY ctxKey = "idempotency_key_valid"
	CSRF_TOKEN_KEY      	ctxKey = "csrf_token"
	CSRF_VALID_KEY      	ctxKey = "csrf_token_valid"
	LOGGER_KEY          	ctxKey = "logger"
	OTEL_LOGGER_KEY     	ctxKey = "otel_logger"
)
