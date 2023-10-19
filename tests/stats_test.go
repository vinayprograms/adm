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
	out, _, log := harness.ReadAndRelease()

	assert.Nil(t, err)
	// only the last line contains the required output
	out_result := strings.Split(out, "\n")[1:]
	log_result := strings.Split(log, "\n")
	assert.Contains(t, log_result[0], "Found 1 file(s)")
	assert.Contains(t, out_result, "MODEL: Client secrets and keys")
	assert.Contains(t, out_result, "\tATTACK: Secrets extracted via path traversal")
	assert.Contains(t, out_result, "\tATTACK: Secrets extracted via web-server path traversal")
	assert.Contains(t, out_result, "\tATTACK: Secrets extracted via encoded path in URL")
	assert.Contains(t, out_result, "\tATTACK: Secrets extracted via encoded path in URL sent to Apache Web Server")
	assert.Contains(t, out_result, "\tATTACK: Compromize web server")
	assert.Contains(t, out_result, "\tATTACK: Secrets extracted from code")
	assert.Contains(t, out_result, "\tATTACK: Secrets extracted from app binary")
	assert.Contains(t, out_result, "\tDEFENSE: Block encoded paths")
	assert.Contains(t, out_result, "\tDEFENSE: Don't allow path traversal to configuration files")
	assert.Contains(t, out_result, "\tDEFENSE: Web-server shouldn't serve config files")
	assert.Contains(t, out_result, "\tDEFENSE: Encrypt secrets before storing in file")
	assert.Contains(t, out_result, "\tDEFENSE: Don't publish code containing secrets")
	assert.Contains(t, out_result, "\tDEFENSE: Don't publish app binary containing secrets")
	assert.Contains(t, out_result, "\tDEFENSE: Don't store secrets in config file")
	assert.Contains(t, out_result, "\tDEFENSE: use Apache Web Server v2.5.51 or later")
	assert.Contains(t, out_result, "\tDEFENSE: Store secrets and keys in vault")
}

func TestStatsWithPolicies(t *testing.T) {
	args := []string{"stat", "./examples/friends.adm"}
	harness := output_interceptor{}
	harness.Hook()

	err := sendToParseArgs(args)
	out, _, log := harness.ReadAndRelease()

	assert.Nil(t, err)
	// only the last line contains the required output
	out_result := strings.Split(out, "\n")[1:]
	log_result := strings.Split(log, "\n")
	assert.Contains(t, log_result[0], "Found 1 file(s)")
	assert.Contains(t, out_result, "MODEL: Friends fight")
	assert.Contains(t, out_result, "\tATTACK: Adam cheats Bob")
	assert.Contains(t, out_result, "\tATTACK: Adam hides the cheat")
	assert.Contains(t, out_result, "\tDEFENSE: Adam's cheating is caught")
	assert.Contains(t, out_result, "\tDEFENSE: Bob tries to verify Adam's story")
	assert.Contains(t, out_result, "\tPOLICY: Honesty is the best policy")
	assert.Contains(t, out_result, "\t\tDEFENSE: Test honesty")
}

func TestStatsWithPreemtiveDefensesFlag(t *testing.T) {
	args := []string{"stat", "-p", "./examples/friends.adm"}
	harness := output_interceptor{}
	harness.Hook()

	err := sendToParseArgs(args)
	out, _, log := harness.ReadAndRelease()

	assert.Nil(t, err)
	// only the last line contains the required output
	out_result := strings.Split(out, "\n")[1:]
	log_result := strings.Split(log, "\n")
	assert.Contains(t, log_result[0], "Found 1 file(s)")
	assert.Contains(t, out_result, "MODEL: Friends fight")
	assert.Contains(t, out_result, "\tPOLICY: Honesty is the best policy")
	assert.Contains(t, out_result, "\t\tDEFENSE: Test honesty")
}

func TestStatsWithIRFlag(t *testing.T) {
	args := []string{"stat", "-i", "./examples/friends.adm"}
	harness := output_interceptor{}
	harness.Hook()

	err := sendToParseArgs(args)
	out, _, log := harness.ReadAndRelease()

	assert.Nil(t, err)
	// only the last line contains the required output
	out_result := strings.Split(out, "\n")[1:]
	log_result := strings.Split(log, "\n")
	assert.Contains(t, log_result[0], "Found 1 file(s)")
	assert.Contains(t, out_result, "MODEL: Friends fight")
	assert.Contains(t, out_result, "\tDEFENSE: Adam's cheating is caught")
	assert.Contains(t, out_result, "\tDEFENSE: Bob tries to verify Adam's story")
}

func TestStatsWithDefenseAndIRFlags(t *testing.T) {
	args := []string{"stat", "-d", "-i", "./examples/friends.adm"}
	harness := output_interceptor{}
	harness.Hook()

	err := sendToParseArgs(args)
	out, _, log := harness.ReadAndRelease()

	assert.Nil(t, err)
	// only the last line contains the required output
	out_result := strings.Split(out, "\n")[1:]
	log_result := strings.Split(log, "\n")
	assert.Contains(t, log_result[0], "Found 1 file(s)")
	assert.Contains(t, out_result, "MODEL: Friends fight")
	assert.Contains(t, out_result, "\tDEFENSE: Adam's cheating is caught")
	assert.Contains(t, out_result, "\tDEFENSE: Bob tries to verify Adam's story")
	assert.Contains(t, out_result, "\tPOLICY: Honesty is the best policy")
	assert.Contains(t, out_result, "\t\tDEFENSE: Test honesty")
}
