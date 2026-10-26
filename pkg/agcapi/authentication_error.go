package agcapi

// AuthenticationError distinguishes a missing/incompatible credential from
// invalid business input, transport failure and Huawei business result codes.
type AuthenticationError struct{ Message string }

func (e *AuthenticationError) Error() string { return e.Message }
