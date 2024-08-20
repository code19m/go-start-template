package errx

import "errors"

// ReplaceWith replaces the target error with a new error.
// If the given error is equal to the target error, it is replaced with the new error.
// Otherwise, the original error is returned unchanged.
//
// Example usage:
//
//	err := repo.FindOrder(ctx, orderID)
//	if err != nil {
//		return errx.ReplaceWith(err, errx.ErrNotFound, errx.ErrOrderNotFound)
//	}
func ReplaceWith(err error, targetErr error, newErr error) error {
	if errors.Is(err, targetErr) {
		return newErr
	}
	return err
}

// ReplaceWithCode replaces the error with the target code with a new error.
// If the given error has the target code, it is replaced with the new error.
// Otherwise, the original error is returned unchanged.
func ReplaceWithCode(err error, targetCode string, newErr error) error {
	if GetCode(err) == targetCode {
		return newErr
	}
	return err
}

// ReplaceWithType replaces the error with the target type with a new error.
// If the given error has the target type, it is replaced with the new error.
// Otherwise, the original error is returned unchanged.
func ReplaceWithType(err error, targetType Type, newErr error) error {
	if GetType(err) == targetType {
		return newErr
	}
	return err
}
