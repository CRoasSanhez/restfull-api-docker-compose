package auth

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var headerSigned = "Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ5YWxvY2hhdF9hdXRoIiwiYm90IjoiYWVyb21leGljby0xMjM0NSIsImlhdCI6MTU4MzUxNDQ5MCwiaXNzIjoieWFsb2NoYXRfYXV0aCIsImp0aSI6Ijk1YWFjMWI0LTAyZjYtNGQ3NC05NTIyLWQ4MTVkZWJjODViZSIsIm5iZiI6MTU4MzUxNDQ5MCwic3ViIjoiNWUzMDkzNjQ5ZDgwMmNhOGYzZWU4MjAwIiwidHlwIjoiYWNjZXNzIn0.o04yFUkpjuEUpmmHdf_sLq4vUetAZVwi24wfrT6MQd4mA3S2e2Wz6Cl7ny5zD2SodXWsdAmAJDdxXmqoHMjy3g"

func TestSignToken(t *testing.T) {

	userID := "sample_user_id"
	token, err := SignToken(userID)
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("FAIL validating token: %s", err.Error()))
	}
	assert.NotEmpty(t, token, "Token empty")
}

func TestValidateToken(t *testing.T) {
	CallValidateToken(t, headerSigned)
}

func CallValidateToken(t *testing.T, header string) {

	var claims = &Claims{}
	err := ValidateToken(header, claims)
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("FAIL validating token: %s", err.Error()))
	}

	assert.True(t, true, "Valid token")
}

func TestRenewToken(t *testing.T) {
	var claims = &Claims{}
	var header = headerSigned

	token, err := RenewToken(header, claims)
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("Error renewing token: %s", err.Error()))
		return
	}

	assert.NotEmptyf(t, token, "Token renewed")
}

// TEST PERFOMANCE
func TestPerformanceSignToken(t *testing.T) {
	bot := "sample_botID_chat"
	subject := "resource_ID"
	atts := 5
	for idx := 0; idx < atts; idx++ {
		go func(index int) {
			_, err := SignToken(bot, subject)
			if err != nil {
				assert.FailNow(t, fmt.Sprintf("Failing singing token %s at idx= %d from %d", err.Error(), idx, atts))
			}
		}(idx)
	}
	time.Sleep(1 * time.Second)
	assert.True(t, true, fmt.Sprintf("Completed %d attempts", atts))
}

func TestPerformanceValidateToken(t *testing.T) {
	atts := 3
	var claims = &Claims{}

	for idx := 0; idx < atts; idx++ {

		go func(t *testing.T, index int) {

			err := ValidateToken(headerSigned, claims)
			if err != nil {
				assert.FailNow(t, "FAIL Validating token", fmt.Sprintf("idx=%d\n", index))
			}
		}(t, idx)
	}

	assert.True(t, true, "Completed")

}
