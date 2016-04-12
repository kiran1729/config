// Copyright 2014 ZeroStack, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
  "flag"
  "os"
  "testing"

  "github.com/BurntSushi/toml"
  "github.com/stretchr/testify/assert"
)

type mockCfgType struct {
  Address    string `toml:"address"`
  NumSamples uint32 `toml:"num_samples"`
}

// Default config values if not loaded from file and flags.
var mockCfg = &mockCfgType{
  Address:    "localhost:1000",
  NumSamples: 6,
}

type mockSuperCfgType struct {
  MockSection *mockCfgType `toml:"mock_section"`
}

// TestRegister makes sure that the flags can be registered based on the
// struct definition for mockCfg.
func TestRegister(t *testing.T) {

  flagsInf, err := RegisterFlags(mockCfg)
  assert.Nil(t, err)
  assert.NotNil(t, flagsInf)

  flags, ok := flagsInf.(*mockCfgType)
  assert.True(t, ok)

  // Call Parse for the supplied command line flags to be picked up after
  // the RegisterFlags call above.
  flag.Parse()

  sc := mockSuperCfgType{MockSection: mockCfg}
  _, err = toml.DecodeFile("test_config.toml", &sc)
  assert.NoError(t, err)

  // TODO: need to run test with flags injected to test this part.
  err = CheckFlagOverride(mockCfg, flags)
  assert.Nil(t, err)

  addrFlag := flag.Lookup("address")
  assert.NotNil(t, addrFlag)
  assert.Equal(t, addrFlag.DefValue, "localhost:1000")

  // Check flag got overridden by value from file.
  assert.Equal(t, sc.MockSection.Address, "localhost:3000")

  hstatSamplesFlag := flag.Lookup("num_samples")
  assert.NotNil(t, hstatSamplesFlag)
  assert.Equal(t, hstatSamplesFlag.DefValue, "6")

  // Check flag got overridden by value from file.
  assert.Equal(t, sc.MockSection.NumSamples, uint32(20))

  outFile := "/tmp/test_write.toml"
  // modify a value and write it
  sc.MockSection.Address = "notlocalhost:2000"
  err = WriteConfig(sc, outFile)
  defer os.Remove(outFile)

  assert.NoError(t, err)

  newSC := &mockSuperCfgType{MockSection: mockCfg}
  _, err = toml.DecodeFile(outFile, &newSC)
  assert.Nil(t, err)

  assert.Equal(t, newSC.MockSection.Address, "notlocalhost:2000")

}
