package device

import (
	"bytes"
	"crypto/sha256"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func ParseIoRegUUID(ioRegOut string) string {
	for _, line := range strings.Split(ioRegOut, "\n") {
		if strings.Contains(line, "IOPlatformUUID") {
			parts := strings.Split(line, "=")
			if len(parts) >= 2 {
				val := strings.TrimSpace(parts[1])
				val = strings.Trim(val, "\"")
				return strings.ToLower(val)
			}
		}
	}
	return ""
}

func ParseWmicUUID(wmicOut string) string {
	for _, line := range strings.Split(wmicOut, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.EqualFold(line, "uuid") {
			continue
		}
		return strings.ToLower(line)
	}
	return ""
}

func CollectMACs() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	var buf bytes.Buffer
	for _, ifc := range ifaces {
		m := ifc.HardwareAddr.String()
		if m == "" || m == "00:00:00:00:00:00" {
			continue
		}
		buf.WriteString(strings.ToLower(strings.ReplaceAll(m, ":", "")))
	}
	return buf.String()
}

func GetDeviceID() string {
	if v := strings.TrimSpace(os.Getenv("DEVICE_ID_OVERRIDE")); v != "" {
		return v
	}

	switch runtime.GOOS {
	case "linux":
		if id := TryReadFile("/etc/machine-id"); id != "" {
			return id
		}
		if id := TryReadFile("/var/lib/dbus/machine-id"); id != "" {
			return id
		}
		if id := TryReadFile("/sys/class/dmi/id/product_uuid"); id != "" {
			return id
		}
	case "darwin":
		if id := TryRunAndTrim("ioreg", "-rd1", "-c", "IOPlatformExpertDevice"); id != "" {
			if parsed := ParseIoRegUUID(id); parsed != "" {
				return parsed
			}
		}
	case "windows":
		if id := TryRunAndTrim("wmic", "csproduct", "get", "uuid"); id != "" {
			if parsed := ParseWmicUUID(id); parsed != "" {
				return parsed
			}
		}
		if id := TryRunAndTrim("powershell", "-Command", "Get-WmiObject -Class Win32_ComputerSystemProduct | Select-Object -ExpandProperty UUID"); id != "" {
			return id
		}
	}

	if macs := CollectMACs(); macs != "" {
		return macs
	}

	return ""
}

func GetDeviceKey() []byte {
	id := GetDeviceID()
	sum := sha256.Sum256([]byte(id))
	return sum[:]
}

func TryReadFile(path string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.ToLower(strings.TrimSpace(string(b)))
}

func TryRunAndTrim(name string, args ...string) string {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
