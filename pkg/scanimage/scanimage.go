package scanimage

import (
	"bytes"
	"fmt"
	"os/exec"
)

var availableFormats = []string{
	"pnm", "tiff", "png", "jpeg",
}

const FormatDefault = "png"

var availableResolutions = []int{
	75, 150, 300, 600, 1200,
}

const ResolutionDefault = 300

var availableModes = []string{
	"color", "gray", "lineart",
}

const ModeDefault = "color"

const GammaDefault = 2.2

type ParameterError struct {
	Message string
}

func (p *ParameterError) Error() string {
	return fmt.Sprintf("client error: %s", p.Message)
}

type ExecutionError struct {
	Message string
}

func (s *ExecutionError) Error() string {
	return fmt.Sprintf("scan error: %s", s.Message)
}

type ScanImage struct {
	scanning   bool
	Command    string
	Device     string
	Format     string
	Resolution int
	Mode       string
	Gamma      float64
}

func New() *ScanImage {
	return &ScanImage{
		scanning: false,
		Command:  "scanimage",
	}
}

func inSlice[Type comparable](slice []Type, needle Type) bool {
	for _, i := range slice {
		if i == needle {
			return true
		}
	}

	return false
}

func (s *ScanImage) augment() error {
	if s.Command == "" {
		return &ParameterError{Message: "no command specified"}
	}

	if s.Format == "" {
		s.Format = FormatDefault
	} else if !inSlice(availableFormats, s.Format) {
		return &ParameterError{Message: "unsupported format"}
	}

	if s.Resolution == 0 {
		s.Resolution = ResolutionDefault
	} else if !inSlice(availableResolutions, s.Resolution) {
		return &ParameterError{Message: "unsupported resolution"}
	}

	if s.Mode == "" {
		s.Mode = ModeDefault
	} else if !inSlice(availableModes, s.Mode) {
		return &ParameterError{Message: "unsupported mode"}
	}

	if s.Gamma == 0 {
		s.Gamma = GammaDefault
	}

	return nil
}

func (s *ScanImage) GetMimeType() string {
	if s.Format == "" {
		return "image/*"
	}

	if s.Format == "pnm" {
		return "image/x-portable-anymap"
	}

	return fmt.Sprintf("image/%s", s.Format)
}

func (s *ScanImage) Scan() (*bytes.Buffer, error) {
	if s.scanning {
		return nil, &ExecutionError{Message: "scanner is already running"}
	}

	defer func() {
		s.scanning = false
	}()

	s.scanning = true

	if err := s.augment(); err != nil {
		return nil, err
	}

	arguments := []string{
		fmt.Sprintf("--resolution=%d", s.Resolution),
		fmt.Sprintf("--mode=%s", s.Mode),
		fmt.Sprintf("--format=%s", s.Format),
		fmt.Sprintf("--gamma=%f", s.Gamma),
	}

	if s.Device != "" {
		arguments = append(arguments, fmt.Sprintf("--device-name=%s", s.Device))
	}

	cmd := exec.Command(s.Command, arguments...)

	fmt.Printf("[CMD] `%s`\n", cmd.String())

	var stdoutBuf, stderrBuf bytes.Buffer

	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	if err := cmd.Start(); err != nil {
		return nil, &ExecutionError{Message: fmt.Sprintf("cannot start `%s`", s.Command)}
	}

	if err := cmd.Wait(); err != nil {
		return nil, &ExecutionError{Message: err.Error() + " (" + stderrBuf.String() + ")"}
	}

	return &stdoutBuf, nil
}
