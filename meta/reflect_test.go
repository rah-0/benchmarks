package meta

import (
	"reflect"
	"testing"
)

func BenchmarkMethodAccess(b *testing.B) {
	p := &Person{}

	b.Run("SetName", func(b *testing.B) {
		for b.Loop() {
			p.SetName("John")
		}
	})

	b.Run("GetName", func(b *testing.B) {
		p.SetName("John")
		for b.Loop() {
			_ = p.GetName()
		}
	})
}

func BenchmarkReflectionAccess(b *testing.B) {
	p := &Person{}
	val := reflect.ValueOf(p).Elem()
	nameField := val.FieldByName("Name")

	b.Run("SetName", func(b *testing.B) {
		for b.Loop() {
			nameField.SetString("John")
		}
	})

	b.Run("GetName", func(b *testing.B) {
		nameField.SetString("John")
		for b.Loop() {
			_ = nameField.String()
		}
	})
}
