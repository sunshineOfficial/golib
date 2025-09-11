package apierr

import "strings"

// Scope
// Определяет область обращения, при котором произошла ошибка
type Scope string

// Entity
// Определяет сущность, над которой совершались действия, приведшие к ошибке
type Entity string

// Reason
// Определяет причину возникновения ошибки
type Reason string

const (
	CodeDelimiter = "."

	// ScopeInternal
	// Ошибка внутри вызываемого сервиса
	ScopeInternal Scope = "internal"

	// ScopeExternal
	// Сервисы в нашем кластере и сервисы вне (s3, API поставщика и т. д.)
	ScopeExternal Scope = "external"

	// EntityRuntime
	// Сущность runtime или сам сервис, для ошибок технического характера (например, паники)
	EntityRuntime Entity = "runtime"

	// EntityAuth
	// Сущность авторизации, для ошибок аутентификации и авторизации
	EntityAuth Entity = "auth"

	// ReasonUnknown
	// Причина ошибки неизвестна
	ReasonUnknown Reason = "unknown"

	// ReasonNotFound
	// Причина ошибки - сущность не найдена
	ReasonNotFound Reason = "notFound"
)

func Build(scope Scope, entity Entity, reason Reason, flags ...string) string {
	if len(scope) == 0 {
		scope = ScopeInternal
	}
	if len(entity) == 0 {
		entity = EntityRuntime
	}
	if len(reason) == 0 {
		reason = ReasonUnknown
	}

	parts := make([]string, 0, 3+len(flags))
	parts = append(parts, string(scope), string(entity), string(reason))
	parts = append(parts, flags...)

	return strings.Join(parts, CodeDelimiter)
}

const (
	CodePanic        = "internal.runtime.unexpectedError"
	CodeUnknownError = "internal.runtime.unknownError"
	CodeForbidden    = "internal.auth.forbidden"
	CodeUnauthorized = "internal.auth.unauthorized"
)
