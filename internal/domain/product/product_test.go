package product

import "testing"

func TestNewProduct(t *testing.T) {
	p := NewProduct("name", "description", 100.5, 5)

	if p.Name != "name" {
		t.Errorf("expected Name 'Тест', got '%s'", p.Name)
	}
	if p.Description != "description" {
		t.Errorf("expected Description 'Описание', got '%s'", p.Description)
	}
	if p.Price != 100.5 {
		t.Errorf("expected Price 100, got %f", p.Price)
	}
	if p.CategoryId != 5 {
		t.Errorf("expected CategoryId 1, got %d", p.CategoryId)
	}
}
