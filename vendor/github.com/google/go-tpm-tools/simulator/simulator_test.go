/*
 * Copyright 2018 Google Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not
 * use this file except in compliance with the License. You may obtain a copy of
 * the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations under
 * the License.
 */

package simulator

import (
	"crypto/rsa"
	"io"
	"math/big"
	"testing"

	"github.com/google/go-tpm-tools/tpm2tools"
	"github.com/google/go-tpm/tpm2"
)

func getSimulator(t *testing.T) *Simulator {
	t.Helper()
	simulator, err := Get()
	if err != nil {
		t.Fatal(err)
	}
	return simulator
}

func getEKModulus(t *testing.T, rwc io.ReadWriteCloser) *big.Int {
	t.Helper()
	ek, err := tpm2tools.EndorsementKeyRSA(rwc)
	if err != nil {
		t.Fatal(err)
	}
	defer ek.Close()

	return ek.PublicKey().(*rsa.PublicKey).N
}

func TestResetDoesntChangeEK(t *testing.T) {
	s := getSimulator(t)
	defer tpm2tools.CheckedClose(t, s)

	modulus1 := getEKModulus(t, s)
	if err := s.Reset(); err != nil {
		t.Fatal(err)
	}
	modulus2 := getEKModulus(t, s)

	if modulus1.Cmp(modulus2) != 0 {
		t.Fatal("Reset() should not change the EK")
	}
}
func TestManufactureResetChangesEK(t *testing.T) {
	s := getSimulator(t)
	defer tpm2tools.CheckedClose(t, s)

	modulus1 := getEKModulus(t, s)
	if err := s.ManufactureReset(); err != nil {
		t.Fatal(err)
	}
	modulus2 := getEKModulus(t, s)

	if modulus1.Cmp(modulus2) == 0 {
		t.Fatal("ManufactureReset() should change the EK")
	}
}

func TestGetRandom(t *testing.T) {
	s := getSimulator(t)
	defer tpm2tools.CheckedClose(t, s)
	result, err := tpm2.GetRandom(s, 10)
	if err != nil {
		t.Fatalf("GetRandom: %v", err)
	}
	t.Log(result)
}

// The default EK modulus returned by the simulator when using a seed of 0.
func zeroSeedModulus() *big.Int {
	mod := new(big.Int)
	mod.SetString("16916951631746795233120676661491589156159944041454533323301360736206690950055927665898258850365255777475324525235640153431219834851979041935421083247812345676551677241639541392158486693550125570954276972465867114995062336740464652481116557477039581976647612151813804384773839359390083864432536639577227083497558006614244043011423717921293964465162166865351126036685960128739613171620392174911624095420039156957292384191548425395162459332733115699189854006301807847331248289929021522087915411000598437989788501679617747304391662751900488011803826205901900186771991702576478232121332699862815915856148442279432061762451", 10)
	return mod
}

func TestFixedSeedExpectedModulus(t *testing.T) {
	s, err := GetWithFixedSeedInsecure(0)
	if err != nil {
		t.Fatal(err)
	}
	defer tpm2tools.CheckedClose(t, s)

	modulus := getEKModulus(t, s)
	if modulus.Cmp(zeroSeedModulus()) != 0 {
		t.Fatalf("getEKModulus() = %v, want %v", modulus, zeroSeedModulus())
	}
}

func TestDifferentSeedDifferentModulus(t *testing.T) {
	s, err := GetWithFixedSeedInsecure(1)
	if err != nil {
		t.Fatal(err)
	}
	defer tpm2tools.CheckedClose(t, s)

	modulus := getEKModulus(t, s)
	if modulus.Cmp(zeroSeedModulus()) == 0 {
		t.Fatalf("Moduli should not be equal when using different seeds")
	}
}
