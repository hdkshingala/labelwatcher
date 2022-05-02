package server

import (
	"log"
	"net/http"
	"os"

	"github.com/hdkshingala/labelwatcher/controller"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apiserver/pkg/server"
	"k8s.io/apiserver/pkg/server/options"
	"k8s.io/component-base/cli/globalflag"
)

var (
	scheme = runtime.NewScheme()
	codec  = serializer.NewCodecFactory(scheme)
)

const (
	valDeploymentCreator = "valdeploymentcreator"
)

type Options struct {
	SecureServingOptions options.SecureServingOptions
}

func (options *Options) AddFlagSet(fs *pflag.FlagSet) {
	options.SecureServingOptions.AddFlags(fs)
}

func NewDefaultOptions() *Options {
	options := &Options{
		SecureServingOptions: *options.NewSecureServingOptions(),
	}

	options.SecureServingOptions.BindPort = 8443
	options.SecureServingOptions.ServerCert.PairName = valDeploymentCreator

	return options
}

type Config struct {
	SecureServingConfig *server.SecureServingInfo
}

func (options *Options) Config() (*Config, error) {
	err := options.SecureServingOptions.MaybeDefaultWithSelfSignedCerts("0.0.0.0", nil, nil)
	if err != nil {
		return nil, err
	}
	c := &Config{}
	options.SecureServingOptions.ApplyTo(&c.SecureServingConfig)

	return c, nil
}

type Server struct {
	Config     *Config
	Controller *controller.Controller
	Mux        http.ServeMux
}

func NewServer(controller *controller.Controller) (*Server, error) {
	options := NewDefaultOptions()

	fs := pflag.NewFlagSet(valDeploymentCreator, pflag.ExitOnError)
	globalflag.AddGlobalFlags(fs, valDeploymentCreator)

	options.AddFlagSet(fs)

	if err := fs.Parse(os.Args); err != nil {
		log.Printf("Failed to parse flags. Error: %s\n", err.Error())
		return nil, err
	}

	c, err := options.Config()
	if err != nil {
		log.Printf("Failed to prepare config. Error: %s\n", err.Error())
		return nil, err
	}

	ser := &Server{
		Config:     c,
		Controller: controller,
		Mux:        *http.NewServeMux(),
	}

	ser.Mux.HandleFunc("/", ser.ServeDeploymentCreatorValidation)
	return ser, nil
}
