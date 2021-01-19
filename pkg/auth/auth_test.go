package auth

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_PBKDF2HashPassword(t *testing.T) {
	password := "password"
	salt := "TESTSALT"
	iterations := 150000
	keyLenght := 64
	h := pbkdf2HashPassword(password, salt, iterations, keyLenght)
	t.Logf("Result hash of password='%s' and salt='%s' is: %s", password, salt, h)
	require.NotEmpty(t, h, "expected to see some string")
	require.Equal(t,
		"7fc909bc1888eb3b5c717dfcca83a2b1b031ecb40ed6ad4278399d78d29ea0212805336b86df2c7254e9d206eee53b5300edeaeee6f35bb96b8f0890c693d24f",
		h,
		"different hashed values",
	)
}

func Test_PBKDF2HashPassword_Special_Cases(t *testing.T) {
	password := "password"
	salt := "TESTSALT"
	iterations := 1
	keyLenght := 1
	h := pbkdf2HashPassword(password, salt, iterations, keyLenght)
	t.Logf("Result hash of password='%s' and salt='%s' is: %s", password, salt, h)
}

func Test_GenerateRandomString(t *testing.T) {
	lenght := 5
	str1, err := GenerateRandomString(lenght)
	require.NoError(t, err, "generate random string must not fail")
	require.Equal(t, lenght, len(str1), "resulted string length differs from expected")

	str2, err := GenerateRandomString(lenght)
	require.NoError(t, err, "generate random string must not fail")
	require.Equal(t, lenght, len(str2), "resulted string length differs from expected")

	require.NotEqual(t, str1, str2, "generated values must be different")
}
