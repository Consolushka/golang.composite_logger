package internal

import (
	"fmt"
	"runtime"
	"strings"
)

func BuildErrorContextWithStackTrace(ctx map[string]interface{}) map[string]interface{} {
	context := CloneContext(ctx)

	if _, exists := context["stackTrace"]; exists {
		return context
	}

	stackTrace := ExtractStackTraceFromError(context["error"])
	if stackTrace == "" {
		stackTrace = BuildFallbackStackTrace()
	}

	context["stackTrace"] = stackTrace

	return context
}

func CloneContext(ctx map[string]interface{}) map[string]interface{} {
	if ctx == nil {
		return map[string]interface{}{}
	}

	cloned := make(map[string]interface{}, len(ctx))
	for key, value := range ctx {
		cloned[key] = value
	}

	return cloned
}

func ExtractStackTraceFromError(errValue interface{}) string {
	err, ok := errValue.(error)
	if !ok || err == nil {
		return ""
	}

	expanded := fmt.Sprintf("%+v", err)
	if expanded != "" && expanded != err.Error() {
		return expanded
	}

	return ""
}

func BuildFallbackStackTrace() string {
	const maxFrames = 48
	const skipFrames = 3

	pcs := make([]uintptr, maxFrames)
	count := runtime.Callers(skipFrames, pcs)
	if count == 0 {
		return ""
	}

	frames := runtime.CallersFrames(pcs[:count])
	lines := make([]string, 0, count)

	for {
		frame, more := frames.Next()
		if !more {
			if ShouldIncludeFrame(frame.Function) {
				lines = append(lines, FormatFrame(frame))
			}
			break
		}

		if ShouldIncludeFrame(frame.Function) {
			lines = append(lines, FormatFrame(frame))
		}
	}

	return strings.Join(lines, "\n")
}

func ShouldIncludeFrame(function string) bool {
	return !strings.Contains(function, "github.com/Consolushka/golang.composite_logger/pkg/") &&
		!strings.Contains(function, "github.com/Consolushka/golang.composite_logger/internal")
}

func FormatFrame(frame runtime.Frame) string {
	return fmt.Sprintf("%s\n\t%s:%d", frame.Function, frame.File, frame.Line)
}
