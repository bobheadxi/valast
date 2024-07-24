package customtype

import (
	"fmt"
	"go/ast"
	"reflect"
	"sync"
)

var (
	customTypesMux sync.Mutex
	customTypes    = make(map[reflect.Type]func(any) ast.Expr)
)

// Register registers a type that for representation in a custom manner with
// valast. If valast encounters a value or pointer to a value of this type, it
// will use the given render func to generate the appropriate AST representation.
//
// This is useful if a type's fields are private, and can only be represented
// through a constructor - see stdtypes.go for examples.
//
// This mechanism currently only works with struct types.
func Register[T any](render func(value T) ast.Expr) {
	customTypesMux.Lock()
	var zero T
	t := reflect.TypeOf(zero)
	if _, exists := customTypes[t]; exists {
		panic(fmt.Sprintf("%T already registered", zero))
	}
	customTypes[t] = func(value any) ast.Expr { return render(value.(T)) }
	customTypesMux.Unlock()
}

// Is indicates if the given reflect.Type has a custom AST representation
// generator registered.
func Is(rt reflect.Type) (func(any) ast.Expr, bool) {
	customTypesMux.Lock()
	defer customTypesMux.Unlock()

	t, ok := customTypes[rt]
	return t, ok
}
