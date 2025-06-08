package logger

func GinDebugPrintRoute(httpMethod, absolutePath, handlerName string, nuHandlers int) {
	Get().Info().
		Str("method", httpMethod).
		Str("path", absolutePath).
		Str("handler", handlerName).
		Int("handlers_count", nuHandlers).
		Send()
}

func GinDebugPrint(format string, values ...any) {
	Get().Info().Msgf(format, values...)
}
