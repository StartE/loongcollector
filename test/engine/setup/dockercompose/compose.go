// Copyright 2021 iLogtail Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dockercompose

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	docker "github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gopkg.in/yaml.v3"

	"github.com/alibaba/ilogtail/pkg/logger"
	"github.com/alibaba/ilogtail/test/config"
)

const (
	composeCategory = "docker-compose"
	finalFileName   = "testcase-compose.yaml"
	template        = `version: '3.8'
services:
  loongcollectorC:
    image: aliyun/loongcollector:0.0.1
    hostname: loongcollector
    privileged: true
    pid: host
    volumes:
      - %s:/loongcollector/conf/default_flusher.json
      - %s:/loongcollector/conf/continuous_pipeline_config/local
      - /:/logtail_host
      - /var/run/docker.sock:/var/run/docker.sock
      - /sys/:/sys/
    ports:
      - 18689:18689
    environment:
      - LOGTAIL_FORCE_COLLECT_SELF_TELEMETRY=true
      - LOGTAIL_DEBUG_FLAG=true
      - LOGTAIL_AUTO_PROF=false
      - LOGTAIL_HTTP_LOAD_CONFIG=true
      - ALICLOUD_LOG_DOCKER_ENV_CONFIG=true
      - ALICLOUD_LOG_PLUGIN_ENV_CONFIG=false
      - ALIYUN_LOGTAIL_USER_DEFINED_ID=1111
    healthcheck:
      test: "cat /loongcollector/log/loongcollector.LOG"
      interval: 15s
      timeout: 5s
`
)

// ComposeBooter control docker-compose to start or stop containers.
type ComposeBooter struct {
	cli       *client.Client
	logtailID string
}

// StrategyWrapper caches the real wait.StrategyTarget to get the metadata of the running containers.
type StrategyWrapper struct {
	original    wait.Strategy
	target      wait.StrategyTarget
	service     string
	virtualPort nat.Port
}

// NewComposeBooter create a new compose booter.
func NewComposeBooter() *ComposeBooter {
	return &ComposeBooter{}
}

func (c *ComposeBooter) Start(ctx context.Context) error {
	if err := c.createComposeFile(ctx); err != nil {
		return err
	}
	projectName := strings.Split(config.CaseHome, "/")[len(strings.Split(config.CaseHome, "/"))-2]
	hasher := sha256.New()
	hasher.Write([]byte(projectName))
	projectName = fmt.Sprintf("%x", hasher.Sum(nil))
	compose := testcontainers.NewLocalDockerCompose([]string{config.CaseHome + finalFileName}, projectName).WithCommand([]string{"up", "-d", "--build"})
	strategyWrappers := withExposedService(compose)
	// retry 3 times
	for i := 0; i < 3; i++ {
		execError := compose.Invoke()
		if execError.Error == nil {
			break
		}
		logger.Error(context.Background(), "START_DOCKER_COMPOSE_ERROR",
			"stdout", execError.Error.Error())
		if i == 2 {
			return execError.Error
		}
		execError = testcontainers.NewLocalDockerCompose([]string{config.CaseHome + finalFileName}, projectName).Down()
		if execError.Error != nil {
			logger.Error(context.Background(), "DOWN_DOCKER_COMPOSE_ERROR",
				"stdout", execError.Error.Error())
			return execError.Error
		}
	}
	cli, err := CreateDockerClient()
	if err != nil {
		return err
	}
	c.cli = cli

	list, err := cli.ContainerList(context.Background(), types.ContainerListOptions{
		Filters: filters.NewArgs(
			filters.Arg("name", fmt.Sprintf("%s_loongcollectorC*", projectName)),
			filters.Arg("name", fmt.Sprintf("%s-loongcollectorC*", projectName)),
		),
	})
	if len(list) != 1 {
		logger.Errorf(context.Background(), "LOGTAIL_COMPOSE_ALARM", "logtail container size is not equal 1, got %d count", len(list))
		return err
	}
	c.logtailID = list[0].ID

	// the docker engine cannot access host on the linux platform, more details please see: https://github.com/docker/for-linux/issues/264
	cmd := []string{
		"sh",
		"-c",
		"env |grep HOST_OS|grep Linux && ip -4 route list match 0/0|awk '{print $3\" host.docker.internal\"}' >> /etc/hosts",
	}
	if err = c.exec(c.logtailID, cmd); err != nil {
		return err
	}
	err = registerDockerNetMapping(strategyWrappers)
	logger.Debugf(context.Background(), "registered net mapping: %v", networkMapping)
	return err
}

func (c *ComposeBooter) Stop() error {
	projectName := strings.Split(config.CaseHome, "/")[len(strings.Split(config.CaseHome, "/"))-2]
	hasher := sha256.New()
	hasher.Write([]byte(projectName))
	projectName = fmt.Sprintf("%x", hasher.Sum(nil))
	execError := testcontainers.NewLocalDockerCompose([]string{config.CaseHome + finalFileName}, projectName).Down()
	if execError.Error != nil {
		logger.Error(context.Background(), "STOP_DOCKER_COMPOSE_ERROR",
			"stdout", execError.Stdout.Error(), "stderr", execError.Stderr.Error())
		return execError.Error
	}
	_ = os.Remove(config.CaseHome + finalFileName)
	return nil
}

func (c *ComposeBooter) exec(id string, cmd []string) error {
	cfg := dockertypes.ExecConfig{
		User: "root",
		Cmd:  cmd,
	}
	resp, err := c.cli.ContainerExecCreate(context.Background(), id, cfg)
	if err != nil {
		logger.Errorf(context.Background(), "DOCKER_EXEC_ALARM", "cannot create exec config: %v", err)
		return err
	}
	err = c.cli.ContainerExecStart(context.Background(), resp.ID, dockertypes.ExecStartCheck{
		Detach: false,
		Tty:    false,
	})
	if err != nil {
		logger.Errorf(context.Background(), "DOCKER_EXEC_ALARM", "cannot start exec config: %v", err)
		return err
	}
	return nil
}

func (c *ComposeBooter) CopyCoreLogs() {
	if c.logtailID != "" {
		_ = os.Remove(config.LogDir)
		_ = os.Mkdir(config.LogDir, 0750)
		cmd := exec.Command("docker", "cp", c.logtailID+":/loongcollector/log/loongcollector.LOG", config.LogDir)
		output, err := cmd.CombinedOutput()
		logger.Debugf(context.Background(), "\n%s", string(output))
		if err != nil {
			logger.Error(context.Background(), "COPY_LOG_ALARM", "type", "main", "err", err)
		}
		cmd = exec.Command("docker", "cp", c.logtailID+":/loongcollector/log/go_plugin.LOG", config.LogDir)
		output, err = cmd.CombinedOutput()
		logger.Debugf(context.Background(), "\n%s", string(output))
		if err != nil {
			logger.Error(context.Background(), "COPY_LOG_ALARM", "type", "plugin", "err", err)
		}
	}
}

func (s *StrategyWrapper) WaitUntilReady(ctx context.Context, target wait.StrategyTarget) error {
	s.target = target
	return s.original.WaitUntilReady(ctx, target)
}

func (c *ComposeBooter) createComposeFile(ctx context.Context) error {
	// read the case docker compose file.
	if _, err := os.Stat(config.CaseHome); os.IsNotExist(err) {
		if err = os.MkdirAll(config.CaseHome, 0750); err != nil {
			return err
		}
	}
	_, err := os.Stat(config.CaseHome + config.DockerComposeFileName)
	var bytes []byte
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		if bytes, err = os.ReadFile(config.CaseHome + config.DockerComposeFileName); err != nil {
			return err
		}
	}
	cfg := c.getLogtailpluginConfig()
	services := cfg["services"].(map[string]interface{})
	loongcollector := services["loongcollectorC"].(map[string]interface{})
	// merge docker compose file.
	if len(bytes) > 0 {
		caseCfg := make(map[string]interface{})
		if err = yaml.Unmarshal(bytes, &caseCfg); err != nil {
			return err
		}
		// depend on
		loongcollectorDependOn := map[string]interface{}{}
		if dependOnContainers, ok := ctx.Value(config.DependOnContainerKey).([]string); ok {
			for _, container := range dependOnContainers {
				loongcollectorDependOn[container] = map[string]string{
					"condition": "service_healthy",
				}
			}
		}
		newServices := caseCfg["services"].(map[string]interface{})
		for k := range newServices {
			services[k] = newServices[k]
		}
		loongcollector["depends_on"] = loongcollectorDependOn
	}
	// volume
	loongcollectorMount := services["loongcollectorC"].(map[string]interface{})["volumes"].([]interface{})
	if volumes, ok := ctx.Value(config.MountVolumeKey).([]string); ok {
		for _, volume := range volumes {
			loongcollectorMount = append(loongcollectorMount, volume)
		}
	}
	// ports
	loongcollectorPort := services["loongcollectorC"].(map[string]interface{})["ports"].([]interface{})
	if ports, ok := ctx.Value(config.ExposePortKey).([]string); ok {
		for _, port := range ports {
			loongcollectorPort = append(loongcollectorPort, port)
		}
	}
	loongcollector["volumes"] = loongcollectorMount
	loongcollector["ports"] = loongcollectorPort
	yml, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(config.CaseHome+finalFileName, yml, 0600)
}

// getLogtailpluginConfig find the docker compose configuration of the loongcollector.
func (c *ComposeBooter) getLogtailpluginConfig() map[string]interface{} {
	cfg := make(map[string]interface{})
	str := fmt.Sprintf(template, config.FlusherFile, config.ConfigDir)
	if err := yaml.Unmarshal([]byte(str), &cfg); err != nil {
		panic(err)
	}
	bytes, _ := yaml.Marshal(cfg)
	println(string(bytes))
	return cfg
}

// registerDockerNetMapping
func registerDockerNetMapping(wrappers []*StrategyWrapper) error {
	for _, wrapper := range wrappers {
		host, err := wrapper.target.Host(context.Background())
		if err != nil {
			return err
		}
		physicalPort, err := wrapper.target.MappedPort(context.Background(), wrapper.virtualPort)
		if err != nil {
			return err
		}
		physicalPort.Int()
		registerNetMapping(wrapper.service+":"+wrapper.virtualPort.Port(), host+":"+physicalPort.Port())
	}
	return nil
}

// withExposedService add wait.Strategy to the docker compose.
func withExposedService(compose testcontainers.DockerCompose) (wrappers []*StrategyWrapper) {
	localCompose := compose.(*testcontainers.LocalDockerCompose)
	for serv, rawCfg := range localCompose.Services {
		cfg := rawCfg.(map[string]interface{})
		rawPorts, ok := cfg["ports"]
		if !ok {
			continue
		}
		ports := rawPorts.([]interface{})
		for i := 0; i < len(ports); i++ {
			var wrapper StrategyWrapper
			var portNum int
			var np nat.Port
			switch p := ports[i].(type) {
			case int:
				portNum = p
				np = nat.Port(fmt.Sprintf("%d/tcp", p))
			case string:
				relations := strings.Split(p, ":")
				proto, port := nat.SplitProtoPort(relations[len(relations)-1])
				np = nat.Port(fmt.Sprintf("%s/%s", port, proto))
				portNum, _ = strconv.Atoi(port)
			}
			wrapper = StrategyWrapper{
				original:    wait.ForListeningPort(np),
				service:     serv,
				virtualPort: np,
			}
			localCompose.WithExposedService(serv, portNum, &wrapper)
			wrappers = append(wrappers, &wrapper)
		}
	}
	return
}

func CreateDockerClient(opt ...docker.Opt) (client *docker.Client, err error) {
	opt = append(opt, docker.FromEnv)
	client, err = docker.NewClientWithOpts(opt...)
	if err != nil {
		return nil, err
	}
	// add dockerClient connectivity tests
	pingCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	ping, err := client.Ping(pingCtx)
	if err != nil {
		return nil, err
	}
	client.NegotiateAPIVersionPing(ping)
	return
}
