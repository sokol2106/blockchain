package customerrors

import "errors"

var ErrNoDataToVerification = errors.New("no data available for processing verification")
