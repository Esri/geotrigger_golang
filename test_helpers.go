package geotrigger_golang

import (
	"testing"
	"reflect"
	"errors"
	"fmt"
)

// Restorer holds a function that can be used
// to restore some previous state.
type restorer func()

// Restore restores some previous state.
func (r restorer) restore() {
	r()
}

// Patch sets the value pointed to by the given destination to the given
// value, and returns a function to restore it to its original value.  The
// value must be assignable to the element type of the destination.
func patch(destination, v interface{}) (restorer, error) {
	destType := reflect.TypeOf(destination)
	if reflect.TypeOf(destination).Kind() != reflect.Ptr {
		return nil, errors.New("Bad destination, please provide a pointer.")
	}

	// we know destination is a pointer, so get the type of value being pointed to
	destType = destType.Elem()
	// compare that type to the type of v
	providedType := reflect.TypeOf(v)
	if destType != providedType {
		return nil, errors.New(fmt.Sprintf("Provided value of type %s does not match destination type: %s.",
			providedType, destType))
	}

	// get the value being pointed to
	destValue := reflect.ValueOf(destination).Elem()
	// reflect.New creates a new pointer value to provided type, elem gets the pointed to value again
	oldValue := reflect.New(destType).Elem()
	// we then set that value to the current destination value to hold onto it
	oldValue.Set(destValue)

	// the value of the provided... value...
	value := reflect.ValueOf(v)
	if !value.IsValid() {
		// This should be a rare occurrence.
		// the value provided could not be reflected, and we have an invalid Value here
		// so just attempt to use the zero value for the destination type.
		value = reflect.Zero(destValue.Type())
	}

	// replace the destination's current value with the value of the provided v
	// this shouldn't panic, because we have already checked that they are the same type
	destValue.Set(value)
	return func() {
		// restore the destination's value to its original
		destValue.Set(oldValue)
	}, nil
}

func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func refute(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}