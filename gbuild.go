package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// 版本信息（可通过编译时注入）
var (
	version   = "dev"
	buildTime = "unknown"
)

var (
	buildOS     = flag.String("os", "", "Target OS (linux/windows/darwin), default is current")
	buildArch   = flag.String("arch", "", "Target Arch (amd64/arm64), default is current")
	output      = flag.String("o", "", "Output file path (optional)")
	buildTags   = flag.String("tags", "", "Build tags (optional)")
	ldflags     = flag.String("ldflags", "", "Linker flags (optional)")
	showVersion = flag.Bool("version", false, "Show version information")
	showHash    = flag.Bool("hash", false, "Print SHA256 hash of output binary")
	cgoEnabled  = flag.Bool("cgo", false, "Enable CGO (default false)")
)

var supportedPlatforms = map[string][]string{
	"linux":   {"amd64", "arm64"},
	"windows": {"amd64", "arm64"},
	"darwin":  {"amd64", "arm64"},
}

func main() {
	setupUsage()
	flag.Parse()

	if *showVersion {
		printVersion()
		return
	}

	setDefaults()
	validateInput()
	outputPath := getOutputPath()
	compile(outputPath)
	if *showHash {
		printFileHash(outputPath)
	}
}

func setupUsage() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Cross-platform Builder (%s)\n\n", version)
		fmt.Fprintf(os.Stderr, "Usage: %s -os <OS> -arch <Arch> [options]\n\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\nSupported platforms:")
		for supportos, arches := range supportedPlatforms {
			fmt.Fprintf(os.Stderr, "  %-7s → %s\n", supportos, strings.Join(arches, ", "))
		}
		fmt.Fprintln(os.Stderr, "\nExamples:")
		fmt.Fprintln(os.Stderr, "  gbuild -os linux -arch amd64")
		fmt.Fprintln(os.Stderr, "  gbuild -os windows -arch arm64 -o app.exe -cgo")
		fmt.Fprintln(os.Stderr, "  gbuild -version")
	}
}

func printVersion() {
	fmt.Printf("gbuild version: %s\n", version)
	fmt.Printf("Build time:     %s\n", buildTime)
	fmt.Printf("Go version:     %s\n", runtime.Version())
}

func setDefaults() {
	if *buildOS == "" {
		*buildOS = runtime.GOOS
	}
	if *buildArch == "" {
		*buildArch = runtime.GOARCH
	}
}

func validateInput() {
	arches, ok := supportedPlatforms[*buildOS]
	if !ok {
		log.Fatalf("Unsupported OS: %s. Run `gbuild -help` for supported platforms.", *buildOS)
	}
	if !contains(arches, *buildArch) {
		log.Fatalf("Unsupported architecture %q for OS %q.", *buildArch, *buildOS)
	}
}

func getOutputPath() string {
	if *output == "" {
		exeName := fmt.Sprintf("build-%s-%s", *buildOS, *buildArch)
		if *buildOS == "windows" {
			exeName += ".exe"
		}
		return filepath.Join("bin", exeName)
	}
	return filepath.Clean(*output)
}

func compile(outputPath string) {
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	args := []string{"build", "-o", outputPath}
	if *buildTags != "" {
		args = append(args, "-tags", *buildTags)
	}
	if *ldflags != "" {
		args = append(args, "-ldflags", *ldflags)
	}

	cmd := exec.Command("go", args...)
	cmd.Env = buildEnv()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("Compilation failed: %v", err)
	}

	fmt.Printf("\n✅ Successfully built: %s\n", outputPath)
}

func buildEnv() []string {
	env := os.Environ()
	env = append(env,
		"GOOS="+*buildOS,
		"GOARCH="+*buildArch,
		fmt.Sprintf("CGO_ENABLED=%d", boolToInt(*cgoEnabled)),
	)
	return env
}

func printFileHash(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Warning: failed to read file for hashing: %v\n", err)
		return
	}
	hash := sha256.Sum256(data)
	fmt.Printf("SHA256: %x\n", hash)
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
