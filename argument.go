package graphb

import (
	"fmt"
	"time"
)

type argumentValue interface {
	stringChan() <-chan string
}

type Argument struct {
	Name  string
	Value argumentValue
}

func (a *Argument) stringChan() <-chan string {
	tokenChan := make(chan string)
	go func() {
		tokenChan <- a.Name
		tokenChan <- ":"
		for str := range a.Value.stringChan() {
			tokenChan <- str
		}
		close(tokenChan)
	}()
	return tokenChan
}

func ArgumentAny(name string, value interface{}) (Argument, error) {
	switch v := value.(type) {
	case bool:
		return ArgumentBool(name, v), nil
	case []bool:
		return ArgumentBoolSlice(name, v...), nil

	case int:
		return ArgumentInt(name, v), nil
	case []int:
		return ArgumentIntSlice(name, v...), nil

	case string:
		return ArgumentString(name, v), nil
	case []string:
		return ArgumentStringSlice(name, v...), nil

	case time.Time:
		return ArgumentTime(name, v), nil

	default:
		return Argument{}, ArgumentTypeNotSupportedErr{Value: value}
	}
}

func ArgumentBool(name string, value bool) Argument {
	return Argument{name, argBool(value)}
}

func ArgumentInt(name string, value int) Argument {
	return Argument{name, argInt(value)}
}

func ArgumentString(name string, value string) Argument {
	return Argument{name, argString(value)}
}

func ArgumentQuotedString(name string, value string) Argument {
	return Argument{name, argQuotedString(value)}
}

func ArgumentBlockString(name string, value string) Argument {
	return Argument{name, argBlockString(value)}
}

func ArgumentEnum(name string, value string) Argument {
	return Argument{name, argEnum(value)}
}

func ArgumentTime(name string, value time.Time) Argument {
	return Argument{name, argTime(value)}
}

func ArgumentBoolSlice(name string, values ...bool) Argument {
	return Argument{name, argBoolSlice(values)}
}

func ArgumentIntSlice(name string, values ...int) Argument {
	return Argument{name, argIntSlice(values)}
}

func ArgumentStringSlice(name string, values ...string) Argument {
	return Argument{name, argStringSlice(values)}
}

func ArgumentEnumSlice(name string, values ...string) Argument {
	return Argument{name, argEnumSlice(values)}
}

// ArgumentCustomType returns a custom GraphQL type's argument representation, which could be a recursive structure.
func ArgumentCustomType(name string, values ...Argument) Argument {
	return Argument{name, argumentCustom(values)}
}

// ArgumentSlice returns a slice of arguments
func ArgumentSlice(name string, values ...[]Argument) Argument {
	return Argument{name, argArgSlice(values)}
}

/////////////////////////////
// Primitive Wrapper Types //
/////////////////////////////

// argBool represents a boolean value.
type argBool bool

func (v argBool) stringChan() <-chan string {
	tokenChan := make(chan string)
	go func() {
		tokenChan <- fmt.Sprintf("%t", v)
		close(tokenChan)
	}()
	return tokenChan
}

// argInt represents an integer value.
type argInt int

func (v argInt) stringChan() <-chan string {
	tokenChan := make(chan string)
	go func() {
		tokenChan <- fmt.Sprintf("%d", v)
		close(tokenChan)
	}()
	return tokenChan
}

// argString represents a string value.
type argString string

func (v argString) stringChan() <-chan string {
	tokenChan := make(chan string)
	go func() {
		tokenChan <- fmt.Sprintf(`"%s"`, v)
		close(tokenChan)
	}()
	return tokenChan
}

// argQuotedString represents a quoted string value.
type argQuotedString string

func (v argQuotedString) stringChan() <-chan string {
	tokenChan := make(chan string)
	go func() {
		tokenChan <- fmt.Sprintf(`"\\"%s\\""`, v)
		close(tokenChan)
	}()
	return tokenChan
}

// argBlockString represents a block string value.
type argBlockString string

func (v argBlockString) stringChan() <-chan string {
	tokenChan := make(chan string)
	go func() {
		tokenChan <- fmt.Sprintf(`"""%s"""`, v)
		close(tokenChan)
	}()
	return tokenChan
}

// argEnum represents a enum value.
type argEnum string

func (v argEnum) stringChan() <-chan string {
	tokenChan := make(chan string)
	go func() {
		tokenChan <- fmt.Sprintf("%s", v)
		close(tokenChan)
	}()
	return tokenChan
}

// argTime represents a time value
type argTime time.Time

func (v argTime) stringChan() <-chan string {
	tokenChan := make(chan string)
	go func() {
		tokenChan <- fmt.Sprintf(`"%s"`, time.Time(v).Format(time.RFC3339))
		close(tokenChan)
	}()
	return tokenChan
}

//////////////////////////////////
// Primitive List Wrapper Types //
//////////////////////////////////

// argBoolSlice implements valueSlice
type argBoolSlice []bool

func (s argBoolSlice) stringChan() <-chan string {
	tokenChan := make(chan string)
	go func() {
		tokenChan <- "["
		for i, v := range s {
			if i != 0 {
				tokenChan <- ","
			}
			tokenChan <- fmt.Sprintf("%t", v)
		}
		tokenChan <- "]"
		close(tokenChan)
	}()
	return tokenChan
}

// argIntSlice implements valueSlice
type argIntSlice []int

func (s argIntSlice) stringChan() <-chan string {
	tokenChan := make(chan string)
	go func() {
		tokenChan <- "["
		for i, v := range s {
			if i != 0 {
				tokenChan <- ","
			}
			tokenChan <- fmt.Sprintf("%d", v)
		}
		tokenChan <- "]"
		close(tokenChan)
	}()
	return tokenChan
}

// argStringSlice implements valueSlice
type argStringSlice []string

func (s argStringSlice) stringChan() <-chan string {
	tokenChan := make(chan string)
	go func() {
		tokenChan <- "["
		for i, v := range s {
			if i != 0 {
				tokenChan <- ","
			}
			tokenChan <- fmt.Sprintf(`"%s"`, v)
		}
		tokenChan <- "]"
		close(tokenChan)
	}()
	return tokenChan
}

// argEnumSlice implements valueSlice
type argEnumSlice []string

func (s argEnumSlice) stringChan() <-chan string {
	tokenChan := make(chan string)
	go func() {
		tokenChan <- "["
		for i, v := range s {
			if i != 0 {
				tokenChan <- ","
			}
			tokenChan <- fmt.Sprintf("%s", v)
		}
		tokenChan <- "]"
		close(tokenChan)
	}()
	return tokenChan
}

type argumentCustom []Argument

func (s argumentCustom) stringChan() <-chan string {
	tokenChan := make(chan string)
	go func() {
		tokenChan <- "{"
		for i, v := range s {
			if i != 0 {
				tokenChan <- ","
			}
			for str := range v.stringChan() {
				tokenChan <- str
			}
		}
		tokenChan <- "}"
		close(tokenChan)
	}()
	return tokenChan
}

type argArgSlice [][]Argument

func (s argArgSlice) stringChan() <-chan string {
	tokenChan := make(chan string)
	go func() {
		tokenChan <- "["
		for i, v := range s {
			if i != 0 {
				tokenChan <- ","
			}
			for str := range argumentCustom(v).stringChan() {
				tokenChan <- str
			}
		}
		tokenChan <- "]"
		close(tokenChan)
	}()
	return tokenChan
}
