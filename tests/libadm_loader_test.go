package test

import (
	"fmt"
	"libadm/loaders"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

/////////////////////////////////////////////////
// Unit tests


func TestGherkinLoadFailure(t *testing.T) {
	file, err := os.ReadFile("./examples/basic/faulty-admspec.adm")
	if err == nil{
		gherkinModel, err1 := loaders.LoadGherkinContent(string(file))
		if err1 != nil {
			fmt.Println(err1)
		}
		assert.Nil(t, gherkinModel)
	} else {
		t.Error(err)
	}
}
func TestGherkinLoader(t *testing.T) {
	testVectors := map[string][]string{ // The last item in args list is the expected value
		"LoadSuccess":					{"./examples/oauth/secrets-keys.adm", "Client secrets and keys"},
	}

	for name, args := range testVectors {
		t.Run(name, func(t *testing.T) {
			
			expected := args[1]
			
			file, err := os.ReadFile(args[0])
			if err == nil{
				gherkinModel, err1 := loaders.LoadGherkinContent(string(file))
				if err1 != nil {
					fmt.Println(err1)
				}
				assert.Equal(t, gherkinModel.Feature.Name, expected)
			} else {
				fmt.Println(err)
			}
		})
	}
}