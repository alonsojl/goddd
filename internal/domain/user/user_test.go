package user_test

import (
	"goddd/internal/domain/user"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/stretchr/testify/assert"
)

func TestUserName(t *testing.T) {
	table := []struct {
		title     string
		wantError bool
		err       validation.Error
		name      string
	}{
		{
			title:     "cannot be blank",
			wantError: true,
			err:       validation.ErrRequired,
			name:      "",
		},
		{
			title:     "the length must be between 1 and 10",
			wantError: true,
			err:       validation.ErrLengthOutOfRange,
			name:      "Jorge Luis Alonso Hdez",
		},
		{
			title:     "success",
			wantError: false,
			err:       nil,
			name:      "Jorge Luis",
		},
	}
	var verr validation.Error
	for _, tt := range table {
		t.Run(tt.title, func(t *testing.T) {
			err := validation.Validate(tt.name, user.NameRules()...)
			if !tt.wantError {
				assert.Equal(t, tt.err, err)
			} else {
				assert.ErrorAs(t, err, &verr)
				assert.Equal(t, tt.err.Code(), verr.Code())
			}
		})
	}
}
