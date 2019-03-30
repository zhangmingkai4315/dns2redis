package dns2redis

import (
	"testing"

	"github.com/mholt/caddy"
)

func TestSetup(t *testing.T) {
	testCases := []struct {
		input       string
		expectError bool
	}{
		{
			input:       `dns2redis`,
			expectError: false,
		},
		{
			input:       `dns2redis 127.0.0.1`,
			expectError: false,
		},
		{
			input:       `dns2redis 127.0.0.1:6379`,
			expectError: false,
		},
		{
			input:       `dns2redis 6379`,
			expectError: true,
		},
	}
	for _, t := range testCases {
		c := caddy.NewTestController("dns", t.input)
		if err := setup(c)  {
			if err == nil && t.expectError == true{
				t.Fatalln("Expected errors, but got nil")
				continue
			}
			if err !=nil && t.expectError == false{
				t.Fatalf("Expected no errors, but got %s",err)
			}		
		}
	}
}
