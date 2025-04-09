package deeplx_translator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTextToString(t *testing.T) {
	tests := []struct {
		text     any
		expected string
		wantErr  bool
	}{
		{"hello world", "hello world", false},
		{[]string{"line1", "line2", "line3"}, "line1\nline2\nline3", false},
		{123, "", true},
	}
	for _, tt := range tests {
		result, err := textToString(tt.text)
		if tt.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		}
	}
}

func TestTextToStringSlice(t *testing.T) {
	tests := []struct {
		text     any
		expected []string
		wantErr  bool
	}{
		{
			text:     "hello, world!",
			expected: []string{"hello, world!"},
			wantErr:  false,
		},
		{
			text:     "你好。今天天气不错！你要去哪里？再见……",
			expected: []string{"你好。", "今天天气不错！", "你要去哪里？", "再见……"},
			wantErr:  false,
		},
		{
			text:     []string{"foo", "bar"},
			expected: []string{"foo", "bar"},
			wantErr:  false,
		},
		{
			text:     3.14,
			expected: nil,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		result, err := textToStringSlice(tt.text)
		if tt.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		}
	}
}

func TestSplitTextsAfter(t *testing.T) {
	tests := []struct {
		input    string
		seps     []string
		expected []string
	}{
		{
			input:    "Hello. World! How are you?",
			seps:     []string{".", "!", "?"},
			expected: []string{"Hello.", " World!", " How are you?"},
		},
		{
			input:    "Hello. World! I'm ok; How are you?",
			seps:     defaultSentenceTerminators,
			expected: []string{"Hello.", " World!", " I'm ok;", " How are you?"},
		},
		{
			input:    "你好。欢迎你！要喝水吗？再见……",
			seps:     defaultSentenceTerminators,
			expected: []string{"你好。", "欢迎你！", "要喝水吗？", "再见……"},
		},
	}
	for _, tt := range tests {
		result := splitTextsAfter(tt.input, tt.seps...)
		assert.Equal(t, tt.expected, result)
	}
}
