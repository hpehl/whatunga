package shell

/*
#cgo darwin CFLAGS: -I/usr/local/opt/readline/include
#cgo darwin LDFLAGS: -L/usr/local/opt/readline/lib
#cgo LDFLAGS: -lreadline -lhistory

#include <stdio.h>
#include <stdlib.h>
#include <readline/readline.h>
#include <readline/history.h>

extern char *_completion_fn(char *s, int i);

static char *_completion_fn_trans(const char *s, int i) {
	return _completion_fn((char *) s, i);
}

static void register_readline() {
	rl_completion_entry_function = _completion_fn_trans;
	using_history();
}
*/
import "C"

import (
	"io"
	"os"
	"os/signal"
	"syscall"
	"unsafe"
)

// The prompt used by Reader(). The prompt can contain ANSI escape
// sequences, they will be escaped as necessary.
var Prompt = "> "

// The continue prompt used by Reader(). The prompt can contain ANSI escape
// sequences, they will be escaped as necessary.
var Continue = ".."

// The readline package adds a signal handler for SIGINT at init. If
// CatchSigint is true, upon receiving the signal (typically from the
// user pressing Ctrl+C) it will restore the terminal attributes and
// call os.Exit(1).
//
// Applications that install their own SIGINT handler should set this
// variable to false, and call Cleanup() manually if the handler
// causes the application to terminate while a String() call is
// running.
var CatchSigint = true

// If CompletionAppendChar is non-zero, readline will append the
// corresponding character to the prompt after each completion. A
// typical value would be a space.
var CompletionAppendChar = 0

// This function provides entries for the tab completer.
var Completer = func(query, ctx string) []string {
	return nil
}

var entries []*C.char

// Read a line with the given prompt. The prompt can contain ANSI
// escape sequences, they will be escaped as necessary.
func String(prompt string) (string, error) {
	p := C.CString(prompt)
	rp := C.readline(p)
	s := C.GoString(rp)
	C.free(unsafe.Pointer(p))
	if rp != nil {
		C.free(unsafe.Pointer(rp))
		return s, nil
	}
	return s, io.EOF
}

// This function can be assigned to the Completer variable to use
// readline's default filename completion, or it can be called by a
// custom completer function to get a list of files and filter it.
func FilenameCompleter(query, ctx string) []string {
	var compls []string
	var c *C.char
	q := C.CString(query)

	for i := 0; ; i++ {
		if c = C.rl_filename_completion_function(q, C.int(i)); c == nil {
			break
		}
		compls = append(compls, C.GoString(c))
		C.free(unsafe.Pointer(c))
	}

	C.free(unsafe.Pointer(q))

	return compls
}

//export _completion_fn
func _completion_fn(p *C.char, _i C.int) *C.char {
	C.rl_completion_append_character = C.int(CompletionAppendChar)
	i := int(_i)
	if i == 0 {
		es := Completer(C.GoString(p), C.GoString(C.rl_line_buffer))
		entries = make([]*C.char, len(es))
		for i, x := range es {
			entries[i] = C.CString(x)
		}
	}
	if i >= len(entries) {
		return nil
	}
	return entries[i]
}

func SetWordBreaks(cs string) {
	C.rl_completer_word_break_characters = C.CString(cs)
}

// Add an item to the history.
func AddHistory(s string) {
	n := HistorySize()
	if n == 0 || s != GetHistory(n-1) {
		C.add_history(C.CString(s))
	}
}

// Retrieve a line from the history.
func GetHistory(i int) string {
	e := C.history_get(C.int(i + 1))
	if e == nil {
		return ""
	}
	return C.GoString(e.line)
}

// Deletes all the items in the history.
func ClearHistory() {
	C.clear_history()
}

// Returns the number of items in the history.
func HistorySize() int {
	return int(C.history_length)
}

// Load the history from a file.
func LoadHistory(path string) error {
	p := C.CString(path)
	e := C.read_history(p)
	C.free(unsafe.Pointer(p))

	if e == 0 {
		return nil
	}
	return syscall.Errno(e)
}

// Save the history to a file.
func SaveHistory(path string) error {
	p := C.CString(path)
	e := C.write_history(p)
	C.free(unsafe.Pointer(p))

	if e == 0 {
		return nil
	}
	return syscall.Errno(e)
}

// Cleanup() frees internal memory and restores terminal
// attributes. This function should be called when program execution
// stops before the return of a String() call, so as not to leave the
// terminal in a corrupted state.
//
// If the CatchSigint variable is set to true (default), Cleanup() is
// called automatically on reception of a SIGINT signal.
func Cleanup() {
	C.rl_free_line_state()
	C.rl_cleanup_after_signal()
}

func handleSignals() {
	C.rl_catch_signals = 0
	C.rl_catch_sigwinch = 0

	signals := make(chan os.Signal, 2)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGWINCH)

	for s := range signals {
		switch s {
		case syscall.SIGWINCH:
			C.rl_resize_terminal()

		case syscall.SIGINT:
			if CatchSigint {
				Cleanup()
				os.Exit(1)
			}
		}
	}
}

func init() {
	go handleSignals()
	C.register_readline()
}
