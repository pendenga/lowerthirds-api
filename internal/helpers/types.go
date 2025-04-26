package helpers

// ContextKey defines a unique type for a request context key, to ensure type unique-ness
type ContextKey string

const ErrorsResponseKey = ContextKey("errorResponse")
const RequestIDKey = ContextKey("requestID")
const AccountIDKey = ContextKey("accountID")
const AuthKey = ContextKey("authorization")
const UsernameKey = ContextKey("username")
const UserIDKey = ContextKey("userID")
const StartDateKey = ContextKey("startDate")
const EndDateKey = ContextKey("endDate")
const GeometryKey = ContextKey("geometry")
const DeviceIDsKey = ContextKey("deviceIDs")
const DeviceIDsHideLocationMapKey = ContextKey("deviceIDsHideLocationMap")
