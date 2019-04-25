package errors

import "github.com/lab259/errors"

var (
	CommonModule = errors.Module("common")
	TodoModule   = errors.Module("todos")
)

var (
	// 	 ____
	// 	/ ___|___  _ __ ___  _ __ ___   ___  _ __
	// | |   / _ \| '_ ` _ \| '_ ` _ \ / _ \| '_ \
	// | |__| (_) | | | | | | | | | | | (_) | | | |
	// 	\____\___/|_| |_| |_|_| |_| |_|\___/|_| |_|

	ErrNotFound = errors.Wrap(errors.New("not found"), CommonModule, errors.Code("not-found"), errors.Message("We could not find the resource you requested."))

	//  _____         _
	// |_   _|__   __| | ___  ___
	//   | |/ _ \ / _` |/ _ \/ __|
	//   | | (_) | (_| | (_) \__ \
	//   |_|\___/ \__,_|\___/|___/

	ErrDescriptionRequired = errors.Wrap(errors.New("description is required"), TodoModule, errors.Code("description-required"), errors.Message("You must provide a description for a Todo."))
)
