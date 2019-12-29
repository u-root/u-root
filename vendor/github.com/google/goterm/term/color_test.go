// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package term

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"
)

type testcol struct {
	str      stringer
	exp      string
	expNoCol string
}

// ExampleColor()
func ExampleColor() {
	g := Green("Green world")
	fmt.Println("Hello", g)
	fmt.Println(Red("Warning!"))
	if string(g) == "Green world" {
		fmt.Println(Blinking(string(g)))
	}
	var col fmt.Stringer
	atk := 20
	switch {
	case atk == 0:
		col = Blue("5 FADE OUT")
	case atk < 4:
		col = Green("4 DOUBLE TAKE")
	case atk < 10:
		col = Yellow("3 ROUND HOUSE")
	case atk < 50:
		col = Red("2 FAST PACE")
	case atk >= 50:
		col = Blinking("1 COCKED PISTOL")
	}
	fmt.Println("Defcon: ", col)
	// Output:
	// Hello [32mGreen world[39m
	// [31mWarning![39m
	// [5mGreen world[0m
	// Defcon:  [31m2 FAST PACE[39m
}

// TestColorType tests the base color types.
func TestColorType(t *testing.T) {
	// Testing through all the simple Color types.
	var res string
	colors := []testcol{
		{Green("Green"), FgGreen, "Green"},
		{Blue("Blue"), FgBlue, "Blue"},
		{Red("Red"), FgRed, "Red"},
		{Yellow("Yellow"), FgYellow, "Yellow"},
		{Magenta("Magenta"), FgMagenta, "Magenta"},
		{Cyan("Cyan"), FgCyan, "Cyan"},
		{White("White"), FgWhite, "White"},
		{Black("Black"), FgBlack, "Black"},
		{Random("Random"), "", "Random"},
		// Background
		{BGreen("BGreen"), BgGreen, "BGreen"},
		{BBlue("BBlue"), BgBlue, "BBlue"},
		{BRed("BRed"), BgRed, "BRed"},
		{BYellow("BYellow"), BgYellow, "BYellow"},
		{BRandom("BRandom"), "", "BRandom"},
		{BMagenta("BMagenta"), BgMagenta, "BMagenta"},
		{BCyan("BCyan"), BgCyan, "BCyan"},
		{BWhite("BWhite"), BgWhite, "BWhite"},
		{BBlack("BBlack"), BgBlack, "BBlack"}}

	ColorDisable()
	for i := 0; i < 100; i++ {
		for _, col := range colors {
			if col.expNoCol != col.str.String() {
				t.Errorf("ColorDisable failed got: %q want: %q", col.str.String(), col.expNoCol)
			}
		}
	}

	ColorEnable()
	for i := 0; i < 100; i++ {
		for _, col := range colors {
			pattern := fmt.Sprintf("%sm.*(39|49)m", col.exp)
			match, err := regexp.MatchString(pattern, col.str.String())
			if err != nil || !match {
				t.Errorf("TestColorType received wrong basecolor got: %q want: %q err: %v", col.str.String(), col.exp, err)
			}
			res += col.str.String() + " "
		}
	}
	t.Log(res)
}

// TestColorFmt tests out the fmt printers.
func TestColorFmt(t *testing.T) {
	nr := 37
	array := []string{"opel", "ascona"}
	colors := []testcol{
		{Green(fmt.Sprintf("Green nr: %d Array: %v", nr, array)), Greenf("Green nr: %d Array: %v", nr, array), ""},
		{Blue(fmt.Sprintf("Blue nr: %d Array: %v", nr, array)), Bluef("Blue nr: %d Array: %v", nr, array), ""},
		{Red(fmt.Sprintf("Red nr: %d Array: %v", nr, array)), Redf("Red nr: %d Array: %v", nr, array), ""},
		{Yellow(fmt.Sprintf("Yellow nr: %d Array: %v", nr, array)), Yellowf("Yellow nr: %d Array: %v", nr, array), ""},
		{Magenta(fmt.Sprintf("Magenta nr: %d Array: %v", nr, array)), Magentaf("Magenta nr: %d Array: %v", nr, array), ""},
		{Cyan(fmt.Sprintf("Cyan nr: %d Array: %v", nr, array)), Cyanf("Cyan nr: %d Array: %v", nr, array), ""},
		{White(fmt.Sprintf("White nr: %d Array: %v", nr, array)), Whitef("White nr: %d Array: %v", nr, array), ""},
		{Black(fmt.Sprintf("Black nr: %d Array: %v", nr, array)), Blackf("Black nr: %d Array: %v", nr, array), ""},
		// Background
		{BGreen(fmt.Sprintf("BGreen nr: %d Array: %v", nr, array)), BGreenf("BGreen nr: %d Array: %v", nr, array), ""},
		{BBlue(fmt.Sprintf("BBlue nr: %d Array: %v", nr, array)), BBluef("BBlue nr: %d Array: %v", nr, array), ""},
		{BRed(fmt.Sprintf("BRed nr: %d Array: %v", nr, array)), BRedf("BRed nr: %d Array: %v", nr, array), ""},
		{BYellow(fmt.Sprintf("BYellow nr: %d Array: %v", nr, array)), BYellowf("BYellow nr: %d Array: %v", nr, array), ""},
		{BMagenta(fmt.Sprintf("BMagenta nr: %d Array: %v", nr, array)), BMagentaf("BMagenta nr: %d Array: %v", nr, array), ""},
		{BCyan(fmt.Sprintf("BCyan nr: %d Array: %v", nr, array)), BCyanf("BCyan nr: %d Array: %v", nr, array), ""},
		{BWhite(fmt.Sprintf("BWhite nr: %d Array: %v", nr, array)), BWhitef("BWhite nr: %d Array: %v", nr, array), ""},
		{BBlack(fmt.Sprintf("BBlack nr: %d Array: %v", nr, array)), BBlackf("BBlack nr: %d Array: %v", nr, array), ""}}

	for _, col := range colors {
		if got, want := col.str.String(), col.exp; got != want {
			t.Errorf("got: %v want: %v", got, want)
			continue
		}
		t.Log(col.exp)
	}
}

// TestModType tests out the attributes.
func TestModType(t *testing.T) {
	var res string
	var mods = []testcol{
		{Blinking("Blinking"), Blink, "Blinking"},
		{Underline("Underline"), Underln, "Underline"},
		{Bold("Bold"), Bld, "Bold"},
		// Bright ("Bright"), Save that one for now
		{Italic("Italic"), Ital, "Italic"}}
	for _, mod := range mods {
		pattern := mod.exp + "m.*" + NoMode + "m"
		if match, err := regexp.MatchString(pattern, mod.str.String()); err != nil || !match {
			t.Errorf("TestModTupe received wrong modstring got: %q want: %q err: %v", mod.str.String(), mod.exp, err)
		}
		res += mod.str.String() + "\n"
	}
	ColorDisable()
	for _, mod := range mods {
		if mod.expNoCol != mod.str.String() {
			t.Errorf("ColorDisable failed got: %q want: %q:", mod.str.String(), mod.expNoCol)
		}
	}
	ColorEnable()
	t.Log(res)
}

// TestColor256 tests the terminal 256 color modes.
func TestColor256(t *testing.T) {
	var rstr string
	for c := 0; c <= 1; c++ {
		for i := 0; i < 256; i++ {
			col, err := NewColor256(strconv.Itoa(i), strconv.Itoa(i), "")
			if err != nil {
				t.Fatal("err: ", err)
			}
			if !colorEnable {
				if col.String() != strconv.Itoa(i) {
					t.Errorf("ColorDisable failed got: %q want: %q", col.String(), strconv.Itoa(i))
				}
				continue
			}
			pattern := F256 + ";5;" + strconv.Itoa(i) + ";" + Bg256 + ";5;m" + strconv.Itoa(i)
			pattern += ".*" + FgDefault + ";5;" + BgDefault + ";5;" + "m"
			if match, err := regexp.MatchString(pattern, col.String()); err != nil || !match {
				t.Errorf("TestColor256 received wrong foreground 256color string got: %q want: %d err: %v", col.String(), i, err)
			}
			rstr += string(col) + " "
		}
		rstr += "\n"
		for i := 0; i < 256; i++ {
			col, err := NewColor256(strconv.Itoa(i), "", strconv.Itoa(i))
			if err != nil {
				t.Fatal("err: ", err)
			}
			if !colorEnable {
				if col.String() != strconv.Itoa(i) {
					t.Errorf("ColorDisable failed got: %q want: %q", col.String(), strconv.Itoa(i))
				}
				continue
			}
			pattern := F256 + ";5;;" + Bg256 + ";5;" + strconv.Itoa(i) + "m" + strconv.Itoa(i)
			pattern += ".*" + FgDefault + ";5;" + BgDefault + ";5;" + "m"
			if match, err := regexp.MatchString(pattern, col.String()); err != nil || !match {
				t.Errorf("TestColor256 received wrong background 256color string got: %q want: %d err: %v", col.String(), i, err)
			}
			rstr += string(col) + " "
		}
		ColorDisable()
	}
	ColorEnable()
	t.Log("256 cols:\n", rstr)
}

// BenchMarkCombo tests out the speed of the combo function
func BenchmarkCombo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewCombo("RedBlinking", FgRed, Blink)
		NewCombo("GreenUnderline", FgGreen, Underln)
		NewCombo("BlueItalic", FgBlue, Ital)
		NewCombo("RedGreen", FgRed, BgGreen)
		NewCombo("BlueRedUnderlineItalic", FgBlue, BgRed, Underln, Ital)
		NewCombo("FgBlueFgRedBgYellowBgBlue", FgBlue, FgRed, BgYellow, BgBlue)
		NewCombo("UnderlnUnderlnItalicItalic", Underln, Underln, Ital, Ital)
		NewCombo("GreenYellowUnderlineBold", BgYellow, FgGreen, Underln, Bld)
	}
}

// TestCombo tests out the combination of different modes.
func TestCombo(t *testing.T) {
	var res string
	combo := []testcol{
		{NewCombo("RedBlinking", FgRed, Blink), "\x1b[31;5mRedBlinking\x1b[39;0m", "RedBlinking"},
		{NewCombo("GreenUnderline", FgGreen, Underln), "\x1b[32;4mGreenUnderline\x1b[39;0m", "GreenUnderline"},
		{NewCombo("BlueItalic", FgBlue, Ital), "\x1b[34;3mBlueItalic\x1b[39;0m", "BlueItalic"},
		{NewCombo("RedGreen", FgRed, BgGreen), "\x1b[31;42mRedGreen\x1b[39;49m", "RedGreen"},
		{NewCombo("BlueRedUnderlineItalic", FgBlue, BgRed, Underln, Ital), "\x1b[34;41;4;3mBlueRedUnderlineItalic\x1b[39;49;0m", "BlueRedUnderlineItalic"},
		{NewCombo("FgBlueFgRedBgYellowBgBlue", FgBlue, FgRed, BgYellow, BgBlue), "\x1b[34;43mFgBlueFgRedBgYellowBgBlue\x1b[39;49m", "FgBlueFgRedBgYellowBgBlue"},
		{NewCombo("UnderlnUnderlnItalicItalic", Underln, Underln, Ital, Ital), "\x1b[4;3mUnderlnUnderlnItalicItalic\x1b[0m", "UnderlnUnderlnItalicItalic"},
		{NewCombo("GreenYellowUnderlineBold", BgYellow, FgGreen, Underln, Bld), "\x1b[43;32;4;1mGreenYellowUnderlineBold\x1b[39;49;0m", "GreenYellowUnderlineBold"}}
	for _, com := range combo {
		if fmt.Sprintf("%q", com.str) != fmt.Sprintf("%q", com.exp) {
			t.Errorf("TestCombo received wrong combination string string got: %q want: %q ", com.str, com.exp)
		}
		res += com.str.String() + "\n"
	}
	ColorDisable()
	combo = []testcol{
		{NewCombo("RedBlinking", FgRed, Blink), "\x1b[31;5mRedBlinking\x1b[39;0m", "RedBlinking"},
		{NewCombo("GreenUnderline", FgGreen, Underln), "\x1b[32;4mGreenUnderline\x1b[39;0m", "GreenUnderline"},
		{NewCombo("BlueItalic", FgBlue, Ital), "\x1b[34;3mBlueItalic\x1b[39;0m", "BlueItalic"},
		{NewCombo("RedGreen", FgRed, BgGreen), "\x1b[31;42mRedGreen\x1b[39;49m", "RedGreen"},
		{NewCombo("BlueRedUnderlineItalic", FgBlue, BgRed, Underln, Ital), "\x1b[34;41;4;3mBlueRedUnderlineItalic\x1b[39;49;0m", "BlueRedUnderlineItalic"},
		{NewCombo("FgBlueFgRedBgYellowBgBlue", FgBlue, FgRed, BgYellow, BgBlue), "\x1b[34;43mFgBlueFgRedBgYellowBgBlue\x1b[39;49m", "FgBlueFgRedBgYellowBgBlue"},
		{NewCombo("UnderlnUnderlnItalicItalic", Underln, Underln, Ital, Ital), "\x1b[4;3mUnderlnUnderlnItalicItalic\x1b[0m", "UnderlnUnderlnItalicItalic"},
		{NewCombo("GreenYellowUnderlineBold", BgYellow, FgGreen, Underln, Bld), "\x1b[43;32;4;1mGreenYellowUnderlineBold\x1b[39;49;0m", "GreenYellowUnderlineBold"}}
	for _, com := range combo {
		if com.expNoCol != com.str.String() {
			t.Errorf("ColorDisable failed got: %q want: %q", com.str.String(), com.expNoCol)
		}
	}
	ColorEnable()
	t.Log(res)
}

// TestColorRGB generates a bunch of RGB colors and logs it.
// This is for a human to look at to figure out if their term handles it
// the function does not fail.
func TestColorRGB(t *testing.T) {
	var red, green, blue uint8
	var rstr string
	for c := 0; c <= 1; c++ {
		// Go through the reds
		for ; red <= 254; red++ {
			col := NewColorRGB("#", red, green, blue)
			if !colorEnable {
				if "#" != col.String() {
					t.Errorf("ColorDisable failed got: %q want: #", col.String())
				}
				continue
			}
			rstr = col.String()
		}
		rstr = "\n"
		// Green
		for ; green <= 254; green++ {
			col := NewColorRGB("#", red, green, blue)
			if !colorEnable {
				if "#" != col.String() {
					t.Errorf("ColorDisable failed got: %q want: #", col.String())
				}
				continue
			}
			rstr += col.String()

		}
		rstr += "\n"
		// Blue
		for ; blue <= 254; blue++ {
			col := NewColorRGB("#", red, green, blue)
			if !colorEnable {
				if "#" != col.String() {
					t.Errorf("ColorDisable failed got: %q want: #", col.String())
				}
				continue
			}
			rstr += col.String()

		}
		rstr += "\n"
		for grey := uint8(0); grey <= 254; grey++ {
			col := NewColorRGB("#", grey, grey, grey)
			if !colorEnable {
				if "#" != col.String() {
					t.Errorf("ColorDisable failed got: %q want: #", col.String())
				}
				continue
			}
			rstr += col.String()
		}
		ColorDisable()
	}
	ColorEnable()
	t.Log(rstr)
}

// TestTheTest runs the big test again logging it for a human to look at.
func TestTheTest(t *testing.T) {
	t.Log(TestTerm())
	t.Log("ColorDisable")
	ColorDisable()
	t.Log(TestTerm())
	ColorEnable()
}
