package store

import (
	"bytes"
	"os/user"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/lima-vm/lima/pkg/limayaml"
	"gotest.tools/v3/assert"
)

const separator = string(filepath.Separator)

var vmtype = limayaml.QEMU
var goarch = limayaml.NewArch(runtime.GOARCH)

var instance = Instance{
	Name:   "foo",
	Status: StatusStopped,
	VMType: vmtype,
	Arch:   goarch,
	Dir:    "dir",
}

var table = "NAME    STATUS     SSH            CPUS    MEMORY    DISK    DIR\n" +
	"foo     Stopped    127.0.0.1:0    0       0B        0B      dir\n"

var tableEmu = "NAME    STATUS     SSH            ARCH       CPUS    MEMORY    DISK    DIR\n" +
	"foo     Stopped    127.0.0.1:0    unknown    0       0B        0B      dir\n"

var tableHome = "NAME    STATUS     SSH            CPUS    MEMORY    DISK    DIR\n" +
	"foo     Stopped    127.0.0.1:0    0       0B        0B      ~" + separator + "dir\n"

var tableAll = "NAME    STATUS     SSH            VMTYPE    ARCH      CPUS    MEMORY    DISK    DIR\n" +
	"foo     Stopped    127.0.0.1:0    " + vmtype + "      " + goarch + "    0       0B        0B      dir\n"

// for width 60, everything is hidden
var table60 = "NAME    STATUS     SSH            CPUS    MEMORY    DISK\n" +
	"foo     Stopped    127.0.0.1:0    0       0B        0B\n"

// for width 80, identical is hidden (type/arch)
var table80i = "NAME    STATUS     SSH            CPUS    MEMORY    DISK    DIR\n" +
	"foo     Stopped    127.0.0.1:0    0       0B        0B      dir\n"

// for width 80, different arch is still shown (not dir)
var table80d = "NAME    STATUS     SSH            ARCH       CPUS    MEMORY    DISK\n" +
	"foo     Stopped    127.0.0.1:0    unknown    0       0B        0B\n"

// for width 100, nothing is hidden
var table100 = "NAME    STATUS     SSH            VMTYPE    ARCH      CPUS    MEMORY    DISK    DIR\n" +
	"foo     Stopped    127.0.0.1:0    " + vmtype + "      " + goarch + "    0       0B        0B      dir\n"

// for width 80, directory is hidden (if not identical)
var tableTwo = "NAME    STATUS     SSH            VMTYPE    ARCH       CPUS    MEMORY    DISK\n" +
	"foo     Stopped    127.0.0.1:0    qemu      x86_64     0       0B        0B\n" +
	"bar     Stopped    127.0.0.1:0    vz        aarch64    0       0B        0B\n"

func TestPrintInstanceTable(t *testing.T) {
	var buf bytes.Buffer
	instances := []*Instance{&instance}
	PrintInstances(&buf, instances, "table", nil)
	assert.Equal(t, table, buf.String())
}

func TestPrintInstanceTableEmu(t *testing.T) {
	var buf bytes.Buffer
	instance1 := instance
	instance1.Arch = "unknown"
	instances := []*Instance{&instance1}
	PrintInstances(&buf, instances, "table", nil)
	assert.Equal(t, tableEmu, buf.String())
}

func TestPrintInstanceTableHome(t *testing.T) {
	var buf bytes.Buffer
	u, err := user.Current()
	assert.NilError(t, err)
	instance1 := instance
	instance1.Dir = filepath.Join(u.HomeDir, "dir")
	instances := []*Instance{&instance1}
	PrintInstances(&buf, instances, "table", nil)
	assert.Equal(t, tableHome, buf.String())
}

func TestPrintInstanceTable60(t *testing.T) {
	var buf bytes.Buffer
	instances := []*Instance{&instance}
	options := PrintOptions{TerminalWidth: 60}
	PrintInstances(&buf, instances, "table", &options)
	assert.Equal(t, table60, buf.String())
}

func TestPrintInstanceTable80SameArch(t *testing.T) {
	var buf bytes.Buffer
	instances := []*Instance{&instance}
	options := PrintOptions{TerminalWidth: 80}
	PrintInstances(&buf, instances, "table", &options)
	assert.Equal(t, table80i, buf.String())
}

func TestPrintInstanceTable80DiffArch(t *testing.T) {
	var buf bytes.Buffer
	instance1 := instance
	instance1.Arch = limayaml.NewArch("unknown")
	instances := []*Instance{&instance1}
	options := PrintOptions{TerminalWidth: 80}
	PrintInstances(&buf, instances, "table", &options)
	assert.Equal(t, table80d, buf.String())
}

func TestPrintInstanceTable100(t *testing.T) {
	var buf bytes.Buffer
	instances := []*Instance{&instance}
	options := PrintOptions{TerminalWidth: 100}
	PrintInstances(&buf, instances, "table", &options)
	assert.Equal(t, table100, buf.String())
}

func TestPrintInstanceTableAll(t *testing.T) {
	var buf bytes.Buffer
	instances := []*Instance{&instance}
	options := PrintOptions{TerminalWidth: 40, AllFields: true}
	PrintInstances(&buf, instances, "table", &options)
	assert.Equal(t, tableAll, buf.String())
}

func TestPrintInstanceTableTwo(t *testing.T) {
	var buf bytes.Buffer
	instance1 := instance
	instance1.Name = "foo"
	instance1.VMType = limayaml.QEMU
	instance1.Arch = limayaml.X8664
	instance2 := instance
	instance2.Name = "bar"
	instance2.VMType = limayaml.VZ
	instance2.Arch = limayaml.AARCH64
	instances := []*Instance{&instance1, &instance2}
	options := PrintOptions{TerminalWidth: 80}
	PrintInstances(&buf, instances, "table", &options)
	assert.Equal(t, tableTwo, buf.String())
}