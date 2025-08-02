// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"os"
	"testing"
)

var testdata = `To be fair, this is just random weirdo stuff going on.
We learn something new every day.`

// TestEDCOmmandsNonInput tests ed without commands which invoke input mode.
// Input mode requires a different, more complex test setup to work.
func TestEdCommandsNonInput(t *testing.T) {
	for _, tt := range []struct {
		name     string
		cmd      string
		prompt   string
		suppress bool
		wantOut  string
	}{
		{
			name:    "CmdErr_Only_On",
			cmd:     "H\nq\n",
			wantOut: "exit\n",
		},
		{
			name:    "CmdErr_On_Off",
			cmd:     "H\nH\nq\n",
			wantOut: "exit\n",
		},
		{
			name:    "CmdErr_printErr",
			cmd:     "h\nq\n",
			wantOut: "exit\nexit\n",
		},
		{
			name:    "cmdDelete_undo",
			cmd:     "-1 d\nu\nq\n",
			wantOut: "exit\n",
		},
		{
			name:    "cmdDelete_invalidAddr",
			cmd:     "4 d\nu\nq\n",
			wantOut: "line is out of bounds\nexit\n",
		},
		{
			name:    "CmdPrint_List",
			cmd:     "-1 l\nq\n",
			wantOut: "To be fair, this is just random weirdo stuff going on.$\nexit\n",
		},
		{
			name:    "CmdPrint_Print_previous_line",
			cmd:     "-1 p\nq\n",
			wantOut: "To be fair, this is just random weirdo stuff going on.\nexit\n",
		},
		{
			name:    "CmdPrint_Print_addressed line",
			cmd:     "-1 n\nq\n",
			wantOut: "1\tTo be fair, this is just random weirdo stuff going on.\nexit\n",
		},
		{
			name:    "CmdPrint_InvalidAddr",
			cmd:     "4 l\nq\n",
			wantOut: "line is out of bounds\nexit\n",
		},
		{
			name:    "CmdWrite_simple",
			cmd:     "w\nq\n",
			wantOut: "exit\n",
		},
		{
			name:    "CmdWrite_simple_invalidAddr",
			cmd:     "4 w\nq\n",
			wantOut: "line is out of bounds\nexit\n",
		},
		{
			name:    "CmdWrite_Quit",
			cmd:     "wq\nq\n",
			wantOut: "exit\n",
		},
		{
			name:    "CmdWrite_W",
			cmd:     "W\nq\n",
			wantOut: "exit\n",
		},
		{
			name:    "CmdMark_simple",
			cmd:     "k 0\nq\n",
			wantOut: "exit\n",
		},
		{
			name:    "CmdMark_No_mark_character_supplied",
			cmd:     "-1 k\nq\n",
			wantOut: "no mark character supplied\nexit\n",
		},
		{
			name:    "cmdLine_PrintLine",
			cmd:     "- =\np\nq\n", // This equals to the input of `- =`[Enter press] `p` [Enter press]
			wantOut: "1\nWe learn something new every day.\nexit\n",
		},
		{
			name:    "cmdFile_setFile_default",
			cmd:     "f\nq\n",
			wantOut: "exit\n",
		},
		{
			name:    "CmdMove_Move_simple",
			cmd:     "1 m 2\nu\nq\n",
			wantOut: "exit\n",
		},
		{
			name:    "CmdCopy_Cut",
			cmd:     "y\nq\n",
			wantOut: "exit\n",
		},
		{
			name:    "CmdScroll",
			cmd:     "-1 z\nu\nq\n",
			wantOut: "To be fair, this is just random weirdo stuff going on.\nWe learn something new every day.\nexit\n",
		},
		{
			name:    "CmdScroll_invalidAddr",
			cmd:     "4 z\nu\nq\n",
			wantOut: "line is out of bounds\nexit\n",
		},
		{
			name:    "CmdScroll_invalid_windowsize",
			cmd:     "-1 z 0\nu\nq\n",
			wantOut: "invalid window size:  0\nexit\n",
		},
		{
			name:    "CmdSub_You_We_in_line2_1_p",
			cmd:     "2 s/(We)/You/p\nu\nq\n",
			wantOut: "You learn something new every day.\nexit\n",
		},
		{
			name:    "CmdSub_You_We_in_line2_2_l",
			cmd:     "2 s/(We)/You/l\nu\nq\n",
			wantOut: "You learn something new every day.$\nexit\n",
		},
		{
			name:    "CmdSub_You_We_in_line2_3_n",
			cmd:     "2 s/(We)/You/n\nu\nq\n",
			wantOut: "1\tYou learn something new every day.\nexit\n",
		},
		{
			name:    "CmdSub_You_We_in_line2_3_g",
			cmd:     "2 s/(We)/You/g\np\nu\nq\n",
			wantOut: "You learn something new every day.\nexit\n",
		},
		{
			name:    "CmdSub_invalidAddr",
			cmd:     "4 s/(We)/You/n\nu\nq\n",
			wantOut: "line is out of bounds\nexit\n",
		},
		{
			name:    "CmdQuit_Buffer_dirty",
			cmd:     "2 s/(We)/You/n\nq\nu\nq",
			wantOut: "1\tYou learn something new every day.\nwarning: file modified\nexit\n",
		},
		{
			name:    "CmdEdit_undo",
			cmd:     "e\nu\nq\n",
			wantOut: "87\nexit\n",
		},
		{
			name:    "CmdPaste",
			cmd:     "x\nq\n",
			wantOut: "exit\n",
		},
		{
			name:    "CmdJoin",
			cmd:     "-1,2 j\np\nu\nq\n",
			wantOut: "To be fair, this is just random weirdo stuff going on.We learn something new every day.\nexit\n",
		},
		{
			name:    "CmdJoin_invalidAddr",
			cmd:     "-3 j\nu\nq\n",
			wantOut: "line is out of bounds\nexit\n",
		},
		{
			name:    "CmdDump",
			cmd:     "D\nq\n",
			wantOut: "&{[] [To be fair, this is just random weirdo stuff going on. We learn something new every day.] [0 1] [] [0 1] false false false false 1 0 1 map[]}\nexit\n",
		},
	} {
		t.Run("Command:"+tt.cmd, func(t *testing.T) {
			tmpDir := t.TempDir()
			tmpFile, err := os.CreateTemp(tmpDir, "testfile-")
			if err != nil {
				t.Errorf("os.CreateTemp(tmpdir, `testfile-`)=file, %q, want file, nil", err)
			}
			var in, out bytes.Buffer
			in.WriteString(tt.cmd)
			tmpFile.WriteString(testdata)
			if err := runEd(&in, &out, tt.suppress, tt.prompt, tmpFile.Name()); err != nil {
				t.Errorf(`runEd(&in, &out, tt.suppress, tt.prompt, "")=%q, want nil`, err)
			}
			if out.String() != tt.wantOut {
				t.Errorf("%s failed. Got: %s, Want: %s", tt.cmd, out.String(), tt.wantOut)
			}
		})
	}
}
