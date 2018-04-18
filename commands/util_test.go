package commands

import (
	"io/ioutil"
	"os"
	"runtime"
	"testing"
)

func TestGetLinuxDist(t *testing.T) {
	tests := map[string]struct {
		content      []byte
		writeTmpFile bool
		expected     string
	}{
		"os-release does not exist": {
			writeTmpFile: false,
			expected:     "",
		},
		"no ID defined": {
			content: []byte(`
NAME="Arch Linux"
PRETTY_NAME="Arch Linux"
ID_LIKE=archlinux
ANSI_COLOR="0;36"
HOME_URL="https://www.archlinux.org/"
SUPPORT_URL="https://bbs.archlinux.org/"
BUG_REPORT_URL="https://bugs.archlinux.org/"
`),
			writeTmpFile: true,
			expected:     "",
		},
		"Fedora example": {
			content: []byte(`
NAME=Fedora
VERSION="17 (Beefy Miracle)"
ID=fedora
VERSION_ID=17
PRETTY_NAME="Fedora 17 (Beefy Miracle)"
ANSI_COLOR="0;34"
CPE_NAME="cpe:/o:fedoraproject:fedora:17"
HOME_URL="https://fedoraproject.org/"
BUG_REPORT_URL="https://bugzilla.redhat.com/"
`),
			writeTmpFile: true,
			expected:     "fedora",
		},
		"Ubuntu example": {
			content: []byte(`
NAME="Ubuntu"
VERSION="13.10, Saucy Salamander"
ID=ubuntu
ID_LIKE=debian
PRETTY_NAME="Ubuntu 13.10"
VERSION_ID="13.10"
HOME_URL="http://www.ubuntu.com/"
SUPPORT_URL="http://help.ubuntu.com/"
BUG_REPORT_URL="http://bugs.launchpad.net/ubuntu/"
`),
			writeTmpFile: true,
			expected:     "ubuntu",
		},
		"Arch Linux example": {
			content: []byte(`
NAME="Arch Linux"
PRETTY_NAME="Arch Linux"
ID=arch
ID_LIKE=archlinux
ANSI_COLOR="0;36"
HOME_URL="https://www.archlinux.org/"
SUPPORT_URL="https://bbs.archlinux.org/"
BUG_REPORT_URL="https://bugs.archlinux.org/"
`),
			writeTmpFile: true,
			expected:     "arch",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			tmp := linuxOSReleasePath

			if test.writeTmpFile {
				f, err := ioutil.TempFile(os.TempDir(), "linuxdist")
				if err != nil {
					t.Fatalf("Expected nil, got %s", err.Error())
				}

				f.Write(test.content)
				defer os.Remove(f.Name())
				defer f.Close()

				linuxOSReleasePath = f.Name()
			} else {
				linuxOSReleasePath = "/does/not/exist"
			}
			defer func() { linuxOSReleasePath = tmp }()

			dist := getLinuxDist()
			if dist != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, dist)
			}
		})
	}
}

func TestLinuxCommandExists(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skipf("skipping test for platform %s", runtime.GOOS)
	}

	tests := map[string]struct {
		cmd      string
		expected bool
	}{
		"command exists":         {cmd: "true", expected: true},
		"command does not exist": {cmd: "does_not_exist", expected: false},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			exists := linuxCommandExists(test.cmd)
			if exists != test.expected {
				t.Errorf("Expected %t, got %t", test.expected, exists)
			}
		})
	}
}
