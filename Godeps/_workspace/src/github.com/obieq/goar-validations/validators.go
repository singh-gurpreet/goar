package goar_validations

import (
	"fmt"
	"reflect"
	"regexp"
	"time"
)

type Validator interface {
	IsSatisfied(interface{}) bool
	DefaultMessage() string
	GetKey() string
}

type Required struct {
	Key string
}

func (v Required) GetKey() string {
	return v.Key
}

func ValidRequired() Required {
	return Required{}
}

func (r Required) IsSatisfied(obj interface{}) bool {
	if obj == nil {
		return false
	}

	if str, ok := obj.(string); ok {
		return len(str) > 0
	}
	if b, ok := obj.(bool); ok {
		return b
	}
	if i, ok := obj.(int); ok {
		return i != 0
	}
	if t, ok := obj.(time.Time); ok {
		return !t.IsZero()
	}
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Slice {
		return v.Len() > 0
	}
	return true
}

func (r Required) DefaultMessage() string {
	return "Required"
}

type Min struct {
	Key string
	Min int
}

func (v Min) GetKey() string {
	return v.Key
}

func ValidMin(min int) Min {
	return Min{Min: min}
}

func (m Min) IsSatisfied(obj interface{}) bool {
	num, ok := obj.(int)
	if ok {
		return num >= m.Min
	}
	return false
}

func (m Min) DefaultMessage() string {
	return fmt.Sprintln("Minimum is", m.Min)
}

type Max struct {
	Key string
	Max int
}

func (v Max) GetKey() string {
	return v.Key
}

func ValidMax(max int) Max {
	return Max{Max: max}
}

func (m Max) IsSatisfied(obj interface{}) bool {
	num, ok := obj.(int)
	if ok {
		return num <= m.Max
	}
	return false
}

func (m Max) DefaultMessage() string {
	return fmt.Sprintln("Maximum is", m.Max)
}

// Requires an integer to be within Min, Max inclusive.
type Range struct {
	Key string
	Min
	Max
}

func (v Range) GetKey() string {
	return v.Key
}

func ValidRange(min, max int) Range {
	return Range{Min: Min{Min: min}, Max: Max{Max: max}}
}

func (r Range) IsSatisfied(obj interface{}) bool {
	return r.Min.IsSatisfied(obj) && r.Max.IsSatisfied(obj)
}

func (r Range) DefaultMessage() string {
	return fmt.Sprintln("Range is", r.Min.Min, "to", r.Max.Max)
}

// Requires an array or string to be at least a given length.
type MinSize struct {
	Key string
	Min int
}

func (v MinSize) GetKey() string {
	return v.Key
}

func ValidMinSize(min int) MinSize {
	return MinSize{Min: min}
}

func (m MinSize) IsSatisfied(obj interface{}) bool {
	if str, ok := obj.(string); ok {
		return len(str) >= m.Min
	}
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Slice {
		return v.Len() >= m.Min
	}
	return false
}

func (m MinSize) DefaultMessage() string {
	return fmt.Sprintln("Minimum size is", m.Min)
}

// Requires an array or string to be at most a given length.
type MaxSize struct {
	Key string
	Max int
}

func (v MaxSize) GetKey() string {
	return v.Key
}

func ValidMaxSize(max int) MaxSize {
	return MaxSize{Max: max}
}

func (m MaxSize) IsSatisfied(obj interface{}) bool {
	if str, ok := obj.(string); ok {
		return len(str) <= m.Max
	}
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Slice {
		return v.Len() <= m.Max
	}
	return false
}

func (m MaxSize) DefaultMessage() string {
	return fmt.Sprintln("Maximum size is", m.Max)
}

// Requires an array or string to be exactly a given length.
type Length struct {
	Key string
	N   int
}

func (v Length) GetKey() string {
	return v.Key
}

func ValidLength(n int) Length {
	return Length{N: n}
}

func (s Length) IsSatisfied(obj interface{}) bool {
	if str, ok := obj.(string); ok {
		return len(str) == s.N
	}
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Slice {
		return v.Len() == s.N
	}
	return false
}

func (s Length) DefaultMessage() string {
	return fmt.Sprintln("Required length is", s.N)
}

// Requires a string to match a given regex.
type Match struct {
	Key    string
	Regexp *regexp.Regexp
}

func (v Match) GetKey() string {
	return v.Key
}

func ValidMatch(regex *regexp.Regexp) Match {
	return Match{Regexp: regex}
}

func (m Match) IsSatisfied(obj interface{}) bool {
	str := obj.(string)
	return m.Regexp.MatchString(str)
}

func (m Match) DefaultMessage() string {
	return fmt.Sprintln("Must match", m.Regexp)
}

var emailPattern = regexp.MustCompile("^[\\w!#$%&'*+/=?^_`{|}~-]+(?:\\.[\\w!#$%&'*+/=?^_`{|}~-]+)*@(?:[\\w](?:[\\w-]*[\\w])?\\.)+[a-zA-Z0-9](?:[\\w-]*[\\w])?$")

type Email struct {
	Key string
	Match
}

func (v Email) GetKey() string {
	return v.Key
}

func ValidEmail() Email {
	return Email{Match: Match{Regexp: emailPattern}}
}

func (e Email) DefaultMessage() string {
	return fmt.Sprintln("Must be a valid email address")
}
