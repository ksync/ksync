package input

import (
  "testing"

  "github.com/stretchr/testify/assert"

)

var tl *Locator

func init(){
  tl = &Locator{
      PodName:        "testname",
      Selector:       "testselector",
      ContainerName:  "testcontainername",
    }
}

func TestValidator(t *testing.T) {
  assert.NotPanics(t, func() {tl.Validator()})
}

func TestContainers(t *testing.T) {
  _, err := &tl.Containers()

  assert.NoError(t, err)
  tl = &Locator{PodName: "",}
}
