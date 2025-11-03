package config

import (
	"math"
	"os"
	"testing"
)

func TestParseConfig(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		fileName    string
		setupFile   func() string
		cleanupFile func(string)
		wantErr     bool
		validate    func(*dotNotationMap) bool
	}{
		{
			name: "valid config file",
			content: `[log]
level=info

[device]
id=test-device-123

[mqtt]
broker=tcp://localhost:1883
qos=1`,
			setupFile: func() string {
				tmpfile, err := os.CreateTemp("", "test_config_*.toml")
				if err != nil {
					t.Fatal(err)
				}
				return tmpfile.Name()
			},
			cleanupFile: func(name string) {
				os.Remove(name)
			},
			wantErr: false,
			validate: func(m *dotNotationMap) bool {
				return m.Get("log.level") == "info" &&
					m.Get("device.id") == "test-device-123" &&
					m.Get("mqtt.broker") == "tcp://localhost:1883"
			},
		},
		{
			name:        "non-existent file",
			fileName:    "non_existent.toml",
			wantErr:     true,
			validate:    func(*dotNotationMap) bool { return true },
			setupFile:   func() string { return "" },
			cleanupFile: func(string) {},
		},
		{
			name: "multiline string",
			content: `[device]
description="""This is a
multiline
description"""`,
			setupFile: func() string {
				tmpfile, err := os.CreateTemp("", "test_config_*.toml")
				if err != nil {
					t.Fatal(err)
				}
				return tmpfile.Name()
			},
			cleanupFile: func(name string) {
				os.Remove(name)
			},
			wantErr: false,
			validate: func(m *dotNotationMap) bool {
				val := m.Get("device.description")
				str, ok := val.(string)
				return ok && str == "This is a\nmultiline\ndescription"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fileName string
			if tt.setupFile != nil {
				fileName = tt.setupFile()
				if fileName != "" && tt.content != "" {
					if err := os.WriteFile(fileName, []byte(tt.content), 0644); err != nil {
						t.Fatal(err)
					}
				}
			} else {
				fileName = tt.fileName
			}

			defer tt.cleanupFile(fileName)

			got, err := parseConfig(fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !tt.validate(got) {
				t.Errorf("parseConfig() validation failed")
			}
		})
	}
}

func TestRemoveComments(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "no comment",
			input: "key=value",
			want:  "key=value",
		},
		{
			name:  "with comment",
			input: "key=value # this is a comment",
			want:  "key=value",
		},
		{
			name:  "only comment",
			input: "# this is a comment",
			want:  "",
		},
		{
			name:  "comment with spaces",
			input: "  key=value  # comment",
			want:  "key=value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeComments(tt.input); got != tt.want {
				t.Errorf("removeComments() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseValue(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  interface{}
	}{
		{
			name:  "boolean true",
			input: "true",
			want:  true,
		},
		{
			name:  "boolean false",
			input: "false",
			want:  false,
		},
		{
			name:  "null",
			input: "null",
			want:  nil,
		},
		{
			name:  "empty string",
			input: "",
			want:  nil,
		},
		{
			name:  "integer",
			input: "42",
			want:  42,
		},
		{
			name:  "float",
			input: "3.14",
			want:  3.14,
		},
		{
			name:  "hexadecimal",
			input: "0xFF",
			want:  int64(255),
		},
		{
			name:  "octal",
			input: "0o755",
			want:  int64(493),
		},
		{
			name:  "binary",
			input: "0b1010",
			want:  int64(10),
		},
		{
			name:  "quoted string",
			input: `"hello world"`,
			want:  "hello world",
		},
		{
			name:  "single quoted string",
			input: `'hello world'`,
			want:  "hello world",
		},
		{
			name:  "plus nan",
			input: "+nan",
			want:  math.NaN(),
		},
		{
			name:  "minus nan",
			input: "-nan",
			want:  -math.NaN(),
		},
		{
			name:  "plain string",
			input: "hello",
			want:  "hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseValue(tt.input)
			if tt.name == "plus nan" || tt.name == "minus nan" {
				if gotFloat, ok := got.(float64); ok {
					if !math.IsNaN(gotFloat) && !math.IsNaN(-gotFloat) {
						t.Errorf("parseValue() = %v, want NaN", got)
					}
				} else {
					t.Errorf("parseValue() = %v, want NaN", got)
				}
			} else if got != tt.want {
				t.Errorf("parseValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJoinKeys(t *testing.T) {
	tests := []struct {
		name string
		keys []string
		want string
	}{
		{
			name: "single key",
			keys: []string{"key"},
			want: "key",
		},
		{
			name: "multiple keys",
			keys: []string{"mqtt", "topics", "data"},
			want: "mqtt.topics.data",
		},
		{
			name: "keys with empty strings",
			keys: []string{"mqtt", "", "topics", "  ", "data"},
			want: "mqtt.topics.data",
		},
		{
			name: "empty keys",
			keys: []string{"", "  "},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := JoinKeys(tt.keys...); got != tt.want {
				t.Errorf("JoinKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnquoteMaybe(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "double quoted string",
			input:   `"hello"`,
			want:    "hello",
			wantErr: false,
		},
		{
			name:    "single quoted string",
			input:   `'world'`,
			want:    "world",
			wantErr: false,
		},
		{
			name:    "not quoted",
			input:   "hello",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := unquoteMaybe(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("unquoteMaybe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("unquoteMaybe() = %v, want %v", got, tt.want)
			}
		})
	}
}
