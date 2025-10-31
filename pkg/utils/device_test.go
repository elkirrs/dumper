package utils_test

import (
	"dumper/pkg/utils"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseIoRegUUID(t *testing.T) {
	out := `Some info
	"IOPlatformUUID" = "ABCDEF-12345-67890"
	Other info`
	uuid := utils.ParseIoRegUUID(out)
	assert.Equal(t, "abcdef-12345-67890", uuid)

	assert.Empty(t, utils.ParseIoRegUUID("no uuid here"))
}

func TestParseWmicUUID(t *testing.T) {
	out := `UUID
	ABCDEF-12345-67890`
	uuid := utils.ParseWmicUUID(out)
	assert.Equal(t, "abcdef-12345-67890", uuid)

	assert.Empty(t, utils.ParseWmicUUID("uuid\n"))
}

func TestCollectMACs(t *testing.T) {
	macs := utils.CollectMACs()
	assert.IsType(t, "", macs)
}

func TestGetDeviceID_Override(t *testing.T) {
	os.Setenv("DEVICE_ID_OVERRIDE", "override123")
	defer os.Unsetenv("DEVICE_ID_OVERRIDE")

	id := utils.GetDeviceID()
	assert.Equal(t, "override123", id)
}

func TestGetDeviceKey(t *testing.T) {
	os.Setenv("DEVICE_ID_OVERRIDE", "my-device")
	defer os.Unsetenv("DEVICE_ID_OVERRIDE")

	key := utils.GetDeviceKey()
	assert.Equal(t, 32, len(key)) // SHA256 always 32 bytes
}

func TestTryReadFile(t *testing.T) {
	tmpFile := t.TempDir() + "/file.txt"
	os.WriteFile(tmpFile, []byte(" Hello \n"), 0644)

	val := utils.TryReadFile(tmpFile)
	assert.Equal(t, "hello", val)

	assert.Empty(t, utils.TryReadFile("/nonexistent/file"))
}

func TestTryRunAndTrim(t *testing.T) {
	out := utils.TryRunAndTrim("echo", "hello")
	assert.Equal(t, "hello", out)

	out2 := utils.TryRunAndTrim("nonexistent-command")
	assert.Empty(t, out2)
}
