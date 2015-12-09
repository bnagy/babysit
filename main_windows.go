package main

import (
	"flag"
	"fmt"
	"github.com/bnagy/w32"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	MAX_WAIT = uint32(3600 * 1000) // 1hr
)

var (
	subRegex = regexp.MustCompile("@@")

	flagTimeout = flag.Int("t", 10000, "Process wait timeout (ms), -1 for INFINITE") // 10s
	flagInputs  = flag.String("i", "", "Glob for input files, eg: c:\\files\\*.pdf")
)

func RunTarget(cmd, input string) (code uintptr, e error) {

	cmd = subRegex.ReplaceAllString(cmd, input)

	pi, e := w32.CreateProcessQuick(cmd)
	if e != nil {
		log.Fatalf("[!!] Failed to create %s: %s", cmd, e)
	} else {
		// log.Printf("[OK] Created process %q with handle 0x%x, PID %d", cmd, pi.Process, pi.ProcessId)
	}

	e = w32.WaitForSingleObject(pi.Process, uint32(*flagTimeout))
	if e != nil {
		log.Fatalf("[!!] Failed in WaitForSingleObject: %s", e)
	}

	code, e = w32.GetExitCodeProcess(pi.Process)
	if e != nil {
		log.Fatalf("[!!] Failed to get exit code for PID %d: %s", pi.ProcessId, e)
	} else {
		if code != w32.ERROR_SUCCESS {
			log.Printf("[!!] %q 0x%x\n", input, code)
		}
	}

	w32.CloseHandle(pi.Process)
	w32.CloseHandle(pi.Thread)
	return
}

func main() {

	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			"\n%s runs a command with a set of inputs and records nonzero exit codes\n",
			path.Base(os.Args[0]),
		)
		fmt.Fprintf(
			os.Stderr,
			"Usage: %s -i c:\\inputs\\*.txt C:\\path\\to\\target -in @@ -other -targetopts\n",
			path.Base(os.Args[0]),
		)
		fmt.Fprintf(os.Stderr, "( @@ will be substituted with each input )\n\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
	}

	flag.Parse()

	// Validate timeout
	if uint32(*flagTimeout) > MAX_WAIT && *flagTimeout != -1 {
		fmt.Fprintf(os.Stderr, "[!!] Wait timeout (%d)s too long!\n", *flagTimeout/1000)
		flag.Usage()
		os.Exit(1)
	}

	//Make sure we have some input files
	mark := time.Now()
	matches, err := filepath.Glob(*flagInputs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[!!] Bad pattern %q: %s\n", *flagInputs, err)
		flag.Usage()
		os.Exit(1)
	}
	if len(matches) == 0 {
		fmt.Fprintf(os.Stderr, "[!!] No matching inputs found for %q.\n", *flagInputs)
		flag.Usage()
		os.Exit(1)
	}
	log.Printf("[OK] Found %d input files.", len(matches))

	// Make sure we have at least one substitution marker
	targetCmd := strings.Join(flag.Args(), " ")
	if !subRegex.MatchString(targetCmd) {
		fmt.Fprintf(os.Stderr, "[!!] No substitute markers (@@) in target command!\n")
		flag.Usage()
		os.Exit(1)
	}

	for _, s := range matches {
		RunTarget(targetCmd, s)
	}

	log.Printf("[OK] All done. %d files in %s (%.2f/s)",
		len(matches),
		time.Since(mark).String(),
		float64(len(matches))/float64(time.Since(mark).Seconds()),
	)

}
