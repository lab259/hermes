package errors

import "github.com/lab259/errors"

var (
	TodoModule = errors.Module("todos")
)

var (

	//  _____         _
	// |_   _|__   __| | ___  ___
	//   | |/ _ \ / _` |/ _ \/ __|
	//   | | (_) | (_| | (_) \__ \
	//   |_|\___/ \__,_|\___/|___/

	ErrTodoNotFound        = errors.Wrap(errors.New("todo not found"), TodoModule, errors.Code("todo-not-found"), errors.Message("We could not find the Todo you requested."))
	ErrDescriptionRequired = errors.Wrap(errors.New("description is required"), TodoModule, errors.Code("description-required"), errors.Message("You must provide a description for a Todo."))
)
