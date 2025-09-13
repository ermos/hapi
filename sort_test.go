package hapi

import (
	"reflect"
	"testing"
)

func TestParseSortFromString(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    Sort
		wantErr bool
	}{
		{
			name:  "Valid ascending sort",
			value: "name:asc",
			want: Sort{
				Field:     "name",
				Direction: SortDirectionAsc,
			},
			wantErr: false,
		},
		{
			name:  "Valid descending sort",
			value: "age:desc",
			want: Sort{
				Field:     "age",
				Direction: SortDirectionDesc,
			},
			wantErr: false,
		},
		{
			name:  "Field with underscores",
			value: "created_at:asc",
			want: Sort{
				Field:     "created_at",
				Direction: SortDirectionAsc,
			},
			wantErr: false,
		},
		{
			name:  "Field with dots",
			value: "user.name:desc",
			want: Sort{
				Field:     "user.name",
				Direction: SortDirectionDesc,
			},
			wantErr: false,
		},
		{
			name:    "Invalid format - no colon",
			value:   "name",
			want:    Sort{},
			wantErr: true,
		},
		{
			name:    "Invalid format - multiple colons",
			value:   "name:asc:extra",
			want:    Sort{},
			wantErr: true,
		},
		{
			name:    "Invalid direction",
			value:   "name:invalid",
			want:    Sort{},
			wantErr: true,
		},
		{
			name:  "Empty field",
			value: ":asc",
			want: Sort{
				Field:     "",
				Direction: SortDirectionAsc,
			},
			wantErr: false,
		},
		{
			name:    "Empty direction",
			value:   "name:",
			want:    Sort{},
			wantErr: true,
		},
		{
			name:    "Empty string",
			value:   "",
			want:    Sort{},
			wantErr: true,
		},
		{
			name:    "Only colon",
			value:   ":",
			want:    Sort{},
			wantErr: true,
		},
		{
			name:    "Case sensitive direction",
			value:   "name:ASC",
			want:    Sort{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseSortFromString(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseSortFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseSortFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}
