//go:build integration && kubernetes

package kubernetes_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/gitlab-runner/common"
	"gitlab.com/gitlab-org/gitlab-runner/common/buildtest"
	"gitlab.com/gitlab-org/gitlab-runner/executors/kubernetes"
	"gitlab.com/gitlab-org/gitlab-runner/executors/kubernetes/internal/pull"
	"gitlab.com/gitlab-org/gitlab-runner/helpers"
	"gitlab.com/gitlab-org/gitlab-runner/helpers/dns"
	"gitlab.com/gitlab-org/gitlab-runner/helpers/featureflags"
	"gitlab.com/gitlab-org/gitlab-runner/helpers/retry"
	"gitlab.com/gitlab-org/gitlab-runner/session"
	"gitlab.com/gitlab-org/gitlab-runner/shells"
	"gitlab.com/gitlab-org/gitlab-runner/shells/shellstest"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	testFeatureFlag      string
	testFeatureFlagValue bool
	ciNamespace          = "k8s-runner-integration-tests"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	testFeatureFlag = os.Getenv("CI_RUNNER_TEST_FEATURE_FLAG")
	if testFeatureFlag != "" {
		var err error
		testFeatureFlagValue, err = strconv.ParseBool(os.Getenv("CI_RUNNER_TEST_FEATURE_FLAG_VALUE"))
		if err != nil {
			panic(err)
		}
	}

	if os.Getenv("CI_RUNNER_TEST_NAMESPACE") != "" {
		ciNamespace = os.Getenv("CI_RUNNER_TEST_NAMESPACE")
	}
}

type kubernetesNamespaceManagerAction int64

const (
	createNamespace kubernetesNamespaceManagerAction = iota
	deleteNamespace
	// counterServiceImage counts to 10 and exits
	counterServiceImage = "registry.gitlab.com/gitlab-org/gitlab-runner/test/counter-service:v1"
)

type namespaceManager struct {
	action      kubernetesNamespaceManagerAction
	namespace   string
	client      *k8s.Clientset
	maxAttempts int
	timeout     time.Duration
}

func newNamespaceManager(client *k8s.Clientset, action kubernetesNamespaceManagerAction, namespace string) *namespaceManager {
	return &namespaceManager{
		namespace:   namespace,
		action:      action,
		client:      client,
		maxAttempts: 3,
		timeout:     time.Minute,
	}
}

func (n *namespaceManager) Run() (*v1.Namespace, error) {
	var err error
	var ns *v1.Namespace

	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	switch n.action {
	case createNamespace:
		k8sNamespace := &v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name:         n.namespace,
				GenerateName: n.namespace,
			},
		}
		ns, err = n.client.CoreV1().Namespaces().Create(ctx, k8sNamespace, metav1.CreateOptions{})
	case deleteNamespace:
		err = n.client.CoreV1().Namespaces().Delete(ctx, n.namespace, metav1.DeleteOptions{})
	}

	return ns, err
}

func (n namespaceManager) ShouldRetry(tries int, err error) bool {
	return tries < n.maxAttempts && err != nil
}

type featureFlagTest func(t *testing.T, flagName string, flagValue bool)

func TestRunIntegrationTestsWithFeatureFlag(t *testing.T) {
	t.Parallel()

	tests := map[string]featureFlagTest{
		"testKubernetesSuccessRun":                                testKubernetesSuccessRunFeatureFlag,
		"testKubernetesMultistepRun":                              testKubernetesMultistepRunFeatureFlag,
		"testKubernetesTimeoutRun":                                testKubernetesTimeoutRunFeatureFlag,
		"testKubernetesBuildFail":                                 testKubernetesBuildFailFeatureFlag,
		"testKubernetesBuildCancel":                               testKubernetesBuildCancelFeatureFlag,
		"testKubernetesBuildLogLimitExceeded":                     testKubernetesBuildLogLimitExceededFeatureFlag,
		"testKubernetesBuildMasking":                              testKubernetesBuildMaskingFeatureFlag,
		"testKubernetesBuildPassingEnvsMultistep":                 testKubernetesBuildPassingEnvsMultistep,
		"testKubernetesCustomClonePath":                           testKubernetesCustomClonePathFeatureFlag,
		"testKubernetesNoRootImage":                               testKubernetesNoRootImageFeatureFlag,
		"testKubernetesMissingImage":                              testKubernetesMissingImageFeatureFlag,
		"testKubernetesMissingTag":                                testKubernetesMissingTagFeatureFlag,
		"testKubernetesFailingToPullImageTwiceFeatureFlag":        testKubernetesFailingToPullImageTwiceFeatureFlag,
		"testKubernetesFailingToPullServiceImageTwiceFeatureFlag": testKubernetesFailingToPullSvcImageTwiceFeatureFlag,
		"testKubernetesFailingToPullHelperTwiceFeatureFlag":       testKubernetesFailingToPullHelperTwiceFeatureFlag,
		"testOverwriteNamespaceNotMatch":                          testOverwriteNamespaceNotMatchFeatureFlag,
		"testOverwriteServiceAccountNotMatch":                     testOverwriteServiceAccountNotMatchFeatureFlag,
		"testInteractiveTerminal":                                 testInteractiveTerminalFeatureFlag,
		"testKubernetesReplaceEnvFeatureFlag":                     testKubernetesReplaceEnvFeatureFlag,
		"testKubernetesReplaceMissingEnvVarFeatureFlag":           testKubernetesReplaceMissingEnvVarFeatureFlag,
		"testKubernetesWithNonRootSecurityContext":                testKubernetesWithNonRootSecurityContext,
		"testBuildsDirDefaultVolumeFeatureFlag":                   testBuildsDirDefaultVolumeFeatureFlag,
		"testBuildsDirVolumeMountEmptyDirFeatureFlag":             testBuildsDirVolumeMountEmptyDirFeatureFlag,
		"testBuildsDirVolumeMountHostPathFeatureFlag":             testBuildsDirVolumeMountHostPathFeatureFlag,
		"testKubernetesBashFeatureFlag":                           testKubernetesBashFeatureFlag,
		"testKubernetesContainerHookFeatureFlag":                  testKubernetesContainerHookFeatureFlag,
		"testKubernetesGarbageCollection":                         testKubernetesGarbageCollection,
		"testKubernetesNamespaceIsolation":                        testKubernetesNamespaceIsolation,
		"testKubernetesPublicInternalVariables":                   testKubernetesPublicInternalVariables,
		"testKubernetesWaitResources":                             testKubernetesWaitResources,
		"testKubernetesLongLogsFeatureFlag":                       testKubernetesLongLogsFeatureFlag,
		"testKubernetesHugeScriptAndAfterScriptFeatureFlag":       testKubernetesHugeScriptAndAfterScriptFeatureFlag,
		"testKubernetesCustomPodSpec":                             testKubernetesCustomPodSpec,
		"testKubernetesClusterWarningEvent":                       testKubernetesClusterWarningEvent,
		"testKubernetesFailingBuildForBashAndPwshFeatureFlag":     testKubernetesFailingBuildForBashAndPwshFeatureFlag,
		"testKubernetesPodEvents":                                 testKubernetesPodEvents,
		"testKubernetesDumbInitSuccessRun":                        testKubernetesDumbInitSuccessRun,
		"testKubernetesDisableUmask":                              testKubernetesDisableUmask,
		"testKubernetesNoAdditionalNewLines":                      testKubernetesNoAdditionalNewLines,
		"testJobRunningAndPassingWhenServiceStops":                testJobRunningAndPassingWhenServiceStops,
		"testJobAgainstServiceContainerBehaviour":                 testJobAgainstServiceContainerBehaviour,
	}

	ff := featureflags.UseLegacyKubernetesExecutionStrategy
	ffValue := false
	if testFeatureFlag != "" {
		ff = testFeatureFlag
		ffValue = testFeatureFlagValue
	}

	for name, testFunc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			toggleText := "off"
			if ffValue {
				toggleText = "on"
			}

			t.Run(ff+":"+toggleText, func(t *testing.T) {
				t.Parallel()
				testFunc(t, ff, ffValue)
			})
		})
	}
}

func testKubernetesSuccessRunFeatureFlag(t *testing.T, featureFlagName string, featureFlagValue bool) {
	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	build := getTestBuild(t, common.GetRemoteSuccessfulBuild)
	build.Image.Name = common.TestDockerGitImage
	buildtest.SetBuildFeatureFlag(build, featureFlagName, featureFlagValue)

	err := build.Run(&common.Config{}, &common.Trace{Writer: os.Stdout})
	assert.NoError(t, err)
}

func testKubernetesPodEvents(t *testing.T, featureFlagName string, featureFlagValue bool) {
	t.Skip("TODO: Fix events not properly tested for or waited for - https://gitlab.com/gitlab-org/gitlab-runner/-/jobs/8532408889")
	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	build := getTestBuild(t, common.GetRemoteSuccessfulBuild)
	buildtest.SetBuildFeatureFlag(build, featureFlagName, featureFlagValue)
	build.Image.Name = common.TestAlpineImage
	build.Variables = append(
		build.Variables,
		common.JobVariable{Key: featureflags.PrintPodEvents, Value: "true"},
	)

	out, err := buildtest.RunBuildReturningOutput(t, build)
	require.NoError(t, err)

	expectedLines := []string{
		"Type     Reason      Message",
		"Normal   Scheduled   Successfully assigned",
	}

	if build.Variables.Get(featureflags.UseLegacyKubernetesExecutionStrategy) == "false" {
		expectedLines = append(
			expectedLines,
			"Normal   Created   Created container init-permissions",
			"Normal   Started   Started container init-permissions",
		)
	}

	expectedLines = append(
		expectedLines,
		"Normal   Pulling   Pulling image|Normal   Pulled   Successfully pulled image|Normal   Pulled   Container image .* already present on machine",
		"Normal   Created   Created container build",
		"Normal   Started   Started container build",
		"Normal   Created   Created container helper",
		"Normal   Started   Started container helper",
	)

	for _, l := range expectedLines {
		assert.Regexp(t, regexp.MustCompile(fmt.Sprintf(`(?m)%s`, l)), out)
	}
}

func testKubernetesBuildPassingEnvsMultistep(t *testing.T, featureFlagName string, featureFlagValue bool) {
	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	shellstest.OnEachShell(t, func(t *testing.T, shell string) {
		build := getTestBuild(t, func() (common.JobResponse, error) {
			return common.JobResponse{}, nil
		})
		build.Runner.RunnerSettings.Shell = shell
		if shell == "pwsh" {
			t.Skip("TODO: Fix pwsh fails")
		}

		buildtest.RunBuildWithPassingEnvsMultistep(
			t,
			build.Runner,
			func(_ *testing.T, build *common.Build) {
				buildtest.SetBuildFeatureFlag(build, featureFlagName, featureFlagValue)
			},
		)
	})
}

func testKubernetesDumbInitSuccessRun(t *testing.T, featureFlagName string, featureFlagValue bool) {
	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	build := getTestBuild(t, common.GetRemoteSuccessfulBuild)
	build.Image.Name = common.TestDockerGitImage
	buildtest.SetBuildFeatureFlag(build, featureFlagName, featureFlagValue)
	buildtest.SetBuildFeatureFlag(build, featureflags.UseDumbInitWithKubernetesExecutor, true)

	err := build.Run(&common.Config{}, &common.Trace{Writer: os.Stdout})
	assert.NoError(t, err)
}

func testKubernetesDisableUmask(t *testing.T, featureFlagName string, featureFlagValue bool) {
	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	customBuildDir := "/custom_builds_dir"
	tests := map[string]struct {
		image        string
		shell        string
		buildDir     string
		script       string
		runAsUser    int64
		runAsGroup   int64
		disableUmask bool
		envars       common.JobVariables
		verifyFn     func(t *testing.T, out string)
	}{
		"umask enabled": {
			image:      common.TestAlpineImage,
			shell:      "bash",
			script:     "ls -lR /builds/gitlab-org/ci-cd/gitlab-runner-pipeline-tests",
			runAsUser:  int64(1234),
			runAsGroup: int64(5678),
			verifyFn: func(t *testing.T, out string) {
				assert.NotContains(t, out, "1234")
				assert.NotContains(t, out, "5678")
				assert.Regexp(t, regexp.MustCompile(`(?m)^.*root\s*root.*gitlab-test.*$`), out)
			},
		},
		"umask disabled": {
			image:        common.TestAlpineImage,
			shell:        "bash",
			script:       "ls -lR /builds/gitlab-org/ci-cd/gitlab-runner-pipeline-tests",
			runAsUser:    int64(1234),
			runAsGroup:   int64(5678),
			disableUmask: true,
			verifyFn: func(t *testing.T, out string) {
				assert.NotContains(t, out, "root")
				assert.Regexp(t, regexp.MustCompile(`(?m)^.*1234\s*5678.*gitlab-test.*$`), out)
			},
		},
		"umask enabled with custom builds_dir": {
			image:      common.TestAlpineImage,
			shell:      "bash",
			buildDir:   customBuildDir,
			script:     "ls -lR $BUILDS_DIRECTORY/gitlab-org/ci-cd/gitlab-runner-pipeline-tests",
			runAsUser:  int64(1234),
			runAsGroup: int64(5678),
			envars: common.JobVariables{
				common.JobVariable{Key: "BUILDS_DIRECTORY", Value: customBuildDir},
			},
			verifyFn: func(t *testing.T, out string) {
				assert.NotContains(t, out, "1234")
				assert.NotContains(t, out, "5678")
				assert.Regexp(t, regexp.MustCompile(`(?m)^.*root\s*root.*gitlab-test.*$`), out)
			},
		},
		"umask disabled with custom builds_dir": {
			image:        common.TestAlpineImage,
			shell:        "bash",
			buildDir:     customBuildDir,
			script:       "ls -lR $BUILDS_DIRECTORY/gitlab-org/ci-cd/gitlab-runner-pipeline-tests",
			runAsUser:    int64(1234),
			runAsGroup:   int64(5678),
			disableUmask: true,
			envars: common.JobVariables{
				common.JobVariable{Key: "BUILDS_DIRECTORY", Value: customBuildDir},
			},
			verifyFn: func(t *testing.T, out string) {
				assert.NotContains(t, out, "root")
				assert.Regexp(t, regexp.MustCompile(`(?m)^.*1234\s*5678.*gitlab-test.*$`), out)
			},
		},
	}

	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			t.Parallel()

			build := getTestBuild(t, func() (common.JobResponse, error) {
				return common.GetRemoteBuildResponse(tc.script)
			})

			build.Variables = append(build.Variables, tc.envars...)
			build.Runner.RunnerSettings.Shell = tc.shell
			build.JobResponse.Image.Name = tc.image

			if tc.buildDir != "" {
				build.Runner.BuildsDir = tc.buildDir
				build.Runner.Kubernetes.Volumes = common.KubernetesVolumes{
					EmptyDirs: []common.KubernetesEmptyDir{
						common.KubernetesEmptyDir{
							Name:      "repo",
							MountPath: "$BUILDS_DIRECTORY",
						},
					},
				}
			}

			build.Runner.Kubernetes.BuildContainerSecurityContext = common.KubernetesContainerSecurityContext{
				RunAsUser:  &tc.runAsUser,
				RunAsGroup: &tc.runAsGroup,
			}
			build.Runner.Kubernetes.HelperImage = "registry.gitlab.com/gitlab-org/gitlab-runner/gitlab-runner-helper:x86_64-latest"

			buildtest.SetBuildFeatureFlag(build, featureFlagName, featureFlagValue)
			buildtest.SetBuildFeatureFlag(build, "FF_DISABLE_UMASK_FOR_KUBERNETES_EXECUTOR", tc.disableUmask)

			var buf bytes.Buffer
			err := build.Run(&common.Config{}, &common.Trace{Writer: &buf})
			assert.NoError(t, err)

			tc.verifyFn(t, buf.String())
		})
	}
}

func testKubernetesNoAdditionalNewLines(t *testing.T, featureFlagName string, featureFlagValue bool) {
	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	build := getTestBuild(t, func() (common.JobResponse, error) {
		return common.GetRemoteBuildResponse("for i in $(seq 1 120); do printf .; sleep 0.02; done; echo")
	})

	build.Runner.RunnerSettings.Shell = "bash"
	build.JobResponse.Image.Name = common.TestAlpineImage
	build.Runner.Kubernetes.HelperImage = "registry.gitlab.com/gitlab-org/gitlab-runner/gitlab-runner-helper:x86_64-latest"

	buildtest.SetBuildFeatureFlag(build, featureFlagName, featureFlagValue)

	var buf bytes.Buffer
	err := build.Run(&common.Config{}, &common.Trace{Writer: &buf})
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "........................................................................................................................")
}

func TestBuildScriptSections(t *testing.T) {
	t.Parallel()

	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	shellstest.OnEachShell(t, func(t *testing.T, shell string) {
		if shell != "bash" {
			t.Skip("TODO: fix this test for non-bash shells. This wasn't working before anyways because the image was never set correctly.")
		}

		t.Parallel()
		build := getTestBuild(t, func() (common.JobResponse, error) {
			return common.GetRemoteBuildResponse(`echo "Hello
World"`)
		})
		build.Runner.RunnerSettings.Shell = shell
		if shell != "bash" {
			build.Runner.Kubernetes.Image = common.TestPwshImage
		}

		buildtest.RunBuildWithSections(t, build)
	})
}

func TestEntrypointNotIgnored(t *testing.T) {
	t.Parallel()

	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	buildTestJob := func() (common.JobResponse, error) {
		return common.GetRemoteBuildResponse(
			"if [ -f /tmp/debug.log ]; then",
			"cat /tmp/debug.log",
			"else",
			"echo 'file not found'",
			"fi",
			"echo \"I am now `whoami`\"",
		)
	}

	helperTestJob := func() (common.JobResponse, error) {
		return common.GetRemoteBuildResponse(
			"if [ -f /builds/debug.log ]; then",
			"cat /builds/debug.log",
			"else",
			"echo 'file not found'",
			"fi",
			"echo \"I am now `whoami`\"",
		)
	}

	testCases := map[string]struct {
		jobResponse          func() (common.JobResponse, error)
		buildImage           string
		helperImage          string
		useHonorEntrypointFF bool
		expectedOutputLines  []string
	}{
		"build image with entrypoint feature flag off": {
			jobResponse:          buildTestJob,
			buildImage:           common.TestAlpineEntrypointImage,
			useHonorEntrypointFF: false,
			expectedOutputLines:  []string{"I am now root", "file not found"},
		},
		"build image with entrypoint feature flag on": {
			jobResponse:          buildTestJob,
			buildImage:           common.TestAlpineEntrypointImage,
			useHonorEntrypointFF: true,
			expectedOutputLines:  []string{"I am now nobody", "this has been executed through a custom entrypoint"},
		},
		"helper image with entrypoint feature flag off": {
			jobResponse:          helperTestJob,
			helperImage:          common.TestHelperEntrypointImage,
			useHonorEntrypointFF: false,
			expectedOutputLines:  []string{"I am now root", "file not found"},
		},
		"helper image with entrypoint feature flag on": {
			jobResponse:          helperTestJob,
			helperImage:          common.TestHelperEntrypointImage,
			useHonorEntrypointFF: true,
			expectedOutputLines:  []string{"I am now nobody", "this has been executed through a custom entrypoint"},
		},
	}

	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			t.Parallel()

			build := getTestBuildWithImage(t, common.TestAlpineEntrypointImage, func() (common.JobResponse, error) {
				jobResponse, err := tc.jobResponse()
				if err != nil {
					return common.JobResponse{}, err
				}

				jobResponse.Image = common.Image{
					Name: common.TestAlpineEntrypointImage,
				}

				return jobResponse, nil
			})

			if tc.helperImage != "" {
				build.Runner.Kubernetes.HelperImage = common.TestHelperEntrypointImage
			}

			build.Variables = append(
				build.Variables,
				common.JobVariable{Key: featureflags.KubernetesHonorEntrypoint, Value: strconv.FormatBool(tc.useHonorEntrypointFF)},
			)

			out, err := buildtest.RunBuildReturningOutput(t, build)
			require.NoError(t, err)

			t.Log(out)

			for _, expectedLine := range tc.expectedOutputLines {
				assert.Contains(t, out, expectedLine)
			}
		})
	}
}

func testKubernetesMultistepRunFeatureFlag(t *testing.T, featureFlagName string, featureFlagValue bool) {
	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	successfulBuild, err := common.GetRemoteSuccessfulMultistepBuild()
	require.NoError(t, err)

	failingScriptBuild, err := common.GetRemoteFailingMultistepBuild(common.StepNameScript)
	require.NoError(t, err)

	failingReleaseBuild, err := common.GetRemoteFailingMultistepBuild("release")
	require.NoError(t, err)

	successfulBuild.Image.Name = common.TestDockerGitImage
	failingScriptBuild.Image.Name = common.TestDockerGitImage
	failingReleaseBuild.Image.Name = common.TestDockerGitImage

	tests := map[string]struct {
		jobResponse    common.JobResponse
		expectedOutput []string
		unwantedOutput []string
		errExpected    bool
	}{
		"Successful build with release and after_script step": {
			jobResponse: successfulBuild,
			expectedOutput: []string{
				"echo Hello World",
				"echo Release",
				"echo After Script",
			},
			errExpected: false,
		},
		"Failure on script step. Release is skipped. After script runs.": {
			jobResponse: failingScriptBuild,
			expectedOutput: []string{
				"echo Hello World",
				"echo After Script",
			},
			unwantedOutput: []string{
				"echo Release",
			},
			errExpected: true,
		},
		"Failure on release step. After script runs.": {
			jobResponse: failingReleaseBuild,
			expectedOutput: []string{
				"echo Hello World",
				"echo Release",
				"echo After Script",
			},
			errExpected: true,
		},
	}

	for tn, tt := range tests {
		t.Run(tn, func(t *testing.T) {
			t.Parallel()

			build := getTestBuild(t, func() (common.JobResponse, error) {
				return tt.jobResponse, nil
			})
			buildtest.SetBuildFeatureFlag(build, featureFlagName, featureFlagValue)

			var buf bytes.Buffer
			err := build.Run(&common.Config{}, &common.Trace{Writer: &buf})

			out := buf.String()
			for _, output := range tt.expectedOutput {
				assert.Contains(t, out, output)
			}

			for _, output := range tt.unwantedOutput {
				assert.NotContains(t, out, output)
			}

			if tt.errExpected {
				var buildErr *common.BuildError
				assert.ErrorAs(t, err, &buildErr)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func testKubernetesTimeoutRunFeatureFlag(t *testing.T, featureFlagName string, featureFlagValue bool) {
	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	build := getTestBuild(t, common.GetRemoteLongRunningBuild)
	build.Image.Name = common.TestDockerGitImage
	build.RunnerInfo.Timeout = 10 // seconds
	buildtest.SetBuildFeatureFlag(build, featureFlagName, featureFlagValue)

	err := build.Run(&common.Config{}, &common.Trace{Writer: os.Stdout})
	require.Error(t, err)
	var buildError *common.BuildError
	assert.ErrorAs(t, err, &buildError)
	assert.Equal(t, common.JobExecutionTimeout, buildError.FailureReason)
}

func testKubernetesLongLogsFeatureFlag(t *testing.T, featureFlagName string, featureFlagValue bool) {
	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	tests := map[string]struct {
		getLine func() string
	}{
		"short log": {
			getLine: func() string {
				return "Regular log"
			},
		},
		"buffer size log": {
			getLine: func() string {
				return strings.Repeat("1", common.DefaultReaderBufferSize)
			},
		},
		"long log": {
			getLine: func() string {
				return strings.Repeat("lorem ipsum", common.DefaultReaderBufferSize)
			},
		},
	}

	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			t.Parallel()

			line := tc.getLine()
			build := getTestBuild(t, func() (common.JobResponse, error) {
				return common.GetRemoteBuildResponse(fmt.Sprintf(`echo "%s"`, line))
			})
			buildtest.SetBuildFeatureFlag(build, featureFlagName, featureFlagValue)

			outBuffer := new(bytes.Buffer)
			err := build.Run(&common.Config{}, &common.Trace{Writer: outBuffer})
			require.NoError(t, err)
			assert.Contains(t, outBuffer.String(), fmt.Sprintf(`$ echo "%s"`, line))
			// We check the whole line is found in the log without any newline within
			assert.Regexp(t, regexp.MustCompile(fmt.Sprintf(`(?m)^%s$`, line)), outBuffer.String())
		})
	}
}

func testKubernetesHugeScriptAndAfterScriptFeatureFlag(t *testing.T, featureFlagName string, featureFlagValue bool) {
	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	getAfterScript := func(featureFlag bool, script ...string) common.Step {
		as := common.Step{
			Name: "after_script",
			Script: common.StepScript{
				"echo $CI_JOB_STATUS",
			},
			Timeout:      3600,
			When:         common.StepWhenAlways,
			AllowFailure: true,
		}

		if !featureFlag {
			as.Script = append(as.Script, "ls -l /scripts-0-0/*")
		}

		as.Script = append(as.Script, script...)

		return as
	}

	tests := map[string]struct {
		image       string
		shell       string
		getScript   func() common.StepScript
		afterScript []string
		verifyFn    func(t *testing.T, out string)
	}{
		"bash normal script": {
			image: common.TestAlpineImage,
			shell: "bash",
			getScript: func() common.StepScript {
				return []string{
					`echo "My normal string"`,
				}
			},
			verifyFn: func(t *testing.T, out string) {
				assert.Contains(t, out, "success")
				assert.Contains(t, out, "My normal string")
			},
		},
		"pwsh unicode script": {
			image: common.TestPwshImage,
			shell: "pwsh",
			getScript: func() common.StepScript {
				return []string{
					"echo \"`“ `“ `” `” `„ ‘ ’ ‚ ‛ ‘ ’ ; < ( ) & ^ # [ ] { } ' < > | @ % „",
				}
			},
			verifyFn: func(t *testing.T, out string) {
				assert.Contains(t, out, "success")
				assert.Contains(t, out, "“ “ ” ” „ ‘ ’ ‚ ‛ ‘ ’ ; < ( ) & ^ # [ ] { } ' < > | @ %")
			},
		},
		"bash nested here string": {
			image:       common.TestAlpineImage,
			shell:       "bash",
			afterScript: []string{"cat ./print.sh"},
			getScript: func() common.StepScript {
				return []string{
					`cat <<EOF > ./print.sh
#!/bin/bash
echo "My nested here-string"
EOF
`,
				}
			},
			verifyFn: func(t *testing.T, out string) {
				assert.Contains(t, out, "success")
				assert.Contains(t, out, "echo \"My nested here-string\"")
			},
		},
		"pwsh nested here-string": {
			image: common.TestPwshImage,
			shell: "pwsh",
			getScript: func() common.StepScript {
				return []string{
					`echo @'
My nested here-string
echo @"
My nested nested here-string
“ “ ” ” „ ‘ ’ ‚ ‛ ‘ ’ ; < ( ) & ^ # [ ] { } ' < > | @ %
"@
'@`,
				}
			},
			verifyFn: func(t *testing.T, out string) {
				assert.Contains(t, out, "success")
				assert.Contains(t, out, "My nested here-string")
				assert.Contains(t, out, "My nested nested here-string")
				assert.Contains(t, out, "“ “ ” ” „ ‘ ’ ‚ ‛ ‘ ’ ; < ( ) & ^ # [ ] { } ' < > | @ %")
			},
		},
		"bash huge script": {
			image: common.TestAlpineImage,
			shell: "bash",
			getScript: func() common.StepScript {
				s := strings.Repeat(
					"echo \"Lorem ipsum dolor sit amet, consectetur adipiscing elit\"\n",
					10*1024,
				)
				return strings.Split(s, "\n")
			},
			verifyFn: func(t *testing.T, out string) {
				assert.Contains(t, out, "success")
				assert.Contains(t, out, "Lorem ipsum dolor sit amet, consectetur adipiscing elit")
			},
		},
		"pwsh special script with special character": {
			image: common.TestPwshImage,
			shell: "pwsh",
			getScript: func() common.StepScript {
				return []string{
					`& {$Calendar = Get-Date; If ($Calendar.Month -eq '0') {"This is wrong"} Else {echo "not happening" > test.txt}; ls; Get-Content test.txt;}`,
				}
			},
			verifyFn: func(t *testing.T, out string) {
				assert.Contains(t, out, "test.txt")
				assert.Contains(t, out, "not happening")
			},
		},
		"pwsh multiple instructions in the script": {
			image: common.TestPwshImage,
			shell: "pwsh",
			getScript: func() common.StepScript {
				return []string{
					`$Calendar = Get-Date`,
					`If ($Calendar.Month -eq '0') {"This is wrong"} Else {echo "not happening" > test.txt}`,
					`ls`,
					`Get-Content test.txt`,
					`&{ echo "Display special characters () {} <> [] \ | ;"}`,
				}
			},
			verifyFn: func(t *testing.T, out string) {
				assert.Contains(t, out, "test.txt")
				assert.Contains(t, out, "not happening")
				assert.Contains(t, out, "Display special characters () {} <> [] \\ | ;")
			},
		},
		"pwsh instruction with arrays": {
			image: common.TestPwshImage,
			shell: "pwsh",
			getScript: func() common.StepScript {
				return []string{
					`$data = @('Zero','One')`,
					`$data | % {"$PSItem"}`,
					`$data = @{two = "Two"; three = "Three"; }`,
					`$data.values | % {"$PSItem"}`,
				}
			},
			verifyFn: func(t *testing.T, out string) {
				assert.Contains(t, out, "Zero")
				assert.Contains(t, out, "One")
				assert.Contains(t, out, "Two")
				assert.Contains(t, out, "Three")
			},
		},
	}

	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			t.Parallel()

			build := getTestBuild(t, func() (common.JobResponse, error) {
				return common.GetRemoteBuildResponse("echo \"Hello World\"")
			})

			build.Runner.RunnerSettings.Shell = tc.shell
			build.JobResponse.Image.Name = tc.image
			build.JobResponse.Steps[0].Script = append(
				build.JobResponse.Steps[0].Script,
				tc.getScript()...,
			)
			build.JobResponse.Steps = append(
				build.JobResponse.Steps,
				getAfterScript(featureFlagValue, tc.afterScript...),
			)

			buildtest.SetBuildFeatureFlag(build, featureFlagName, featureFlagValue)

			outBuffer := new(bytes.Buffer)
			err := build.Run(&common.Config{}, &common.Trace{Writer: outBuffer})
			require.NoError(t, err)

			if !featureFlagValue {
				assert.Contains(t, outBuffer.String(), "echo $CI_JOB_STATUS")
				assert.Contains(t, outBuffer.String(), "/scripts-0-0/step_script")
				assert.Contains(t, outBuffer.String(), "/scripts-0-0/after_script")
			}

			if tc.verifyFn != nil {
				tc.verifyFn(t, outBuffer.String())
			}
		})
	}
}

func testKubernetesCustomPodSpec(t *testing.T, featureFlagName string, featureFlagValue bool) {
	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	ctxTimeout := time.Minute
	client := getTestKubeClusterClient(t)

	init := func(t *testing.T, _ *common.Build, client *k8s.Clientset) {
		credentials, err := getSecrets(client, ciNamespace, "")
		require.NoError(t, err)
		configMaps, err := getConfigMaps(client, ciNamespace, "")
		require.NoError(t, err)

		assert.Empty(t, credentials)
		assert.Empty(t, configMaps)
	}

	tests := map[string]struct {
		podSpec  []common.KubernetesPodSpec
		verifyFn func(*testing.T, v1.Pod)
	}{
		"change hostname with custom podSpec": {
			podSpec: []common.KubernetesPodSpec{
				{
					Patch: `
[
	{
		"op": "add",
		"path": "/hostname",
		"value": "my-custom-hostname"
	}
]
`,
					PatchType: common.PatchTypeJSONPatchType,
				},
			},
			verifyFn: func(t *testing.T, pod v1.Pod) {
				assert.Equal(t, "my-custom-hostname", pod.Spec.Hostname)
			},
		},
		"update build container with resources limit through custom podSpec using strategic patch type": {
			podSpec: []common.KubernetesPodSpec{
				{
					Patch: `
containers:
- name: "build"
  securityContext:
    runAsUser: 1010
`,
					PatchType: common.PatchTypeStrategicMergePatchType,
				},
			},
			verifyFn: func(t *testing.T, pod v1.Pod) {
				var buildContainer v1.Container
				var user int64 = 1010

				for _, c := range pod.Spec.Containers {
					if c.Name == "build" {
						buildContainer = c
						break
					}
				}
				assert.Equal(t, user, *buildContainer.SecurityContext.RunAsUser)
			},
		},
	}

	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			t.Parallel()

			build := getTestBuild(t, func() (common.JobResponse, error) {
				jobResponse, err := common.GetRemoteBuildResponse(
					"sleep 5000",
				)
				require.NoError(t, err)

				return jobResponse, nil
			})
			build.Runner.Kubernetes.PodSpec = tc.podSpec
			buildtest.SetBuildFeatureFlag(build, featureFlagName, featureFlagValue)
			buildtest.SetBuildFeatureFlag(build, featureflags.UseAdvancedPodSpecConfiguration, true)

			init(t, build, client)

			deletedPodNameCh := make(chan string)
			defer buildtest.OnUserStage(build, func() {
				ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
				defer cancel()
				pods, err := client.CoreV1().Pods(ciNamespace).List(
					ctx,
					metav1.ListOptions{
						LabelSelector: labels.Set(build.Runner.Kubernetes.PodLabels).String(),
					},
				)
				require.NoError(t, err)
				require.NotEmpty(t, pods.Items)
				pod := pods.Items[0]

				tc.verifyFn(t, pod)

				err = client.
					CoreV1().
					Pods(ciNamespace).
					Delete(ctx, pod.Name, metav1.DeleteOptions{
						PropagationPolicy: &kubernetes.PropagationPolicy,
					})
				require.NoError(t, err)

				deletedPodNameCh <- pod.Name
			})()

			err := build.Run(&common.Config{}, &common.Trace{Writer: os.Stdout})
			assert.Error(t, err)

			<-deletedPodNameCh
		})
	}
}

func testKubernetesFailingBuildForBashAndPwshFeatureFlag(t *testing.T, featureFlagName string, featureFlagValue bool) {
	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	tests := map[string]struct {
		image string
		shell string
	}{
		"bash failing script": {
			image: common.TestAlpineImage,
			shell: "bash",
		},
		"pwsh failing script": {
			image: common.TestPwshImage,
			shell: "pwsh",
		},
	}

	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			t.Parallel()
			if tc.shell == "pwsh" {
				t.Skip("TODO: Fix pwsh fails")
			}

			build := getTestBuild(t, func() (common.JobResponse, error) {
				return common.GetRemoteBuildResponse("invalid_command")
			})

			build.Runner.RunnerSettings.Shell = tc.shell
			build.JobResponse.Image.Name = tc.image

			buildtest.SetBuildFeatureFlag(build, featureFlagName, featureFlagValue)
			err := build.Run(&common.Config{}, &common.Trace{Writer: os.Stdout})
			require.Error(t, err)
		})
	}
}

func testKubernetesBuildFailFeatureFlag(t *testing.T, featureFlagName string, featureFlagValue bool) {
	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	build := getTestBuild(t, common.GetRemoteFailedBuild)
	build.Image.Name = common.TestDockerGitImage
	buildtest.SetBuildFeatureFlag(build, featureFlagName, featureFlagValue)

	err := build.Run(&common.Config{}, &common.Trace{Writer: os.Stdout})
	require.Error(t, err, "error")
	var buildError *common.BuildError
	require.ErrorAs(t, err, &buildError)
	assert.Contains(t, err.Error(), "command terminated with exit code 1")
	assert.Equal(t, 1, buildError.ExitCode)
}

func testKubernetesBuildCancelFeatureFlag(t *testing.T, featureFlagName string, featureFlagValue bool) {
	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	build := getTestBuild(t, func() (common.JobResponse, error) {
		return common.JobResponse{}, nil
	})
	buildtest.RunBuildWithCancel(
		t,
		build.Runner,
		func(_ *testing.T, build *common.Build) {
			buildtest.SetBuildFeatureFlag(build, featureFlagName, featureFlagValue)
		},
	)
}

func testKubernetesBuildLogLimitExceededFeatureFlag(t *testing.T, featureFlagName string, featureFlagValue bool) {
	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	build := getTestBuild(t, func() (common.JobResponse, error) {
		return common.JobResponse{}, nil
	})
	buildtest.RunRemoteBuildWithJobOutputLimitExceeded(
		t,
		build.Runner,
		func(_ *testing.T, build *common.Build) {
			buildtest.SetBuildFeatureFlag(build, featureFlagName, featureFlagValue)
		},
	)
}

func testKubernetesBuildMaskingFeatureFlag(t *testing.T, featureFlagName string, featureFlagValue bool) {
	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	build := getTestBuild(t, func() (common.JobResponse, error) {
		return common.JobResponse{}, nil
	})
	buildtest.RunBuildWithMasking(
		t,
		build.Runner,
		func(_ *testing.T, build *common.Build) {
			buildtest.SetBuildFeatureFlag(build, featureFlagName, featureFlagValue)
		},
	)
}

func testKubernetesCustomClonePathFeatureFlag(t *testing.T, featureFlagName string, featureFlagValue bool) {
	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	tests := map[string]struct {
		clonePath   string
		expectedErr bool
	}{
		"uses custom clone path": {
			clonePath:   "$CI_BUILDS_DIR/go/src/gitlab.com/gitlab-org/repo",
			expectedErr: false,
		},
		"path has to be within CI_BUILDS_DIR": {
			clonePath:   "/unknown/go/src/gitlab.com/gitlab-org/repo",
			expectedErr: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			build := getTestBuild(t, func() (common.JobResponse, error) {
				return common.GetRemoteBuildResponse(
					"ls -al $CI_BUILDS_DIR/go/src/gitlab.com/gitlab-org/repo",
				)
			})
			build.Runner.Environment = []string{
				"GIT_CLONE_PATH=" + test.clonePath,
			}
			buildtest.SetBuildFeatureFlag(build, featureFlagName, featureFlagValue)

			err := build.Run(&common.Config{}, &common.Trace{Writer: os.Stdout})
			if test.expectedErr {
				var buildErr *common.BuildError
				assert.ErrorAs(t, err, &buildErr)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func testKubernetesNoRootImageFeatureFlag(t *testing.T, featureFlagName string, featureFlagValue bool) {
	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	build := getTestBuild(t, common.GetRemoteSuccessfulBuildWithDumpedVariables)
	build.Image.Name = common.TestAlpineNoRootImage
	buildtest.SetBuildFeatureFlag(build, featureFlagName, featureFlagValue)

	err := build.Run(&common.Config{}, &common.Trace{Writer: os.Stdout})
	assert.NoError(t, err)
}

func testKubernetesMissingImageFeatureFlag(t *testing.T, featureFlagName string, featureFlagValue bool) {
	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	build := getTestBuild(t, common.GetRemoteFailedBuild)
	build.Image.Name = "some/non-existing/image"
	buildtest.SetBuildFeatureFlag(build, featureFlagName, featureFlagValue)

	err := build.Run(&common.Config{}, &common.Trace{Writer: os.Stdout})
	require.Error(t, err)
	assert.ErrorIs(t, err, &common.BuildError{FailureReason: common.ImagePullFailure})
	assert.Contains(t, err.Error(), "image pull failed")
}

func testKubernetesMissingTagFeatureFlag(t *testing.T, featureFlagName string, featureFlagValue bool) {
	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	build := getTestBuild(t, common.GetRemoteFailedBuild)
	build.Image.Name = "docker:missing-tag"
	buildtest.SetBuildFeatureFlag(build, featureFlagName, featureFlagValue)

	err := build.Run(&common.Config{}, &common.Trace{Writer: os.Stdout})
	require.Error(t, err)
	assert.ErrorIs(t, err, &common.BuildError{FailureReason: common.ImagePullFailure})
	assert.Contains(t, err.Error(), "image pull failed")
}

func testKubernetesFailingToPullImageTwiceFeatureFlag(t *testing.T, featureFlagName string, featureFlagValue bool) {
	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	build := getTestBuild(t, common.GetRemoteFailedBuild)
	build.Image.Name = "some/non-existing/image"
	buildtest.SetBuildFeatureFlag(build, featureFlagName, featureFlagValue)

	err := runMultiPullPolicyBuild(t, build)

	var imagePullErr *pull.ImagePullError
	require.ErrorAs(t, err, &imagePullErr)
	assert.Equal(t, build.Image.Name, imagePullErr.Image)
}

func testKubernetesFailingToPullSvcImageTwiceFeatureFlag(t *testing.T, featureFlagName string, featureFlagValue bool) {
	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	build := getTestBuild(t, common.GetRemoteFailedBuild)
	build.Services = common.Services{
		{
			Name: "some/non-existing/image",
		},
	}
	buildtest.SetBuildFeatureFlag(build, featureFlagName, featureFlagValue)

	err := runMultiPullPolicyBuild(t, build)

	var imagePullErr *pull.ImagePullError
	require.ErrorAs(t, err, &imagePullErr)
	assert.Equal(t, build.Services[0].Name, imagePullErr.Image)
}

func testKubernetesFailingToPullHelperTwiceFeatureFlag(t *testing.T, featureFlagName string, featureFlagValue bool) {
	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	build := getTestBuild(t, common.GetRemoteFailedBuild)
	build.Runner.Kubernetes.HelperImage = "some/non-existing/image"
	buildtest.SetBuildFeatureFlag(build, featureFlagName, featureFlagValue)

	err := runMultiPullPolicyBuild(t, build)

	var imagePullErr *pull.ImagePullError
	require.ErrorAs(t, err, &imagePullErr)
	assert.Equal(t, build.Runner.Kubernetes.HelperImage, imagePullErr.Image)
}

func testOverwriteNamespaceNotMatchFeatureFlag(t *testing.T, featureFlagName string, featureFlagValue bool) {
	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	build := getTestBuild(t, func() (common.JobResponse, error) {
		return common.JobResponse{
			GitInfo: common.GitInfo{
				Sha: "1234567890",
			},
			Image: common.Image{
				Name: "test-image",
			},
			Variables: []common.JobVariable{
				{Key: kubernetes.NamespaceOverwriteVariableName, Value: "namespace"},
			},
		}, nil
	})
	build.Runner.Kubernetes.NamespaceOverwriteAllowed = "^not_a_match$"
	build.SystemInterrupt = make(chan os.Signal, 1)
	build.Image.Name = common.TestDockerGitImage
	buildtest.SetBuildFeatureFlag(build, featureFlagName, featureFlagValue)

	err := build.Run(&common.Config{}, &common.Trace{Writer: os.Stdout})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not match")
}

func testOverwriteServiceAccountNotMatchFeatureFlag(t *testing.T, featureFlagName string, featureFlagValue bool) {
	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	build := getTestBuild(t, func() (common.JobResponse, error) {
		return common.JobResponse{
			GitInfo: common.GitInfo{
				Sha: "1234567890",
			},
			Image: common.Image{
				Name: "test-image",
			},
			Variables: []common.JobVariable{
				{Key: kubernetes.ServiceAccountOverwriteVariableName, Value: "service-account"},
			},
		}, nil
	})
	build.Runner.Kubernetes.ServiceAccountOverwriteAllowed = "^not_a_match$"
	build.SystemInterrupt = make(chan os.Signal, 1)
	build.Image.Name = common.TestDockerGitImage
	buildtest.SetBuildFeatureFlag(build, featureFlagName, featureFlagValue)

	err := build.Run(&common.Config{}, &common.Trace{Writer: os.Stdout})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not match")
}

func testInteractiveTerminalFeatureFlag(t *testing.T, featureFlagName string, featureFlagValue bool) {
	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	if os.Getenv("GITLAB_CI") == "true" {
		t.Skip("Skipping inside of GitLab CI check https://gitlab.com/gitlab-org/gitlab-runner/-/issues/26421")
	}

	client := getTestKubeClusterClient(t)
	secrets, err := client.
		CoreV1().
		Secrets(ciNamespace).
		List(context.Background(), metav1.ListOptions{})
	require.NoError(t, err)

	build := getTestBuild(t, func() (common.JobResponse, error) {
		return common.GetRemoteBuildResponse("sleep 5")
	})
	build.Image.Name = "docker:git"
	build.Runner.Kubernetes.BearerToken = string(secrets.Items[0].Data["token"])
	buildtest.SetBuildFeatureFlag(build, featureFlagName, featureFlagValue)

	sess, err := session.NewSession(nil)
	build.Session = sess

	outBuffer := bytes.NewBuffer(nil)
	outCh := make(chan string)

	go func() {
		err = build.Run(
			&common.Config{
				SessionServer: common.SessionServer{
					SessionTimeout: 2,
				},
			},
			&common.Trace{Writer: outBuffer},
		)
		require.NoError(t, err)

		outCh <- outBuffer.String()
	}()

	srv := httptest.NewServer(build.Session.Handler())
	defer srv.Close()

	u := url.URL{
		Scheme: "ws",
		Host:   srv.Listener.Addr().String(),
		Path:   build.Session.Endpoint + "/exec",
	}
	headers := http.Header{
		"Authorization": []string{build.Session.Token},
	}
	conn, resp, err := websocket.DefaultDialer.Dial(u.String(), headers)
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
		if conn != nil {
			_ = conn.Close()
		}
	}()
	require.NoError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusSwitchingProtocols)

	out := <-outCh
	t.Log(out)

	assert.Contains(t, out, "Terminal is connected, will time out in 2s...")
}

func testKubernetesReplaceEnvFeatureFlag(t *testing.T, featureFlagName string, featureFlagValue bool) {
	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")
	build := getTestBuild(t, common.GetRemoteSuccessfulBuild)
	build.Image.Name = "$IMAGE:$VERSION"
	build.JobResponse.Variables = append(
		build.JobResponse.Variables,
		common.JobVariable{Key: "IMAGE", Value: "alpine"},
		common.JobVariable{Key: "VERSION", Value: "latest"},
	)
	buildtest.SetBuildFeatureFlag(build, featureFlagName, featureFlagValue)
	out, err := buildtest.RunBuildReturningOutput(t, build)
	require.NoError(t, err)
	assert.Contains(t, out, "alpine:latest")
}

func testKubernetesReplaceMissingEnvVarFeatureFlag(t *testing.T, featureFlagName string, featureFlagValue bool) {
	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")
	build := getTestBuild(t, common.GetRemoteSuccessfulBuild)
	build.Image.Name = "alpine:$NOT_EXISTING_VARIABLE"
	buildtest.SetBuildFeatureFlag(build, featureFlagName, featureFlagValue)

	err := build.Run(&common.Config{}, &common.Trace{Writer: os.Stdout})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "image pull failed: Failed to apply default image tag \"alpine:\"")
}

func testBuildsDirDefaultVolumeFeatureFlag(t *testing.T, featureFlagName string, featureFlagValue bool) {
	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	build := getTestBuild(t, common.GetRemoteSuccessfulBuild)
	build.Image.Name = common.TestDockerGitImage
	build.Runner.BuildsDir = "/path/to/builds/dir"
	buildtest.SetBuildFeatureFlag(build, featureFlagName, featureFlagValue)

	err := build.Run(&common.Config{}, &common.Trace{Writer: os.Stdout})
	assert.NoError(t, err)

	assert.Equal(t, "/path/to/builds/dir/gitlab-org/ci-cd/gitlab-runner-pipeline-tests/gitlab-test", build.BuildDir)
}

func testBuildsDirVolumeMountEmptyDirFeatureFlag(t *testing.T, featureFlagName string, featureFlagValue bool) {
	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	tests := map[string]struct {
		emptyDir   common.KubernetesEmptyDir
		hasWarning bool
	}{
		"emptyDir with empty size": {
			emptyDir: common.KubernetesEmptyDir{
				Name:      "repo",
				MountPath: "/path/to/builds/dir",
				Medium:    "Memory",
			},
		},
		"emptyDir with untrimed empty size": {
			emptyDir: common.KubernetesEmptyDir{
				Name:      "repo",
				MountPath: "/path/to/builds/dir",
				Medium:    "Memory",
				SizeLimit: "  ",
			},
		},
		"emptyDir with valid size": {
			emptyDir: common.KubernetesEmptyDir{
				Name:      "repo",
				MountPath: "/path/to/builds/dir",
				Medium:    "Memory",
				SizeLimit: "1G",
			},
		},
		"emptyDir with invalid emptyDir": {
			emptyDir: common.KubernetesEmptyDir{
				Name:      "repo",
				MountPath: "/path/to/builds/dir",
				Medium:    "Memory",
				SizeLimit: "invalid",
			},
			hasWarning: true,
		},
	}

	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			t.Parallel()

			build := getTestBuild(t, common.GetRemoteSuccessfulBuild)
			build.Image.Name = common.TestDockerGitImage
			buildtest.SetBuildFeatureFlag(build, featureFlagName, featureFlagValue)
			build.Runner.BuildsDir = "/path/to/builds/dir"
			build.Runner.Kubernetes.Volumes = common.KubernetesVolumes{
				EmptyDirs: []common.KubernetesEmptyDir{
					tc.emptyDir,
				},
			}

			outBuffer := bytes.NewBuffer(nil)
			err := build.Run(&common.Config{}, &common.Trace{Writer: outBuffer})
			assert.NoError(t, err)

			if tc.hasWarning {
				assert.Contains(t, outBuffer.String(), "invalid limit quantity")
			}

			assert.Equal(t, "/path/to/builds/dir/gitlab-org/ci-cd/gitlab-runner-pipeline-tests/gitlab-test", build.BuildDir)
		})
	}
}

func testBuildsDirVolumeMountHostPathFeatureFlag(t *testing.T, featureFlagName string, featureFlagValue bool) {
	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	build := getTestBuild(t, common.GetRemoteSuccessfulBuild)
	build.Image.Name = common.TestDockerGitImage
	buildtest.SetBuildFeatureFlag(build, featureFlagName, featureFlagValue)
	build.Runner.Kubernetes.Volumes = common.KubernetesVolumes{
		HostPaths: []common.KubernetesHostPath{
			{
				Name:      "repo-host",
				MountPath: "/builds",
				HostPath:  "/tmp/builds",
			},
		},
	}

	err := build.Run(&common.Config{}, &common.Trace{Writer: os.Stdout})
	assert.NoError(t, err)

	allVariables := build.GetAllVariables()

	assert.Equal(t, fmt.Sprintf("/builds/%s/gitlab-org/ci-cd/gitlab-runner-pipeline-tests/gitlab-test", allVariables.Value("CI_CONCURRENT_ID")), build.BuildDir)
}

// testKubernetesGarbageCollection tests the deletion of resources via garbage collector once the owning pod is deleted
func testKubernetesGarbageCollection(t *testing.T, featureFlagName string, featureFlagValue bool) {
	t.Skip("TODO: Fix flaky test expected error not always matches https://gitlab.com/gitlab-org/gitlab-runner/-/jobs/8543529098#L226")
	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	ctxTimeout := time.Minute
	client := getTestKubeClusterClient(t)

	validateResourcesCreated := func(
		t *testing.T,
		client *k8s.Clientset,
		featureFlagValue bool,
		namespace string,
		podName string,
	) {
		credentials, err := getSecrets(client, namespace, podName)
		require.NoError(t, err)
		configMaps, err := getConfigMaps(client, namespace, podName)
		require.NoError(t, err)

		assert.NotEmpty(t, credentials)
		assert.Empty(t, configMaps)
	}

	validateResourcesDeleted := func(t *testing.T, client *k8s.Clientset, namespace string, podName string) {
		// The deletion propagation policy has been shifted to Background
		// in the MR https://gitlab.com/gitlab-org/gitlab-runner/-/merge_requests/4339
		// This means the dependant will be deleted in background and not immediately
		// A retry is needed to ensure the dependant is deleted at some point
		creds, err := retry.NewValue(retry.New(), func() ([]v1.Secret, error) {
			credentials, err := getSecrets(client, namespace, podName)
			require.NoError(t, err)

			if len(credentials) > 0 {
				return credentials, errors.New("secrets still exist")
			}

			return credentials, nil
		}).Run()
		require.NoError(t, err)

		cfgMaps, err := retry.NewValue(retry.New(), func() ([]v1.ConfigMap, error) {
			configMaps, err := getConfigMaps(client, namespace, podName)
			require.NoError(t, err)

			if len(configMaps) > 0 {
				return configMaps, errors.New("configMaps still exist")
			}

			return configMaps, nil
		}).Run()
		require.NoError(t, err)

		assert.Empty(t, creds)
		assert.Empty(t, cfgMaps)
	}

	tests := map[string]struct {
		init     func(t *testing.T, build *common.Build, client *k8s.Clientset)
		finalize func(t *testing.T, client *k8s.Clientset)
	}{
		"pod deletion during build step": {},
		"pod deletion during prepare stage in custom namespace": {
			init: func(t *testing.T, build *common.Build, client *k8s.Clientset) {
				credentials, err := getSecrets(client, ciNamespace, "")
				require.NoError(t, err)
				configMaps, err := getConfigMaps(client, ciNamespace, "")
				require.NoError(t, err)

				assert.Empty(t, credentials)
				assert.Empty(t, configMaps)
			},
		},
	}

	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			t.Parallel()

			build := getTestBuild(t, func() (common.JobResponse, error) {
				jobResponse, err := common.GetRemoteBuildResponse(
					"sleep 5000",
				)
				require.NoError(t, err)

				jobResponse.Credentials = []common.Credentials{
					{
						Type:     "registry",
						URL:      "http://example.com",
						Username: "user",
						Password: "password",
					},
				}

				return jobResponse, nil
			})
			buildtest.SetBuildFeatureFlag(build, featureFlagName, featureFlagValue)

			if tc.init != nil {
				tc.init(t, build, client)
			}

			deletedPodNameCh := make(chan string)
			defer buildtest.OnUserStage(build, func() {
				ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
				defer cancel()
				pods, err := client.CoreV1().Pods(ciNamespace).List(
					ctx,
					metav1.ListOptions{
						LabelSelector: labels.Set(build.Runner.Kubernetes.PodLabels).String(),
					},
				)
				require.NoError(t, err)
				require.NotEmpty(t, pods.Items)
				pod := pods.Items[0]

				validateResourcesCreated(t, client, featureFlagValue, ciNamespace, pod.Name)

				err = client.
					CoreV1().
					Pods(ciNamespace).
					Delete(ctx, pod.Name, metav1.DeleteOptions{
						PropagationPolicy: &kubernetes.PropagationPolicy,
					})
				require.NoError(t, err)

				deletedPodNameCh <- pod.Name
			})()

			out, err := buildtest.RunBuildReturningOutput(t, build)

			podName := <-deletedPodNameCh

			if !featureFlagValue {
				assert.True(t, kubernetes.IsKubernetesPodNotFoundError(err), "expected err NotFound, but got %T", err)
				assert.Contains(
					t,
					out,
					"ERROR: Job failed (system failure):",
				)
			} else {
				assert.Errorf(t, err, "command terminated with exit code 137")
			}
			validateResourcesDeleted(t, client, ciNamespace, podName)
		})
	}
}

func testKubernetesNamespaceIsolation(t *testing.T, featureFlagName string, featureFlagValue bool) {
	t.Skip("TODO: skipping namespace isolation test to add metadata for better cleanup")
	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	jobId := rand.Int()
	expectedNamespace := fmt.Sprintf("ci-job-%d", jobId)

	ctxTimeout := time.Minute
	client := getTestKubeClusterClient(t)

	validateNamespaceDeleted := func(t *testing.T, client *k8s.Clientset, namespace string) {
		ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
		defer cancel()

		ns, err := client.CoreV1().Namespaces().Get(
			ctx,
			namespace,
			metav1.GetOptions{},
		)

		require.NoError(t, err)
		assert.Equal(t, v1.NamespaceTerminating, ns.Status.Phase)
	}

	tests := map[string]struct {
		init     func(t *testing.T, build *common.Build, client *k8s.Clientset, namespace string)
		finalize func(t *testing.T, client *k8s.Clientset, namespace string)
	}{
		"test with default values": {
			init: func(t *testing.T, build *common.Build, client *k8s.Clientset, namespace string) {
				credentials, err := getSecrets(client, namespace, "")
				require.NoError(t, err)
				configMaps, err := getConfigMaps(client, namespace, "")
				require.NoError(t, err)

				assert.Empty(t, credentials)
				assert.Empty(t, configMaps)
			},
			finalize: func(t *testing.T, client *k8s.Clientset, namespace string) {
				_, err := newNamespaceManager(client, deleteNamespace, namespace).Run()
				require.NoError(t, err)
			},
		},
	}

	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			t.Parallel()

			build := getTestBuild(t, func() (common.JobResponse, error) {
				jobResponse, err := common.GetRemoteBuildResponse(
					"sleep 5000",
				)
				require.NoError(t, err)

				jobResponse.Credentials = []common.Credentials{
					{
						Type:     "registry",
						URL:      "http://example.com",
						Username: "user",
						Password: "password",
					},
				}

				return jobResponse, nil
			})
			build.ID = int64(jobId)
			build.Runner.Kubernetes.Namespace = "default"
			build.Runner.Kubernetes.NamespacePerJob = true
			buildtest.SetBuildFeatureFlag(build, featureFlagName, featureFlagValue)

			if tc.init != nil {
				tc.init(t, build, client, expectedNamespace)
			}

			deletedPodNameCh := make(chan string)
			defer buildtest.OnUserStage(build, func() {
				ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
				defer cancel()
				pods, err := client.CoreV1().Pods(expectedNamespace).List(
					ctx,
					metav1.ListOptions{
						LabelSelector: labels.Set(build.Runner.Kubernetes.PodLabels).String(),
					},
				)
				require.NoError(t, err)
				require.NotEmpty(t, pods.Items)
				pod := pods.Items[0]

				assert.Equal(t, pod.GetNamespace(), expectedNamespace)

				err = client.
					CoreV1().
					Pods(expectedNamespace).
					Delete(ctx, pod.Name, metav1.DeleteOptions{
						PropagationPolicy: &kubernetes.PropagationPolicy,
					})
				require.NoError(t, err)

				deletedPodNameCh <- pod.Name
			})()

			out, err := buildtest.RunBuildReturningOutput(t, build)

			<-deletedPodNameCh

			if !featureFlagValue {
				assert.Contains(
					t,
					out,
					"ERROR: Job failed (system failure):",
				)
			} else {
				assert.Errorf(t, err, "command terminated with exit code 137")
			}

			validateNamespaceDeleted(t, client, expectedNamespace)

			if tc.finalize != nil {
				tc.finalize(t, client, expectedNamespace)
			}
		})
	}
}

func testKubernetesPublicInternalVariables(t *testing.T, featureFlagName string, featureFlagValue bool) {
	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	ctxTimeout := time.Minute
	client := getTestKubeClusterClient(t)

	containsVerifyFn := func(t *testing.T, v common.JobVariable, envNames []string, envValues []string) {
		assert.Contains(t, envNames, v.Key)
		assert.Contains(t, envValues, v.Value)
	}

	tests := map[string]struct {
		variable common.JobVariable
		verifyFn func(*testing.T, common.JobVariable, []string, []string)
	}{
		"internal variable": {
			variable: common.JobVariable{
				Key:      "my_internal_variable",
				Value:    "my internal variable",
				Internal: true,
			},
			verifyFn: containsVerifyFn,
		},
		"public variable": {
			variable: common.JobVariable{
				Key:    "my_public_variable",
				Value:  "my public variable",
				Public: true,
			},
			verifyFn: containsVerifyFn,
		},
		"regular variable": {
			variable: common.JobVariable{
				Key:   "my_regular_variable",
				Value: "my regular variable",
			},
			verifyFn: func(t *testing.T, v common.JobVariable, envNames []string, envValues []string) {
				assert.NotContains(t, envNames, v.Key)
				assert.NotContains(t, envValues, v.Value)
			},
		},
	}

	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			t.Parallel()

			build := getTestBuild(t, func() (common.JobResponse, error) {
				jobResponse, err := common.GetRemoteBuildResponse(
					"sleep 5000",
				)
				require.NoError(t, err)

				jobResponse.Credentials = []common.Credentials{
					{
						Type:     "registry",
						URL:      "http://example.com",
						Username: "user",
						Password: "password",
					},
				}

				jobResponse.Variables = []common.JobVariable{
					tc.variable,
				}

				return jobResponse, nil
			})
			buildtest.SetBuildFeatureFlag(build, featureFlagName, featureFlagValue)

			deletedPodNameCh := make(chan string)
			defer buildtest.OnUserStage(build, func() {
				ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
				defer cancel()
				pods, err := client.CoreV1().Pods(ciNamespace).List(
					ctx,
					metav1.ListOptions{
						LabelSelector: labels.Set(build.Runner.Kubernetes.PodLabels).String(),
					},
				)
				require.NoError(t, err)
				require.NotEmpty(t, pods.Items)
				pod := pods.Items[0]

				var c *v1.Container
				for _, container := range pod.Spec.Containers {
					if container.Name == "build" {
						c = &container
						break
					}
				}
				require.NotNil(t, c)

				envNames := make([]string, 0)
				envValues := make([]string, 0)
				for _, env := range c.Env {
					envNames = append(envNames, env.Name)
					envValues = append(envValues, env.Value)
				}
				require.NotEmpty(t, envNames)
				require.NotEmpty(t, envValues)

				tc.verifyFn(t, tc.variable, envNames, envValues)

				err = client.
					CoreV1().
					Pods(ciNamespace).
					Delete(ctx, pod.Name, metav1.DeleteOptions{
						PropagationPolicy: &kubernetes.PropagationPolicy,
					})
				require.NoError(t, err)

				deletedPodNameCh <- pod.Name
			})()

			_, err := buildtest.RunBuildReturningOutput(t, build)

			<-deletedPodNameCh

			assert.Errorf(t, err, "command terminated with exit code 137")
		})
	}
}

func testKubernetesWaitResources(t *testing.T, featureFlagName string, featureFlagValue bool) {
	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	secretName := func() string {
		return fmt.Sprintf("my-secret-1-%d", rand.Uint64())
	}
	saName := func() string {
		return fmt.Sprintf("my-serviceaccount-%d", rand.Uint64())
	}
	client := getTestKubeClusterClient(t)

	tests := map[string]struct {
		init             func(t *testing.T, secretName, saName string, build *common.Build, client *k8s.Clientset)
		finalize         func(t *testing.T, secretName, saName string, client *k8s.Clientset)
		checkMaxAttempts int
		imagePullSecret  []string
		serviceAccount   string
		expectedErr      bool
	}{
		"no resources available": {
			checkMaxAttempts: 1,
			imagePullSecret:  []string{secretName()},
			serviceAccount:   saName(),
			expectedErr:      true,
		},
		"only serviceaccount set": {
			serviceAccount: kubernetes.DefaultResourceIdentifier,
		},
		"secret not set but serviceaccount available": {
			checkMaxAttempts: 1,
			imagePullSecret:  []string{secretName()},
			serviceAccount:   kubernetes.DefaultResourceIdentifier,
			expectedErr:      true,
		},
		"secret made available while waiting for resources": {
			checkMaxAttempts: 10,
			imagePullSecret:  []string{secretName()},
			init: func(t *testing.T, secretName, saName string, build *common.Build, client *k8s.Clientset) {
				time.Sleep(time.Second * 3)
				ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
				defer cancel()

				s := &v1.Secret{
					TypeMeta: metav1.TypeMeta{
						Kind: "Secret",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name: secretName,
					},
					Data: map[string][]byte{},
				}

				_, err := kubernetes.CreateTestKubernetesResource(ctx, client, ciNamespace, s)
				require.NoError(t, err)
			},
			finalize: func(t *testing.T, secretName, saName string, client *k8s.Clientset) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
				defer cancel()

				err := client.
					CoreV1().
					Secrets(ciNamespace).
					Delete(ctx, secretName, metav1.DeleteOptions{})
				require.NoError(t, err)
			},
		},
		"serviceaccount made available while waiting for resources": {
			init: func(t *testing.T, secretName, saName string, build *common.Build, client *k8s.Clientset) {
				time.Sleep(time.Second * 3)
				ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
				defer cancel()

				sa := &v1.ServiceAccount{
					TypeMeta: metav1.TypeMeta{
						Kind: "ServiceAccount",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name: saName,
					},
				}

				_, err := kubernetes.CreateTestKubernetesResource(ctx, client, ciNamespace, sa)
				require.NoError(t, err)
			},
			checkMaxAttempts: 10,
			serviceAccount:   saName(),
			finalize: func(t *testing.T, secretName, saName string, client *k8s.Clientset) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
				defer cancel()

				err := client.
					CoreV1().
					ServiceAccounts(ciNamespace).
					Delete(ctx, saName, metav1.DeleteOptions{})
				require.NoError(t, err)
			},
		},
	}

	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			t.Parallel()

			build := getTestBuild(t, func() (common.JobResponse, error) {
				jobResponse, err := common.GetRemoteBuildResponse(
					"echo Hello World",
				)
				require.NoError(t, err)

				jobResponse.Credentials = []common.Credentials{
					{
						Type:     "registry",
						URL:      "http://example.com",
						Username: "user",
						Password: "password",
					},
				}

				return jobResponse, nil
			})
			build.Runner.Kubernetes.ResourceAvailabilityCheckMaxAttempts = tc.checkMaxAttempts
			build.Runner.Kubernetes.ImagePullSecrets = tc.imagePullSecret
			build.Runner.Kubernetes.ServiceAccount = tc.serviceAccount
			buildtest.SetBuildFeatureFlag(build, featureFlagName, featureFlagValue)

			var secretName string
			if len(tc.imagePullSecret) > 0 {
				secretName = tc.imagePullSecret[0]
			}

			saName := tc.serviceAccount

			if tc.init != nil {
				go tc.init(t, secretName, saName, build, client)
			}

			out, err := buildtest.RunBuildReturningOutput(t, build)

			if tc.finalize != nil {
				tc.finalize(t, secretName, saName, client)
			}

			if tc.expectedErr {
				assert.Error(t, err, "checking ImagePullSecret: couldn't find ImagePullSecret or ServiceAccount")
				return
			}

			assert.NoError(t, err)
			assert.Contains(t, out, "Hello World")
		})
	}
}

func testKubernetesClusterWarningEvent(t *testing.T, featureFlagName string, featureFlagValue bool) {
	t.Skip("TODO: Fix test doesn't seem to wait for all events sometimes, https://gitlab.com/gitlab-org/gitlab-runner/-/jobs/8532408889")
	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	tests := map[string]struct {
		image           string
		retrieveWarning bool
		verifyFn        func(*testing.T, string, error)
	}{
		"invalid image": {
			image:           "alpine:invalid-tag",
			retrieveWarning: true,
			verifyFn: func(t *testing.T, out string, err error) {
				assert.Error(t, err)
				assert.Contains(
					t,
					out,
					"WARNING: Event retrieved from the cluster: Failed to pull image \"alpine:invalid-tag\"",
				)
				assert.Contains(t, out, "WARNING: Event retrieved from the cluster: Error: ErrImagePull")
				assert.Contains(t, out, "WARNING: Event retrieved from the cluster: Error: ImagePullBackOff")
			},
		},
		"invalid image with feature flag disabled": {
			image: "alpine:invalid-tag",
			verifyFn: func(t *testing.T, out string, err error) {
				assert.Error(t, err)
				assert.Contains(
					t,
					out,
					"WARNING: Event retrieved from the cluster: Failed to pull image \"alpine:invalid-tag\"",
				)
				assert.Contains(t, out, "WARNING: Event retrieved from the cluster: Error: ErrImagePull")
				assert.Contains(t, out, "WARNING: Event retrieved from the cluster: Error: ImagePullBackOff")
			},
		},
	}

	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			t.Parallel()

			build := getTestBuild(t, func() (common.JobResponse, error) {
				jobResponse, err := common.GetRemoteBuildResponse(
					"echo Hello World",
				)
				require.NoError(t, err)

				return jobResponse, nil
			})
			build.Runner.Kubernetes.Image = tc.image
			buildtest.SetBuildFeatureFlag(build, featureFlagName, featureFlagValue)
			buildtest.SetBuildFeatureFlag(build, "FF_RETRIEVE_POD_WARNING_EVENTS", tc.retrieveWarning)
			build.Runner.Kubernetes.HelperImage = "registry.gitlab.com/gitlab-org/gitlab-runner/gitlab-runner-helper:x86_64-latest"

			out, err := buildtest.RunBuildReturningOutput(t, build)
			tc.verifyFn(t, out, err)
		})
	}
}

// TestLogDeletionAttach tests the outcome when the log files are all deleted
func TestLogDeletionAttach(t *testing.T) {
	t.Parallel()

	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	t.Skip("Log deletion test temporary skipped: issue https://gitlab.com/gitlab-org/gitlab-runner/-/issues/27755")

	tests := []struct {
		stage            string
		outputAssertions func(t *testing.T, out string, pod string)
	}{
		{
			stage: "step_", // Any script the user defined
			outputAssertions: func(t *testing.T, out string, pod string) {
				assert.Contains(
					t,
					out,
					"ERROR: Job failed: command terminated with exit code 100",
				)
			},
		},
		{
			stage: string(common.BuildStagePrepare),
			outputAssertions: func(t *testing.T, out string, pod string) {
				assert.Contains(
					t,
					out,
					"ERROR: Job failed: command terminated with exit code 100",
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.stage, func(t *testing.T) {
			t.Parallel()

			build := getTestBuild(t, func() (common.JobResponse, error) {
				return common.GetRemoteBuildResponse(
					"sleep 5000",
				)
			})
			buildtest.SetBuildFeatureFlag(build, featureflags.UseLegacyKubernetesExecutionStrategy, false)

			deletedPodNameCh := make(chan string)
			defer buildtest.OnUserStage(build, func() {
				client := getTestKubeClusterClient(t)
				pods, err := client.
					CoreV1().
					Pods(ciNamespace).
					List(context.Background(), metav1.ListOptions{
						LabelSelector: labels.Set(build.Runner.Kubernetes.PodLabels).String(),
					})
				require.NoError(t, err)
				require.NotEmpty(t, pods.Items)
				pod := pods.Items[0]
				config, err := kubernetes.GetKubeClientConfig(new(common.KubernetesConfig))
				require.NoError(t, err)
				logsPath := fmt.Sprintf("/logs-%d-%d", build.JobInfo.ProjectID, build.JobResponse.ID)
				opts := kubernetes.ExecOptions{
					Namespace:  pod.Namespace,
					PodName:    pod.Name,
					KubeClient: client,
					Stdin:      true,
					In:         strings.NewReader(fmt.Sprintf("rm -rf %s/*", logsPath)),
					Out:        io.Discard,
					Command:    []string{"/bin/sh"},
					Config:     config,
					Executor:   &kubernetes.DefaultRemoteExecutor{},
				}
				err = opts.Run()
				require.NoError(t, err)

				deletedPodNameCh <- pod.Name
			})()

			out, err := buildtest.RunBuildReturningOutput(t, build)
			require.Error(t, err)
			assert.True(t, err != nil, "No error returned")

			tt.outputAssertions(t, out, <-deletedPodNameCh)
		})
	}
}

// This test reproduces the bug reported in https://gitlab.com/gitlab-org/gitlab-runner/issues/2583
func TestPrepareIssue2583(t *testing.T) {
	t.Parallel()

	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	namespace := "my_namespace"
	serviceAccount := "my_account"

	runnerConfig := &common.RunnerConfig{
		RunnerSettings: common.RunnerSettings{
			Executor: common.ExecutorKubernetes,
			Kubernetes: &common.KubernetesConfig{
				Image:                          "an/image:latest",
				Namespace:                      namespace,
				NamespaceOverwriteAllowed:      ".*",
				ServiceAccount:                 serviceAccount,
				ServiceAccountOverwriteAllowed: ".*",
			},
		},
	}

	build := getTestBuild(t, func() (common.JobResponse, error) {
		return common.JobResponse{
			Variables: []common.JobVariable{
				{Key: kubernetes.NamespaceOverwriteVariableName, Value: "namespace"},
				{Key: kubernetes.ServiceAccountOverwriteVariableName, Value: "sa"},
			},
		}, nil
	})

	e := kubernetes.NewDefaultExecutorForTest()

	// TODO: handle the context properly with https://gitlab.com/gitlab-org/gitlab-runner/-/issues/27932
	prepareOptions := common.ExecutorPrepareOptions{
		Config:  runnerConfig,
		Build:   build,
		Context: context.TODO(),
	}

	err := e.Prepare(prepareOptions)
	assert.NoError(t, err)
	assert.Equal(t, namespace, runnerConfig.Kubernetes.Namespace)
	assert.Equal(t, serviceAccount, runnerConfig.Kubernetes.ServiceAccount)
}

func TestDeletedPodSystemFailureDuringExecution(t *testing.T) {
	t.Skip("Test is flaky, error is sometimes ... status is Failed, https://gitlab.com/gitlab-org/gitlab-runner/-/jobs/8521362284, second test seems entirely wrong")
	t.Parallel()

	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	tests := []struct {
		stage            string
		outputAssertions func(t *testing.T, out string, pod string)
	}{
		{
			stage: "step_", // Any script the user defined
			outputAssertions: func(t *testing.T, out string, pod string) {
				assert.Contains(
					t,
					out,
					"ERROR: Job failed (system failure):",
				)

				assert.Contains(
					t,
					out,
					fmt.Sprintf("pods %q not found", pod),
				)
			},
		},
		{
			stage: string(common.BuildStagePrepare),
			outputAssertions: func(t *testing.T, out string, pod string) {
				assert.Contains(
					t,
					out,
					"ERROR: Job failed (system failure):",
				)

				assert.Contains(
					t,
					out,
					fmt.Sprintf("pods %q not found", pod),
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.stage, func(t *testing.T) {
			t.Parallel()

			build := getTestBuild(t, common.GetRemoteLongRunningBuild)
			// It's not possible to get this kind of information on the legacy execution path.
			buildtest.SetBuildFeatureFlag(build, featureflags.UseLegacyKubernetesExecutionStrategy, false)

			deletedPodNameCh := make(chan string)
			defer buildtest.OnStage(build, tt.stage, func() {
				client := getTestKubeClusterClient(t)
				pods, err := client.CoreV1().Pods(ciNamespace).List(
					context.Background(),
					metav1.ListOptions{
						LabelSelector: labels.Set(build.Runner.Kubernetes.PodLabels).String(),
					},
				)
				require.NoError(t, err)
				require.NotEmpty(t, pods.Items)
				pod := pods.Items[0]
				err = client.
					CoreV1().
					Pods(ciNamespace).
					Delete(context.Background(), pod.Name, metav1.DeleteOptions{})
				require.NoError(t, err)

				deletedPodNameCh <- pod.Name
			})()

			out, err := buildtest.RunBuildReturningOutput(t, build)
			assert.True(t, kubernetes.IsKubernetesPodNotFoundError(err), "expected err NotFound, but got %T", err)

			tt.outputAssertions(t, out, <-deletedPodNameCh)
		})
	}
}

func testKubernetesWithNonRootSecurityContext(t *testing.T, featureFlagName string, featureFlagValue bool) {
	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	build := getTestBuild(t, func() (common.JobResponse, error) {
		return common.GetRemoteBuildResponse("id")
	})
	build.Image.Name = common.TestAlpineNoRootImage

	runAsNonRoot := true
	runAsUser := int64(1895034)
	build.Runner.Kubernetes.PodSecurityContext = common.KubernetesPodSecurityContext{
		RunAsNonRoot: &runAsNonRoot,
		RunAsUser:    &runAsUser,
	}

	// We override the home directory, else we get this error from the git run:
	// 	```
	// 	Fetching changes...
	// 	error: could not lock config file //.gitconfig: Permission denied
	// 	ERROR: Job failed: command terminated with exit code 1
	// 	```
	build.Variables = append(build.Variables, common.JobVariable{
		Key:   "HOME",
		Value: "/dev/shm",
	})

	buildtest.SetBuildFeatureFlag(build, featureFlagName, featureFlagValue)

	out, err := buildtest.RunBuildReturningOutput(t, build)
	assert.NoError(t, err)
	assert.Contains(t, out, fmt.Sprintf("uid=%d gid=0(root)", runAsUser))
}

func testKubernetesBashFeatureFlag(t *testing.T, featureFlagName string, featureFlagValue bool) {
	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	tests := []struct {
		script              string
		expectedContent     string
		expectedErrExitCode int
	}{
		{
			script:          "export hello=world; echo \"hello $hello\"",
			expectedContent: "hello world",
		},
		{
			script:              "return 129",
			expectedErrExitCode: 129,
		},
		{
			script:              "exit 128",
			expectedErrExitCode: 128,
		},
		{
			script:              "eco 'function error'",
			expectedErrExitCode: 127,
		},
		{
			script:              "{{}",
			expectedErrExitCode: 127,
		},
	}

	for _, tc := range tests {
		t.Run("", func(t *testing.T) {
			t.Parallel()

			build := getTestBuild(t, common.GetRemoteSuccessfulBuild)
			buildtest.SetBuildFeatureFlag(build, featureFlagName, featureFlagValue)

			build.Image.Name = common.TestAlpineImage
			build.Runner.Shell = "bash"
			build.JobResponse.Steps = common.Steps{
				common.Step{
					Name:   common.StepNameScript,
					Script: []string{tc.script},
				},
			}

			out, err := buildtest.RunBuildReturningOutput(t, build)
			assert.Contains(t, out, tc.expectedContent)

			if tc.expectedErrExitCode != 0 {
				var buildError *common.BuildError
				if assert.ErrorAs(t, err, &buildError) {
					assert.Equal(t, tc.expectedErrExitCode, buildError.ExitCode)
				}
			}
		})
	}
}

func testKubernetesContainerHookFeatureFlag(t *testing.T, featureFlagName string, featureFlagValue bool) {
	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	tests := map[string]struct {
		image           string
		shell           string
		lifecycleCfg    common.KubernetesContainerLifecyle
		steps           common.Steps
		validateOutputs func(t *testing.T, out string, err error)
	}{
		"invalid hook configuration: more than one handler type": {
			lifecycleCfg: common.KubernetesContainerLifecyle{
				PreStop: &common.KubernetesLifecycleHandler{
					Exec: &common.KubernetesLifecycleExecAction{
						Command: []string{"touch", "/builds/postStart.txt"},
					},
					HTTPGet: &common.KubernetesLifecycleHTTPGet{
						Host: "localhost",
						Port: 8080,
					},
				},
			},
			validateOutputs: func(t *testing.T, out string, err error) {
				require.Error(t, err)
				assert.Contains(t, out, "ERROR: Job failed (system failure):")
				assert.Contains(t, out, "may not specify more than 1 handler type")
			},
		},
		"postStart exec hook bash": {
			steps: common.Steps{
				common.Step{
					Name: common.StepNameScript,
					Script: []string{
						"ls -l /builds",
					},
				},
			},
			lifecycleCfg: common.KubernetesContainerLifecyle{
				PostStart: &common.KubernetesLifecycleHandler{
					Exec: &common.KubernetesLifecycleExecAction{
						Command: []string{"touch", "/builds/postStart.txt"},
					},
				},
			},
			validateOutputs: func(t *testing.T, out string, err error) {
				require.NoError(t, err)
				assert.Contains(t, out, "Job succeeded")
				assert.Contains(t, out, "postStart.txt")
			},
		},
		"postStart exec hook pwsh": {
			image: common.TestPwshImage,
			shell: shells.SNPwsh,
			steps: common.Steps{
				common.Step{
					Name: common.StepNameScript,
					Script: []string{
						"Get-ChildItem /builds",
					},
				},
			},
			lifecycleCfg: common.KubernetesContainerLifecyle{
				PostStart: &common.KubernetesLifecycleHandler{
					Exec: &common.KubernetesLifecycleExecAction{
						Command: []string{"touch", "/builds/postStart.txt"},
					},
				},
			},
			validateOutputs: func(t *testing.T, out string, err error) {
				require.NoError(t, err)
				assert.Contains(t, out, "Job succeeded")
				assert.Contains(t, out, "postStart.txt")
			},
		},
	}

	for tn, tt := range tests {
		t.Run(tn, func(t *testing.T) {
			t.Parallel()

			build := getTestBuild(t, common.GetRemoteSuccessfulBuild)
			buildtest.SetBuildFeatureFlag(build, featureFlagName, featureFlagValue)
			build.Image.Name = common.TestAlpineImage
			build.Runner.RunnerSettings.Kubernetes.ContainerLifecycle = tt.lifecycleCfg

			if tt.image != "" {
				build.Image.Name = tt.image
			}

			if tt.shell != "" {
				build.Runner.Shell = tt.shell
			}

			if tt.steps != nil {
				build.JobResponse.Steps = tt.steps
			}

			out, err := buildtest.RunBuildReturningOutput(t, build)
			tt.validateOutputs(t, out, err)
		})
	}
}

func getTestBuildWithImage(t *testing.T, image string, getJobResponse func() (common.JobResponse, error)) *common.Build {
	jobResponse, err := getJobResponse()
	assert.NoError(t, err)

	podUUID, err := helpers.GenerateRandomUUID(8)
	require.NoError(t, err)

	tt := strings.Split(t.Name(), "/")
	slices.Reverse(tt)

	nodeSelector := map[string]string{}
	if os.Getenv("GITLAB_CI") == "true" {
		nodeSelector["app"] = "gitlab-runner-job"
	}

	return &common.Build{
		JobResponse: jobResponse,
		Runner: &common.RunnerConfig{
			RunnerSettings: common.RunnerSettings{
				Executor: common.ExecutorKubernetes,
				Kubernetes: &common.KubernetesConfig{
					Image:      image,
					PullPolicy: common.StringOrArray{common.PullPolicyIfNotPresent},
					PodLabels: map[string]string{
						"test.k8s.gitlab.com/name":      podUUID,
						"test.k8s.gitlab.com/test-name": dns.MakeRFC1123Compatible(strings.Join(tt, ".")),
					},
					Namespace:                        ciNamespace,
					CleanupGracePeriodSeconds:        common.Int64Ptr(5),
					PodTerminationGracePeriodSeconds: common.Int64Ptr(5),
					PollTimeout:                      int((time.Minute * 10).Seconds()),
					NodeSelector:                     nodeSelector,

					CPULimit:            "0.3",
					MemoryRequest:       "150Mi",
					HelperCPULimit:      "0.2",
					HelperMemoryRequest: "150Mi",
				},
			},
		},
	}
}

func getTestBuild(t *testing.T, getJobResponse func() (common.JobResponse, error)) *common.Build {
	return getTestBuildWithImage(t, common.TestAlpineImage, getJobResponse)
}

func getTestBuildWithServices(
	t *testing.T,
	getJobResponse func() (common.JobResponse, error),
	services ...string,
) *common.Build {
	build := getTestBuild(t, getJobResponse)

	for _, service := range services {
		build.Services = append(build.Services, common.Image{
			Name: service,
		})
	}

	return build
}

func getTestKubeClusterClient(t *testing.T) *k8s.Clientset {
	// Taken from the k8s client to create a config with a custom token
	// this token is linked to a separate service account that is not the one set as the
	// service account of the pod running the integration tests.
	// The service account set on the pod has all the permissions GitLab Runner needs to execute
	// builds in Kubernetes, but it doesn't have permissions needed for the integration tests to run,
	// such as listing pods or creating secrets. The admin service account is used specifically for that.
	const (
		tokenPath   = "/var/run/secrets/kubernetes.io/serviceaccount_admin/token"
		tokenCAPath = "/var/run/secrets/kubernetes.io/serviceaccount_admin/ca.crt"
	)

	var config *rest.Config
	if _, err := os.Stat(tokenPath); err != nil {
		config, err = kubernetes.GetKubeClientConfig(new(common.KubernetesConfig))
		require.NoError(t, err)
	} else {
		host, port := os.Getenv("KUBERNETES_SERVICE_HOST"), os.Getenv("KUBERNETES_SERVICE_PORT")
		if len(host) == 0 || len(port) == 0 {
			t.Fatal(rest.ErrNotInCluster)
		}

		token, err := os.ReadFile(tokenPath)
		require.NoError(t, err)

		tlsClientConfig := rest.TLSClientConfig{CAFile: tokenCAPath}

		config = &rest.Config{
			Host:            "https://" + net.JoinHostPort(host, port),
			TLSClientConfig: tlsClientConfig,
			BearerToken:     string(token),
			BearerTokenFile: tokenPath,
		}
	}

	client, err := k8s.NewForConfig(config)
	require.NoError(t, err)

	return client
}

// getSecrets retrieves all the secrets found in the given namespace
// with at least one ownerReference name matching the name given
// If ownerName is an empty string, all the secrets resources found are returned
func getSecrets(client *k8s.Clientset, namespace, ownerName string) ([]v1.Secret, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	credList, err := client.CoreV1().Secrets(namespace).List(
		ctx,
		metav1.ListOptions{},
	)
	if err != nil {
		return nil, err
	}

	credentials := []v1.Secret{}
	for _, cred := range credList.Items {
		if len(cred.OwnerReferences) == 1 && cred.OwnerReferences[0].Name == ownerName {
			credentials = append(credentials, cred)
		}
	}

	return credentials, nil
}

// getConfigMaps retrieves all the configMaps found in the given namespace
// with at least one ownerReference name matching the name given
// If ownerName is an empty string, all the configMaps resources found are returned
func getConfigMaps(client *k8s.Clientset, namespace, ownerName string) ([]v1.ConfigMap, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	configMapList, err := client.CoreV1().ConfigMaps(namespace).List(
		ctx,
		metav1.ListOptions{},
	)
	if err != nil {
		return nil, err
	}

	configMaps := []v1.ConfigMap{}
	for _, cfg := range configMapList.Items {
		if len(cfg.OwnerReferences) == 1 && cfg.OwnerReferences[0].Name == ownerName {
			configMaps = append(configMaps, cfg)
		}
	}

	return configMaps, nil
}

func runMultiPullPolicyBuild(t *testing.T, build *common.Build) error {
	build.Runner.Kubernetes.PullPolicy = common.StringOrArray{
		common.PullPolicyAlways,
		common.PullPolicyIfNotPresent,
	}

	outBuffer := bytes.NewBuffer(nil)

	err := build.Run(&common.Config{}, &common.Trace{Writer: outBuffer})
	require.Error(t, err)
	assert.ErrorIs(t, err, &common.BuildError{FailureReason: common.ImagePullFailure})

	quotedImage := regexp.QuoteMeta("some/non-existing/image")
	warningFmt := `WARNING: Failed to pull image "%s" for container "[^"]+" with policy "%s": image pull failed:`
	attemptFmt := `Attempt #%d: Trying "%s" pull policy for "%s" image for container "[^"]+"`

	// We expect
	//  - the warning for the 1st attempt with "Always", telling us about the pull issue
	//  - the log of the 2nd attempt with "IfNotPresent"
	//  - the warning for the 2. attempt with "IfNotPresent", telling us about the 2nd pull issue
	expectedLogRE := fmt.Sprintf(
		`(?s)%s.*%s.*%s`,
		fmt.Sprintf(warningFmt, quotedImage, "Always"),
		fmt.Sprintf(attemptFmt, 2, "IfNotPresent", quotedImage),
		fmt.Sprintf(warningFmt, quotedImage, "IfNotPresent"),
	)

	assert.Regexp(t, expectedLogRE, outBuffer.String())

	return err
}

func TestKubernetesAllowedImages(t *testing.T) {
	t.Parallel()

	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	type testDef struct {
		AllowedImages []string
		Image         string
		VerifyFn      func(*testing.T, error)
	}
	tests := map[string]testDef{
		// allowed image case
		"allowed image case": {
			AllowedImages: []string{"alpine"},
			Image:         "alpine",
			VerifyFn: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		// disallowed image case
		"disallowed image case": {
			AllowedImages: []string{"alpine"},
			Image:         "ubuntu",
			VerifyFn: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.ErrorIs(t, err, common.ErrDisallowedImage)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			build := getTestBuild(t, common.GetRemoteSuccessfulBuild)
			build.Runner.Kubernetes.AllowedImages = test.AllowedImages
			build.Image.Name = test.Image

			err := build.Run(&common.Config{}, &common.Trace{Writer: os.Stdout})
			test.VerifyFn(t, err)
		})
	}
}

func TestKubernetesAllowedServices(t *testing.T) {
	t.Parallel()

	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	type testDef struct {
		AllowedServices []string
		Services        common.Services
		VerifyFn        func(*testing.T, error)
	}
	tests := map[string]testDef{
		"allowed service case": {
			AllowedServices: []string{"alpine", "debian"},
			Services: common.Services{
				common.Image{Name: "alpine"},
			},
			VerifyFn: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		"disallowed service case": {
			AllowedServices: []string{"alpine", "debian"},
			Services: common.Services{
				common.Image{Name: "alpine"},
				common.Image{Name: "ubuntu"},
			},
			VerifyFn: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.ErrorIs(t, err, common.ErrDisallowedImage)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			build := getTestBuild(t, common.GetRemoteSuccessfulBuild)
			build.Runner.Kubernetes.AllowedServices = test.AllowedServices
			build.Services = test.Services

			err := build.Run(&common.Config{}, &common.Trace{Writer: os.Stdout})
			test.VerifyFn(t, err)
		})
	}
}

func TestCleanupProjectGitClone(t *testing.T) {
	t.Parallel()

	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	buildtest.RunBuildWithCleanupGitClone(t, getTestBuild(t, common.GetRemoteSuccessfulBuild))
}

func TestCleanupProjectGitFetch(t *testing.T) {
	t.Parallel()

	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	untrackedFilename := "untracked"

	build := getTestBuild(t, func() (common.JobResponse, error) {
		return common.GetRemoteBuildResponse(
			buildtest.GetNewUntrackedFileIntoSubmodulesCommands(untrackedFilename, "", "")...,
		)
	})

	buildtest.RunBuildWithCleanupGitFetch(t, build, untrackedFilename)
}

func TestCleanupProjectGitSubmoduleNormal(t *testing.T) {
	t.Parallel()

	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	untrackedFile := "untracked"
	untrackedSubmoduleFile := "untracked_submodule"

	build := getTestBuild(t, func() (common.JobResponse, error) {
		return common.GetRemoteBuildResponse(
			buildtest.GetNewUntrackedFileIntoSubmodulesCommands(untrackedFile, untrackedSubmoduleFile, "")...,
		)
	})

	buildtest.RunBuildWithCleanupNormalSubmoduleStrategy(t, build, untrackedFile, untrackedSubmoduleFile)
}

func TestCleanupProjectGitSubmoduleRecursive(t *testing.T) {
	t.Parallel()

	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	untrackedFile := "untracked"
	untrackedSubmoduleFile := "untracked_submodule"
	untrackedSubSubmoduleFile := "untracked_submodule_submodule"

	build := getTestBuild(t, func() (common.JobResponse, error) {
		return common.GetRemoteBuildResponse(
			buildtest.GetNewUntrackedFileIntoSubmodulesCommands(
				untrackedFile,
				untrackedSubmoduleFile,
				untrackedSubSubmoduleFile)...,
		)
	})

	buildtest.RunBuildWithCleanupRecursiveSubmoduleStrategy(
		t,
		build,
		untrackedFile,
		untrackedSubmoduleFile,
		untrackedSubSubmoduleFile,
	)
}

func TestKubernetesPwshFeatureFlag(t *testing.T) {
	t.Parallel()

	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	tests := []struct {
		script              string
		expectedRegexp      string
		expectedErrExitCode int
	}{
		{
			script:         "Write-Output $PSVersionTable",
			expectedRegexp: "PSEdition +Core",
		},
		{
			script:         "return 129",
			expectedRegexp: "Job succeeded",
		},
		{
			script:              "Write-Error 'should fail'",
			expectedErrExitCode: 1,
		},
		{
			script:              "Exit 128",
			expectedErrExitCode: 128,
		},
		{
			script:              "$host.SetShouldExit(130)",
			expectedErrExitCode: 130,
		},
		{
			script:              "eco 'function error'",
			expectedErrExitCode: 1,
		},
		{
			script:              "syntax error {{}",
			expectedErrExitCode: 1,
		},
	}

	for _, tc := range tests {
		t.Run("", func(t *testing.T) {
			t.Parallel()

			build := getTestBuild(t, common.GetRemoteSuccessfulBuild)

			build.Image.Name = common.TestPwshImage
			build.Runner.Shell = shells.SNPwsh
			build.JobResponse.Steps = common.Steps{
				common.Step{
					Name:   common.StepNameScript,
					Script: []string{tc.script},
				},
			}

			out, err := buildtest.RunBuildReturningOutput(t, build)
			assert.Regexp(t, regexp.MustCompile(tc.expectedRegexp), out)

			if tc.expectedErrExitCode != 0 {
				var buildError *common.BuildError
				if assert.ErrorAs(t, err, &buildError) {
					assert.Equal(t, tc.expectedErrExitCode, buildError.ExitCode)
				}
			}
		})
	}
}

func TestKubernetesProcessesInBackground(t *testing.T) {
	t.Parallel()

	// Check fix for https://gitlab.com/gitlab-org/gitlab-runner/-/issues/2880

	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	tests := map[string]struct {
		shell  string
		image  string
		script []string
	}{
		"bash shell": {
			shell: "bash",
			image: common.TestAlpineImage,
		},
		"pwsh shell": {
			shell: shells.SNPwsh,
			image: common.TestPwshImage,
		},
	}

	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			t.Parallel()

			build := getTestBuild(t, common.GetRemoteSuccessfulBuild)

			build.Image.Name = tc.image
			build.Runner.Shell = tc.shell
			build.JobResponse.Steps = common.Steps{
				common.Step{
					Name: common.StepNameScript,
					Script: []string{
						`sleep infinity &`,
						`mkdir out && echo "Hello, world" > out/greeting`,
					},
				},
				common.Step{
					Name: common.StepNameAfterScript,
					Script: []string{
						`echo I should be running`,
					},
				},
			}

			out, err := buildtest.RunBuildReturningOutput(t, build)
			assert.Contains(t, out, "I should be running")
			assert.NoError(t, err)
		})
	}
}

func TestBuildExpandedFileVariable(t *testing.T) {
	t.Parallel()

	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	shellstest.OnEachShell(t, func(t *testing.T, shell string) {
		build := getTestBuild(t, common.GetRemoteSuccessfulBuild)
		buildtest.RunBuildWithExpandedFileVariable(t, build.Runner, nil)
	})
}

func TestConflictingPullPolicies(t *testing.T) {
	t.Parallel()

	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	tests := map[string]struct {
		imagePullPolicies   []common.DockerPullPolicy
		pullPolicy          common.StringOrArray
		allowedPullPolicies []common.DockerPullPolicy
		wantErrMsg          string
	}{
		"allowed_pull_policies configured, default pull_policy": {
			imagePullPolicies:   nil,
			pullPolicy:          nil,
			allowedPullPolicies: []common.DockerPullPolicy{common.PullPolicyIfNotPresent},
			wantErrMsg:          fmt.Sprintf(common.IncompatiblePullPolicy, "[]", "Runner config (default)", "[IfNotPresent]"),
		},
		"allowed_pull_policies and pull_policy configured": {
			imagePullPolicies:   nil,
			pullPolicy:          common.StringOrArray{common.PullPolicyNever},
			allowedPullPolicies: []common.DockerPullPolicy{common.PullPolicyIfNotPresent},
			wantErrMsg:          fmt.Sprintf(common.IncompatiblePullPolicy, "[Never]", "Runner config", "[IfNotPresent]"),
		},
		"allowed_pull_policies and image pull_policy configured": {
			imagePullPolicies:   []common.DockerPullPolicy{common.PullPolicyAlways},
			pullPolicy:          nil,
			allowedPullPolicies: []common.DockerPullPolicy{common.PullPolicyIfNotPresent},
			wantErrMsg:          fmt.Sprintf(common.IncompatiblePullPolicy, "[Always]", "GitLab pipeline config", "[IfNotPresent]"),
		},
		"all configured": {
			imagePullPolicies:   []common.DockerPullPolicy{common.PullPolicyAlways},
			pullPolicy:          common.StringOrArray{common.PullPolicyNever},
			allowedPullPolicies: []common.DockerPullPolicy{common.PullPolicyIfNotPresent},
			wantErrMsg:          fmt.Sprintf(common.IncompatiblePullPolicy, "[Always]", "GitLab pipeline config", "[IfNotPresent]"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			build := getTestBuild(t, common.GetRemoteSuccessfulBuild)
			build.JobResponse.Image.PullPolicies = test.imagePullPolicies
			build.Runner.RunnerSettings.Kubernetes.PullPolicy = test.pullPolicy
			build.Runner.RunnerSettings.Kubernetes.AllowedPullPolicies = test.allowedPullPolicies

			gotErr := build.Run(&common.Config{}, &common.Trace{Writer: os.Stdout})

			require.Error(t, gotErr)
			assert.Contains(t, gotErr.Error(), test.wantErrMsg)
			assert.Contains(t, gotErr.Error(), "invalid pull policy for container")
		})
	}
}

func Test_CaptureServiceLogs(t *testing.T) {
	t.Parallel()

	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	tests := map[string]struct {
		buildVars []common.JobVariable
		assert    func(string, error)
	}{
		"enabled": {
			buildVars: []common.JobVariable{
				{
					Key:    "CI_DEBUG_SERVICES",
					Value:  "1",
					Public: true,
				},
			},
			assert: func(out string, err error) {
				assert.NoError(t, err)
				assert.NotContains(t, out, "WARNING: invalid value '1' for CI_DEBUG_SERVICES variable")
				assert.Regexp(t, `\[service:postgres-db\] .* The files belonging to this database system will be owned by user "postgres"`, out)
				assert.Regexp(t, `\[service:postgres-db\] .* database system is ready to accept connections`, out)
				assert.Regexp(t, `\[service:redis-cache\] .* oO0OoO0OoO0Oo Redis is starting oO0OoO0OoO0O`, out)
				assert.Regexp(t, `\[service:redis-cache\] .* Ready to accept connections`, out)
			},
		},
		"not enabled": {
			assert: func(out string, err error) {
				assert.NoError(t, err)
				assert.NotContains(t, out, "WARNING: invalid value '1' for CI_DEBUG_SERVICES variable")
				assert.NotRegexp(t, `\[service:postgres-db\] .* The files belonging to this database system will be owned by user "postgres"`, out)
				assert.NotRegexp(t, `\[service:postgres-db\] .* database system is ready to accept connections`, out)
				assert.NotRegexp(t, `\[service:redis-cache\] .* oO0OoO0OoO0Oo Redis is starting oO0OoO0OoO0O`, out)
				assert.NotRegexp(t, `\[service:redis-cache\] .* Ready to accept connections`, out)
			},
		},
		"bogus value": {
			buildVars: []common.JobVariable{{
				Key:    "CI_DEBUG_SERVICES",
				Value:  "blammo",
				Public: true,
			}},
			assert: func(out string, err error) {
				assert.NoError(t, err)
				assert.Contains(t, out, `WARNING: CI_DEBUG_SERVICES: expected bool got "blammo", using default value: false`)
				assert.NotRegexp(t, `\[service:postgres-db\] .* The files belonging to this database system will be owned by user "postgres"`, out)
				assert.NotRegexp(t, `\[service:postgres-db\] .* database system is ready to accept connections`, out)
				assert.NotRegexp(t, `\[service:redis-cache\] .* oO0OoO0OoO0Oo Redis is starting oO0OoO0OoO0O`, out)
				assert.NotRegexp(t, `\[service:redis-cache\] .* Ready to accept connections`, out)
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			build := getTestBuildWithServices(t, common.GetRemoteSuccessfulBuild, "postgres:14.4", "redis:7.0")
			build.Services[0].Alias = "db"
			build.Services[1].Alias = "cache"
			build.Variables = tt.buildVars
			build.Variables = append(build.Variables, common.JobVariable{
				Key:    "POSTGRES_PASSWORD",
				Value:  "password",
				Public: true,
			})

			out, err := buildtest.RunBuildReturningOutput(t, build)
			tt.assert(out, err)
		})
	}
}

// When testing with minikube, the following commands may be used to
// properly configure the cluster:
//
// minikube config set container-runtime containerd
// minikube config set feature-gates "ProcMountType=true"
//
// Note that the cluster must be re-initialized after making these changes:
//
// minikube delete
// minikube start
func TestKubernetesProcMount(t *testing.T) {
	t.Parallel()
	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	privileged := false

	// Generate a temporary Pod with procMount set to Unmasked.
	// If the cluster supports the ProcMount feature, then this will be reflected
	// in the PodSpec. If the cluster does not support this feature, the API server
	// will return DefaultProcMount.
	tmpPod := getTestBuild(t, func() (common.JobResponse, error) {
		return common.GetRemoteBuildResponse("cat")
	})

	tmpPod.Runner.RunnerSettings.Kubernetes.BuildContainerSecurityContext = common.KubernetesContainerSecurityContext{
		ProcMount:  v1.UnmaskedProcMount,
		Privileged: &privileged,
	}

	shouldSkipCh := make(chan bool)
	cleanup := buildtest.OnUserStage(tmpPod, func() {
		client := getTestKubeClusterClient(t)

		pods, err := client.
			CoreV1().
			Pods(ciNamespace).
			List(context.Background(), metav1.ListOptions{
				LabelSelector: labels.Set(tmpPod.Runner.Kubernetes.PodLabels).String(),
			})

		require.NoError(t, err)
		require.NotEmpty(t, pods.Items)

		pod := pods.Items[0]

		require.NotEmpty(t, pod.Spec.Containers)

		container := pod.Spec.Containers[0]

		procMount := container.SecurityContext.ProcMount
		shouldSkipCh <- procMount == nil || *procMount != v1.UnmaskedProcMount
	})
	defer cleanup()

	buildtest.RunBuildReturningOutput(t, tmpPod)

	shouldSkip := <-shouldSkipCh
	if shouldSkip {
		t.Skip("ProcMountType feature not supported on cluster -- skipping tests")
		return
	}

	// If we get here, then we have validated that the cluster does indeed support the
	// ProcMount feature, and we can proceed with a more thorough set of tests.

	build := getTestBuild(t, func() (common.JobResponse, error) {
		return common.GetRemoteBuildResponse("unshare --fork -r -p --mount-proc true")
	})

	var buildErr *common.BuildError

	tests := map[string]struct {
		procMount v1.ProcMountType
		validate  func(*testing.T, string, error)
	}{
		"Default": {
			procMount: v1.DefaultProcMount,
			validate: func(t *testing.T, out string, err error) {
				assert.ErrorAs(t, err, &buildErr)
				assert.Contains(t, out, "Job failed")
			},
		},
		"default": {
			procMount: v1.ProcMountType("default"),
			validate: func(t *testing.T, out string, err error) {
				assert.ErrorAs(t, err, &buildErr)
				assert.Contains(t, out, "Job failed")
			},
		},
		"Unmasked": {
			procMount: v1.UnmaskedProcMount,
			validate: func(t *testing.T, out string, err error) {
				require.NoError(t, err)
				assert.Contains(t, out, "Job succeeded")
			},
		},
		"unmasked": {
			procMount: v1.ProcMountType("unmasked"),
			validate: func(t *testing.T, out string, err error) {
				require.NoError(t, err)
				assert.Contains(t, out, "Job succeeded")
			},
		},
		"empty": {
			procMount: v1.ProcMountType("   "),
			validate: func(t *testing.T, out string, err error) {
				assert.ErrorAs(t, err, &buildErr)
				assert.Contains(t, out, "Job failed")
			},
		},
		"invalid": {
			procMount: v1.ProcMountType("invalid"),
			validate: func(t *testing.T, out string, err error) {
				assert.ErrorAs(t, err, &buildErr)
				assert.Contains(t, out, "Job failed")
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			build.Runner.RunnerSettings.Kubernetes.BuildContainerSecurityContext = common.KubernetesContainerSecurityContext{
				ProcMount:  test.procMount,
				Privileged: &privileged,
			}

			out, err := buildtest.RunBuildReturningOutput(t, build)
			test.validate(t, out, err)
		})
	}
}

func Test_ContainerOptionsExpansion(t *testing.T) {
	t.Parallel()

	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	build := getTestBuild(t, common.GetRemoteSuccessfulBuild)

	jobVars := common.JobVariables{
		{Key: "CI_DEBUG_SERVICES", Value: "true", Public: true},
		{Key: "POSTGRES_PASSWORD", Value: "password", Public: true},
		{Key: "JOB_IMAGE", Value: "alpine:latest"},
		{Key: "HELPER_IMAGE", Value: "registry.gitlab.com/gitlab-org/gitlab-runner/gitlab-runner-helper:x86_64-latest"},
		{Key: "HELPER_IMAGE_FLAVOR", Value: "alpine"},
		{Key: "SRVS_IMAGE", Value: "postgres:latest"},
		{Key: "SRVS_IMAGE_ALIAS", Value: "db"},
	}
	build.Variables = append(build.Variables, jobVars...)

	build.Runner.Kubernetes.Image = "$JOB_IMAGE"
	build.Runner.Kubernetes.HelperImage = "$HELPER_IMAGE"
	build.Runner.Kubernetes.HelperImageFlavor = "$HELPER_IMAGE_FLAVOR"
	build.Services = []common.Image{
		{Name: "$SRVS_IMAGE", Alias: "$SRVS_IMAGE_ALIAS"},
	}

	out, err := buildtest.RunBuildReturningOutput(t, build)
	assert.NoError(t, err)
	// the helper image name does not appeart in the logs, but the build will fail if the option was not expanded.
	assert.Contains(t, out, "Using Kubernetes executor with image alpine:latest")
	assert.Regexp(t, `\[service:postgres-db\]`, out)
}

func testJobRunningAndPassingWhenServiceStops(t *testing.T, featureFlagName string, featureFlagValue bool) {
	build := getTestBuild(t, func() (common.JobResponse, error) {
		jobResponse, err := common.GetRemoteBuildResponse("sleep 12")
		if err != nil {
			return common.JobResponse{}, err
		}

		jobResponse.Steps = append(
			jobResponse.Steps,
			common.Step{
				Name:   common.StepNameAfterScript,
				Script: []string{"echo after script"},
			},
		)

		return jobResponse, nil
	})

	build.Runner.Kubernetes.Services = []common.Service{{
		Name: counterServiceImage,
	}}

	buildtest.SetBuildFeatureFlag(build, featureFlagName, featureFlagValue)

	err := buildtest.RunBuild(t, build)
	require.NoError(t, err)
}

func TestEntrypointLogging(t *testing.T) {
	t.Skip("TODO: Flaky, fix with https://gitlab.com/gitlab-org/gitlab-runner/-/merge_requests/5175/diffs?commit_id=d424ae620f90db86bacc3696f3b8727886e1f85b")

	t.Run("succeed", testEntrypointLoggingSuccesses)
	t.Run("fail", testEntrypointLoggingFailures)
}

func testEntrypointLoggingFailures(t *testing.T) {
	t.Parallel()

	// When the pollTimeout is smaller than the time it takes for the entrypoint to start the shell, and thus resolve the
	// startupProbe (roughly 1sec * iterations), then the build should fail but still show _some_ of the entrypoint logs (until
	// the pod gets killed because of the timeout)
	// Note: We only use a startup probe in exec mode
	t.Run("startupProbe does not resolve in time", func(t *testing.T) {
		t.Parallel()

		const pollTimeout = 4
		const iterations = 8

		build := getTestBuild(t, common.GetRemoteSuccessfulBuild)
		build.Runner.Kubernetes.Image = "registry.gitlab.com/gitlab-org/gitlab-runner/alpine-entrypoint-pre-post-trap"
		build.Runner.Kubernetes.PollTimeout = pollTimeout
		build.Runner.FeatureFlags = mapFromKeySlices(true, []string{
			featureflags.KubernetesHonorEntrypoint,
			featureflags.UseLegacyKubernetesExecutionStrategy,
		})
		build.Runner.Environment = []string{
			fmt.Sprintf("LOOP_ITERATIONS=%d", iterations),
		}

		out, err := buildtest.RunBuildReturningOutput(t, build)
		if assert.Error(t, err, "expected build to fail, but did not") {
			assert.Contains(t, err.Error(), "timed out")
		}

		// we see some entrypoint logs
		assert.Contains(t, out, "some pre message")
	})
}

func testEntrypointLoggingSuccesses(t *testing.T) {
	t.Parallel()

	const pollTimeout = 12
	const loopIterations = 5
	defaultFeatureFlags := []string{featureflags.ScriptSections, featureflags.PrintPodEvents}

	expectedLogs := func(phase string, count int) []string {
		expectedLogs := make([]string, count*2)
		for idx := 0; idx < count; idx++ {
			// produces something like: "[entrypoint][post][stdout][5/10] "
			expectedLogs[idx] = fmt.Sprintf("[entrypoint][%s][stdout][%d/%d] ", phase, idx, count)
			expectedLogs[idx+count] = fmt.Sprintf("[entrypoint][%s][stderr][%d/%d] ", phase, idx, count)
		}
		return expectedLogs
	}

	runtimeEnvs := map[string]struct {
		shell string
		image string
	}{
		"bash": {shell: shells.Bash, image: "registry.gitlab.com/gitlab-org/gitlab-runner/alpine-entrypoint-pre-post-trap"},
		"pwsh": {shell: shells.SNPwsh, image: "registry.gitlab.com/gitlab-org/gitlab-runner/powershell-entrypoint-pre-post-trap"},
	}
	modes := map[string][]string{
		"attach mode": { /* attach mode is the default, no additional FFs needed */ },
		"exec mode":   {featureflags.UseLegacyKubernetesExecutionStrategy},
	}
	tests := map[string]struct {
		featureFlags  []string
		expectLogs    []string
		notExpectLogs []string
	}{
		"not honoring entrypoint": {
			notExpectLogs: append(expectedLogs("pre", loopIterations), expectedLogs("post", loopIterations)...),
		},
		"honoring entrypoint": {
			featureFlags: []string{featureflags.KubernetesHonorEntrypoint},
			expectLogs:   expectedLogs("pre", loopIterations),
		},
	}

	for runtimeName, runtimeEnv := range runtimeEnvs {
		t.Run(runtimeName, func(t *testing.T) {
			t.Parallel()
			if runtimeName == "pwsh" {
				t.Skip("TODO: pwsh doesn't work")
			}

			for mode, modeFeatureFlags := range modes {
				t.Run(mode, func(t *testing.T) {
					t.Parallel()

					for testName, testConfig := range tests {
						t.Run(testName, func(t *testing.T) {
							t.Parallel()

							build := getTestBuild(t, common.GetRemoteSuccessfulBuild)
							build.Runner.Kubernetes.Image = runtimeEnv.image
							build.Runner.Kubernetes.PollTimeout = pollTimeout
							build.Runner.FeatureFlags = mapFromKeySlices(true, defaultFeatureFlags, modeFeatureFlags, testConfig.featureFlags)
							build.Runner.Environment = []string{
								fmt.Sprintf("LOOP_ITERATIONS=%d", loopIterations),
							}
							build.Runner.Shell = runtimeEnv.shell

							out, err := buildtest.RunBuildReturningOutput(t, build)
							assert.NoError(t, err)

							for _, s := range testConfig.expectLogs {
								occurrences := strings.Count(out, s)
								assert.Equal(t, 1, occurrences, "expected to find %q exactly once, found it %d times", s, occurrences)
							}
							for _, s := range testConfig.notExpectLogs {
								assert.NotContains(t, out, s, "expected output not to contain %q, but does", s)
							}
						}) // tests
					}
				}) // modes
			}
		}) // runtimeEnvs
	}
}

// mapFromKeySlices gives you a new map, with some keys already set to the value provided.
func mapFromKeySlices[K comparable, V any](value V, keySlices ...[]K) map[K]V {
	m := map[K]V{}

	for _, keySlice := range keySlices {
		for _, key := range keySlice {
			m[key] = value
		}
	}

	return m
}

func TestKubernetesScriptsBaseDir(t *testing.T) {
	t.Parallel()
	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	tests := map[string]struct {
		image    string
		shell    string
		script   string
		baseDir  string
		verifyFn func(t *testing.T, out string)
	}{
		"scripts_base_dir enabled": {
			image:   common.TestAlpineImage,
			shell:   "bash",
			script:  "find /tmp",
			baseDir: "/tmp",
			verifyFn: func(t *testing.T, out string) {
				assert.Regexp(t, regexp.MustCompile(`(?m)^/tmp/scripts-0-0$`), out)
			},
		},
		"scripts_base_dir trailing slash": {
			image:   common.TestAlpineImage,
			shell:   "bash",
			script:  "find /tmp",
			baseDir: "/tmp/",
			verifyFn: func(t *testing.T, out string) {
				assert.Regexp(t, regexp.MustCompile(`(?m)^/tmp/scripts-0-0$`), out)
			},
		},
		"scripts_base_dir disabled": {
			image:   common.TestAlpineImage,
			shell:   "bash",
			script:  "find / -maxdepth 1",
			baseDir: "",
			verifyFn: func(t *testing.T, out string) {
				assert.Regexp(t, regexp.MustCompile(`(?m)^/scripts-0-0$`), out)
			},
		},
	}

	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			t.Parallel()

			build := getTestBuild(t, func() (common.JobResponse, error) {
				return common.GetRemoteBuildResponse(tc.script)
			})

			build.Runner.RunnerSettings.Shell = tc.shell
			build.Runner.RunnerSettings.Kubernetes.ScriptsBaseDir = tc.baseDir
			build.JobResponse.Image.Name = tc.image
			build.Runner.Kubernetes.HelperImage = "registry.gitlab.com/gitlab-org/gitlab-runner/gitlab-runner-helper:x86_64-latest"

			var buf bytes.Buffer
			err := build.Run(&common.Config{}, &common.Trace{Writer: &buf})
			assert.NoError(t, err)

			tc.verifyFn(t, buf.String())
		})
	}
}

func TestKubernetesLogsBaseDir(t *testing.T) {
	t.Parallel()

	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	tests := map[string]struct {
		image    string
		shell    string
		script   string
		baseDir  string
		envVars  common.JobVariables
		verifyFn func(t *testing.T, out string)
	}{
		"logs_base_dir enabled": {
			image:   common.TestAlpineImage,
			shell:   "bash",
			script:  "find /tmp",
			baseDir: "/tmp",
			verifyFn: func(t *testing.T, out string) {
				assert.Regexp(t, regexp.MustCompile(`(?m)^/tmp/logs-0-0$`), out)
			},
		},
		"logs_base_dir trailing slash": {
			image:   common.TestAlpineImage,
			shell:   "bash",
			script:  "find /tmp",
			baseDir: "/tmp/",
			verifyFn: func(t *testing.T, out string) {
				assert.Regexp(t, regexp.MustCompile(`(?m)^/tmp/logs-0-0$`), out)
			},
		},
		"logs_base_dir disabled": {
			image:   common.TestAlpineImage,
			shell:   "bash",
			script:  "find / -maxdepth 1",
			baseDir: "",
			verifyFn: func(t *testing.T, out string) {
				assert.Regexp(t, regexp.MustCompile(`(?m)^/logs-0-0$`), out)
			},
		},
	}

	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			t.Parallel()

			build := getTestBuild(t, func() (common.JobResponse, error) {
				return common.GetRemoteBuildResponse(tc.script)
			})

			build.Runner.RunnerSettings.Shell = tc.shell
			build.Runner.RunnerSettings.Kubernetes.LogsBaseDir = tc.baseDir
			build.JobResponse.Image.Name = tc.image
			build.Runner.Kubernetes.HelperImage = "registry.gitlab.com/gitlab-org/gitlab-runner/gitlab-runner-helper:x86_64-latest"

			var buf bytes.Buffer
			err := build.Run(&common.Config{}, &common.Trace{Writer: &buf})
			assert.NoError(t, err)

			tc.verifyFn(t, buf.String())
		})
	}
}

func testJobAgainstServiceContainerBehaviour(t *testing.T, featureFlagName string, featureFlagValue bool) {
	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	tests := map[string]struct {
		services common.Services
		verifyFn func(t *testing.T, err error)
	}{
		"job fails when waiting for service port readiness and service fails": {
			services: common.Services{
				{
					Name: "postgres:12.17-alpine3.19",
					Variables: common.JobVariables{
						common.JobVariable{Key: "HEALTHCHECK_TCP_PORT", Value: "5432"},
					},
				},
			},
			verifyFn: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		// Postgres service container will fail because the password and database variables are not provided
		"job passes when service fails": {
			services: common.Services{
				{
					Name: "postgres:12.17-alpine3.19",
				},
			},
			verifyFn: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			t.Parallel()

			build := getTestBuild(t, func() (common.JobResponse, error) {
				return common.GetRemoteBuildResponse("sleep 5s")
			})

			buildtest.SetBuildFeatureFlag(build, featureFlagName, featureFlagValue)
			build.JobResponse.Image.Name = common.TestAlpineImage
			build.JobResponse.Services = append(build.JobResponse.Services, tc.services...)
			build.Runner.Kubernetes.HelperImage = "registry.gitlab.com/gitlab-org/gitlab-runner/gitlab-runner-helper:x86_64-latest"

			err := build.Run(&common.Config{}, &common.Trace{Writer: os.Stdout})
			tc.verifyFn(t, err)
		})
	}
}

func TestBuildContainerOOMKilled(t *testing.T) {
	t.Parallel()

	kubernetes.SkipKubectlIntegrationTests(t, "kubectl", "cluster-info")

	tests := map[string]struct {
		script      string
		verifyFn    func(t *testing.T, out string, err error)
		memoryLimit string
	}{
		"job fails because build container is OOMKilled": {
			script: `echo "Starting memory allocation to trigger OOM..."

allocate_memory() {
	while true; do
		# Allocate a large block of memory by creating a large variable
		data=$(printf 'A%.0s' $(seq 1 1000000)) # Adjust this number for more memory allocation
		sleep 1  # Optional: add a small delay to control the speed of allocation
	done
}

allocate_memory
`,
			memoryLimit: "6Mi",
			verifyFn: func(t *testing.T, out string, err error) {
				assert.Contains(t, out, "Error in container build: exit code: 137, reason: 'OOMKilled'")
				assert.Error(t, err)
			},
		},
	}

	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			t.Parallel()

			build := getTestBuild(t, func() (common.JobResponse, error) {
				return common.GetRemoteBuildResponse(tc.script)
			})

			buildtest.SetBuildFeatureFlag(build, featureflags.UseLegacyKubernetesExecutionStrategy, false)
			build.JobResponse.Image.Name = common.TestAlpineImage
			build.Runner.Kubernetes.HelperImage = "registry.gitlab.com/gitlab-org/gitlab-runner/gitlab-runner-helper:x86_64-latest"
			build.Runner.Kubernetes.MemoryLimit = tc.memoryLimit
			build.Runner.Kubernetes.MemoryRequest = tc.memoryLimit

			var buf bytes.Buffer
			err := build.Run(&common.Config{}, &common.Trace{Writer: &buf})
			tc.verifyFn(t, buf.String(), err)
		})
	}
}
