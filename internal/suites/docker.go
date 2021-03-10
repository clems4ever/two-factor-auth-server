package suites

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/authelia/authelia/internal/utils"
)

// DockerEnvironment represent a docker environment.
type DockerEnvironment struct {
	dockerComposeFiles []string
}

// NewDockerEnvironment create a new docker environment.
func NewDockerEnvironment(files []string) *DockerEnvironment {
	if os.Getenv("CI") == stringTrue {
		for i := range files {
			files[i] = strings.ReplaceAll(files[i], "{}", "dist")
		}
	} else {
		for i := range files {
			files[i] = strings.ReplaceAll(files[i], "{}", "dev")
		}
	}

	return &DockerEnvironment{dockerComposeFiles: files}
}

func (de *DockerEnvironment) createCommandWithStdout(cmd string) *exec.Cmd {
	dockerCmdLine := fmt.Sprintf("docker-compose -p authelia -f %s %s", strings.Join(de.dockerComposeFiles, " -f "), cmd)
	log.Trace(dockerCmdLine)

	return utils.CommandWithStdout("bash", "-c", dockerCmdLine)
}

func (de *DockerEnvironment) createCommand(cmd string) *exec.Cmd {
	dockerCmdLine := fmt.Sprintf("docker-compose -p authelia -f %s %s", strings.Join(de.dockerComposeFiles, " -f "), cmd)
	log.Trace(dockerCmdLine)

	return utils.Command("bash", "-c", dockerCmdLine)
}

// Up spawn a docker environment.
func (de *DockerEnvironment) Up() error {
	return de.createCommandWithStdout("up --build -d").Run()
}

// Restart restarts a service.
func (de *DockerEnvironment) Restart(service string) error {
	return de.createCommandWithStdout(fmt.Sprintf("restart %s", service)).Run()
}

// Stop a docker service.
func (de *DockerEnvironment) Stop(service string) error {
	return de.createCommandWithStdout(fmt.Sprintf("stop %s", service)).Run()
}

// Start a docker service.
func (de *DockerEnvironment) Start(service string) error {
	return de.createCommandWithStdout(fmt.Sprintf("start %s", service)).Run()
}

// Down destroy a docker environment.
func (de *DockerEnvironment) Down() error {
	return de.createCommandWithStdout("down -v").Run()
}

// Exec execute a command within a given service of the environment.
func (de *DockerEnvironment) Exec(service string, command []string) (string, error) {
	cmd := de.createCommand(fmt.Sprintf("exec -T %s %s", service, strings.Join(command, " ")))
	content, err := cmd.CombinedOutput()

	return string(content), err
}

// Logs get logs of a given service of the environment.
func (de *DockerEnvironment) Logs(service string, flags []string) (string, error) {
	cmd := de.createCommand(fmt.Sprintf("logs %s %s", strings.Join(flags, " "), service))
	content, err := cmd.Output()

	return string(content), err
}
