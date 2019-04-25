package errors

import "errors"

var (
	// 	 ____
	// 	/ ___|___  _ __ ___  _ __ ___   ___  _ __
	// | |   / _ \| '_ ` _ \| '_ ` _ \ / _ \| '_ \
	// | |__| (_) | | | | | | | | | | | (_) | | | |
	// 	\____\___/|_| |_| |_|_| |_| |_|\___/|_| |_|

	ErrNotFound = errors.New("not found")

	//  _____         _
	// |_   _|__   __| | ___  ___
	//   | |/ _ \ / _` |/ _ \/ __|
	//   | | (_) | (_| | (_) \__ \
	//   |_|\___/ \__,_|\___/|___/

	ErrDescriptionRequired = errors.New("description is required")
)
