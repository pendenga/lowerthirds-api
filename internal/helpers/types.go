package helpers

// ContextKey defines a unique type for a request context key, to ensure type unique-ness
type ContextKey string

const AuthKey = ContextKey("authorization")
const ErrorsResponseKey = ContextKey("errorResponse")
const QueryParametersKey = ContextKey("queryParameters")
const RequestIDKey = ContextKey("requestID")
const SocialIDKey = ContextKey("socialID")
const UsernameKey = ContextKey("username")
const UserIDKey = ContextKey("userID")
