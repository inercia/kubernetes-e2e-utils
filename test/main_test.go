package test

import (
	"flag"
	"os"
	"testing"

	"k8s.io/klog/v2"
	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/envfuncs"

	lenvfuncs "github.com/inercia/kubernetes-e2e-utils/pkg/envfuncs"
)

var (
	testenv env.Environment
)

func TestMain(m *testing.M) {
	klog.InitFlags(flag.CommandLine) // initializing the flags
	defer klog.Flush()               // flushes all pending log I/O

	flag.Set("v", "5")

	testenv = env.New()

	clusterName := "kubernetes-e2e-utils-tests"

	namespace := envconf.RandomName("e2e", 16)

	testenv.Setup(
		lenvfuncs.CreateK3dCluster(clusterName),
		envfuncs.CreateNamespace(namespace),
	)

	testenv.Finish(
		envfuncs.DeleteNamespace(namespace),
		lenvfuncs.DestroyK3dCluster(clusterName),
	)

	// launch package tests
	os.Exit(testenv.Run(m))
}
