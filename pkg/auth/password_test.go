package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	t.Run("Successfully hash password", func(t *testing.T) {
		// 測試密碼
		password := "testPassword123"

		// 加密密碼
		hashedPassword, err := HashPassword(password)
		require.NoError(t, err)
		assert.NotEmpty(t, hashedPassword)
		assert.NotEqual(t, password, hashedPassword)

		// 驗證加密後的密碼
		err = CheckPassword(password, hashedPassword)
		assert.NoError(t, err)
	})

	t.Run("Hash different passwords", func(t *testing.T) {
		// 測試兩個不同的密碼
		password1 := "password1"
		password2 := "password2"

		// 加密兩個密碼
		hashedPassword1, err := HashPassword(password1)
		require.NoError(t, err)
		hashedPassword2, err := HashPassword(password2)
		require.NoError(t, err)

		// 驗證加密後的密碼不同
		assert.NotEqual(t, hashedPassword1, hashedPassword2)

		// 驗證每個密碼只能匹配自己的雜湊值
		err = CheckPassword(password1, hashedPassword1)
		assert.NoError(t, err)
		err = CheckPassword(password2, hashedPassword2)
		assert.NoError(t, err)
		err = CheckPassword(password1, hashedPassword2)
		assert.Error(t, err)
	})

	t.Run("Hash same password multiple times", func(t *testing.T) {
		// 測試密碼
		password := "testPassword123"

		// 多次加密同一個密碼
		hashedPassword1, err := HashPassword(password)
		require.NoError(t, err)
		hashedPassword2, err := HashPassword(password)
		require.NoError(t, err)

		// 驗證每次加密的結果都不同（因為 salt 不同）
		assert.NotEqual(t, hashedPassword1, hashedPassword2)

		// 驗證兩個雜湊值都能正確驗證原始密碼
		err = CheckPassword(password, hashedPassword1)
		assert.NoError(t, err)
		err = CheckPassword(password, hashedPassword2)
		assert.NoError(t, err)
	})

	t.Run("Hash password exceeds bcrypt limit", func(t *testing.T) {
		// 創建一個超過 bcrypt 限制的密碼（100 字節）
		longPassword := make([]byte, 100)
		for i := range longPassword {
			longPassword[i] = 'a'
		}

		// 加密長密碼應該失敗
		hashedPassword, err := HashPassword(string(longPassword))
		assert.Error(t, err)
		assert.Empty(t, hashedPassword)
		assert.Contains(t, err.Error(), "password length exceeds 72 bytes")
	})

	t.Run("Hash password with special characters", func(t *testing.T) {
		// 測試包含特殊字符的密碼
		specialPassword := "!@#$%^&*()_+-=[]{}|;:,.<>?/~`"

		// 加密特殊字符密碼
		hashedPassword, err := HashPassword(specialPassword)
		require.NoError(t, err)
		assert.NotEmpty(t, hashedPassword)

		// 驗證特殊字符密碼
		err = CheckPassword(specialPassword, hashedPassword)
		assert.NoError(t, err)
	})

	t.Run("Hash password with unicode characters", func(t *testing.T) {
		// 測試包含 Unicode 字符的密碼
		unicodePassword := "密碼123!@#$%^&*()_+你好世界"

		// 加密 Unicode 密碼
		hashedPassword, err := HashPassword(unicodePassword)
		require.NoError(t, err)
		assert.NotEmpty(t, hashedPassword)

		// 驗證 Unicode 密碼
		err = CheckPassword(unicodePassword, hashedPassword)
		assert.NoError(t, err)
	})
}

func TestCheckPassword(t *testing.T) {
	t.Run("Correct password", func(t *testing.T) {
		// 測試密碼
		password := "testPassword123"

		// 加密密碼
		hashedPassword, err := HashPassword(password)
		require.NoError(t, err)

		// 驗證正確的密碼
		err = CheckPassword(password, hashedPassword)
		assert.NoError(t, err)
	})

	t.Run("Incorrect password", func(t *testing.T) {
		// 測試密碼
		password := "testPassword123"
		wrongPassword := "wrongPassword123"

		// 加密密碼
		hashedPassword, err := HashPassword(password)
		require.NoError(t, err)

		// 驗證錯誤的密碼
		err = CheckPassword(wrongPassword, hashedPassword)
		assert.Error(t, err)
	})

	t.Run("Empty password", func(t *testing.T) {
		// 測試空密碼
		password := ""
		wrongPassword := "wrongPassword123"

		// 加密空密碼
		hashedPassword, err := HashPassword(password)
		require.NoError(t, err)

		// 驗證空密碼
		err = CheckPassword(password, hashedPassword)
		assert.NoError(t, err)

		// 驗證錯誤的密碼
		err = CheckPassword(wrongPassword, hashedPassword)
		assert.Error(t, err)
	})

	t.Run("Invalid hash format", func(t *testing.T) {
		// 測試無效的雜湊格式
		password := "testPassword123"
		invalidHash := "invalid-hash-format"

		// 驗證無效的雜湊格式
		err := CheckPassword(password, invalidHash)
		assert.Error(t, err)
	})

	t.Run("Check password with different case", func(t *testing.T) {
		// 測試密碼
		password := "TestPassword123"
		passwordLower := "testpassword123"

		// 加密密碼
		hashedPassword, err := HashPassword(password)
		require.NoError(t, err)

		// 驗證不同大小寫的密碼
		err = CheckPassword(passwordLower, hashedPassword)
		assert.Error(t, err)
	})

	t.Run("Check password with extra spaces", func(t *testing.T) {
		// 測試密碼
		password := "testPassword123"
		passwordWithSpaces := " testPassword123 "

		// 加密密碼
		hashedPassword, err := HashPassword(password)
		require.NoError(t, err)

		// 驗證帶有額外空格的密碼
		err = CheckPassword(passwordWithSpaces, hashedPassword)
		assert.Error(t, err)
	})
}
