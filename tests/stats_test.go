package test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

/////////////////////////////////////////////////
// Tests

func TestStatsGeneration(t *testing.T) {
	args := []string{"stat", "./examples/oauth/secrets-keys.adm"}
	harness := output_interceptor{}
	harness.Hook()

	err := sendToParseArgs(args)
	out, _ := harness.ReadAndRelease()

	assert.Nil(t, err)
	// only the last line contains the required output
	result := strings.Split(out, "\n")[1:]
	assert.Contains(t, result, "Found 1 file(s)")
	assert.Contains(t, result, "MODEL: Client secrets and keys")
	assert.Contains(t, result, "\tATTACK: Secrets extracted via path traversal")
	assert.Contains(t, result, "\tATTACK: Secrets extracted via web-server path traversal")
	assert.Contains(t, result, "\tATTACK: Secrets extracted via encoded path in URL")
	assert.Contains(t, result, "\tATTACK: Secrets extracted via encoded path in URL sent to Apache Web Server")
	assert.Contains(t, result, "\tATTACK: Compromize web server")
	assert.Contains(t, result, "\tATTACK: Secrets extracted from code")
	assert.Contains(t, result, "\tATTACK: Secrets extracted from app binary")
	assert.Contains(t, result, "\tDEFENSE: Block encoded paths")
	assert.Contains(t, result, "\tDEFENSE: Don't allow path traversal to configuration files")
	assert.Contains(t, result, "\tDEFENSE: Web-server shouldn't serve config files")
	assert.Contains(t, result, "\tDEFENSE: Encrypt secrets before storing in file")
	assert.Contains(t, result, "\tDEFENSE: Don't publish code containing secrets")
	assert.Contains(t, result, "\tDEFENSE: Don't publish app binary containing secrets")
	assert.Contains(t, result, "\tDEFENSE: Don't store secrets in config file")
	assert.Contains(t, result, "\tDEFENSE: use Apache Web Server v2.5.51 or later")
	assert.Contains(t, result, "\tDEFENSE: Store secrets and keys in vault")
}

func TestStatsWithPolicies(t *testing.T) {
	args := []string{"stat", "./examples/friends.adm"}
	harness := output_interceptor{}
	harness.Hook()

	err := sendToParseArgs(args)
	out, _ := harness.ReadAndRelease()

	assert.Nil(t, err)
	// only the last line contains the required output
	result := strings.Split(out, "\n")[1:]
	assert.Contains(t, result, "Found 1 file(s)")
	assert.Contains(t, result, "MODEL: Friends fight")
	assert.Contains(t, result, "\tATTACK: Adam cheats Bob")
	assert.Contains(t, result, "\tATTACK: Adam hides the cheat")
	assert.Contains(t, result, "\tDEFENSE: Adam's cheating is caught")
	assert.Contains(t, result, "\tDEFENSE: Bob tries to verify Adam's story")
	assert.Contains(t, result, "\tPOLICY: Honesty is the best policy")
	assert.Contains(t, result, "\t\tDEFENSE: Test honesty")
}

func TestStatsWithPreemtiveDefensesFlag(t *testing.T) {
	args := []string{"stat", "-p", "./examples/friends.adm"}
	harness := output_interceptor{}
	harness.Hook()

	err := sendToParseArgs(args)
	out, _ := harness.ReadAndRelease()

	assert.Nil(t, err)
	// only the last line contains the required output
	result := strings.Split(out, "\n")[1:]
	assert.Contains(t, result, "Found 1 file(s)")
	assert.Contains(t, result, "MODEL: Friends fight")
	assert.Contains(t, result, "\tPOLICY: Honesty is the best policy")
	assert.Contains(t, result, "\t\tDEFENSE: Test honesty")
}

func TestStatsWithIRFlag(t *testing.T) {
	args := []string{"stat", "-i", "./examples/friends.adm"}
	harness := output_interceptor{}
	harness.Hook()

	err := sendToParseArgs(args)
	out, _ := harness.ReadAndRelease()

	assert.Nil(t, err)
	// only the last line contains the required output
	result := strings.Split(out, "\n")[1:]
	assert.Contains(t, result, "Found 1 file(s)")
	assert.Contains(t, result, "MODEL: Friends fight")
	assert.Contains(t, result, "\tDEFENSE: Adam's cheating is caught")
	assert.Contains(t, result, "\tDEFENSE: Bob tries to verify Adam's story")
}

func TestStatsWithDefenseAndIRFlags(t *testing.T) {
	args := []string{"stat", "-d", "-i", "./examples/friends.adm"}
	harness := output_interceptor{}
	harness.Hook()

	err := sendToParseArgs(args)
	out, _ := harness.ReadAndRelease()

	assert.Nil(t, err)
	// only the last line contains the required output
	result := strings.Split(out, "\n")[1:]
	assert.Contains(t, result, "Found 1 file(s)")
	assert.Contains(t, result, "MODEL: Friends fight")
	assert.Contains(t, result, "\tDEFENSE: Adam's cheating is caught")
	assert.Contains(t, result, "\tDEFENSE: Bob tries to verify Adam's story")
	assert.Contains(t, result, "\tPOLICY: Honesty is the best policy")
	assert.Contains(t, result, "\t\tDEFENSE: Test honesty")
}
