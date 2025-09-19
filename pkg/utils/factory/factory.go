package factory

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	"gopkg.in/yaml.v3"

	"github.com/AlecAivazis/survey/v2"
	"github.com/berrybytes/zocli/internal/browser"
	"github.com/berrybytes/zocli/internal/config"
	"github.com/berrybytes/zocli/pkg/utils/iostreams"
	"github.com/berrybytes/zocli/pkg/utils/printer"
	"github.com/berrybytes/zocli/pkg/utils/terminal"
	cliBrowser "github.com/cli/browser"
	// "github.com/sirupsen/logrus"
)

// Factory
//
// all the root commands flags will be saved here,
// along with all other default fields
// and then, will be transferred to all the sub-command
type Factory struct {
	// determines the verbose mode of the CLI
	Verbose bool
	// determines the quiet mode of the CLI
	Quiet bool
	// determines the interactive mode of the CLI
	NoInteractive bool

	// determines if the user is logged in or not
	LoggedIn bool
	// configuration for running the CLI
	Config *config.Config
	// context for running the CLI
	Ctx context.Context
	// user's operating system
	UserOS string
	// user's home folder
	UserHomeFolder string
	// determines if the configuration folder has been created or not
	ConfigCreated bool
	// http client for sending requests
	HTTPClient http.Client
	// routes for sending requests to 01cloud-api-server
	Routes *config.Routes

	UserAuthToken string
	UserWebToken  string
	UserEmail     string
	WebTokenUsed  bool

	// default browser of the user
	UserBrowser *browser.Browser
	Promopter   survey.Prompt

	Term terminal.Provider
	IO   *iostreams.IOStreams

	Debug   printer.DebugInterface
	Printer printer.PrinterInterface
	Auth0Token string
}

type FactoryInterface interface {
	createDefaultConfig()
	CleanUp()
	CleanUpTest()
}

type MyConfig struct {
	Auth0Token string `yaml:"auth0_token"`
}
// New
//
// creates a new factory instance.
// which will contain all the default values for the necessary fields
func New(ctx context.Context, conf *config.Config) *Factory {
	f := &Factory{
		Config:     conf,
		Ctx:        ctx,
		UserOS:     runtime.GOOS,
		HTTPClient: http.Client{},
		Routes:     config.Load(),
		Term:       terminal.New(runtime.GOOS),
		IO:         iostreams.System(),
		Printer:    printer.New(),
	}

	if f.Config == nil {
		f.createDefaultConfig()
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// as there was error getting the user home directory
		// change the config folder path to temp
		// as nor then, none of the files can be created at runtime
		return f
	}
	f.UserHomeFolder = homeDir
	f.UserBrowser = browser.New("", cliBrowser.Stdout, cliBrowser.Stderr)
	f.Debug = printer.NewDebug(f.Config.Populated.Debug)

	return f
}

func (f *Factory) createDefaultConfig() {
	f.Config = config.New()
	f.Config.ConfigFolder = "/tmp"
	f.Config.AuthFile = "/auth.yaml"
}

func (f *Factory) CleanUp() {
	err := os.RemoveAll(f.Config.ConfigFolder)
	if err != nil {
		f.Printer.Errorf("Err: %s", err.Error())
	}
}

func (f *Factory) CleanUpTest() {
	err := os.RemoveAll(f.Config.ConfigFolder)
	if err != nil {
		f.Debug.Debug(err)
	}
}

// func (f *Factory) GetAuth() map[string]string {
// 	if f.UserWebToken != "" {
// 		header := map[string]string{"X-PERSONAL-TOKEN": f.UserWebToken}
// 		f.WebTokenUsed = true
// 		return header
// 	}
// 	f.WebTokenUsed = false
// 	header := map[string]string{"Authorization": "Basic " + f.UserAuthToken,"X-CUSTOM-AUTH": f.Auth0Token }
// 	return header
// }

// func GetAuth0Token() string {
// 	filename := os.Getenv("CONFIG_FILE")
// 	if filename == "" {
// 		log.Fatal("CONFIG_FILE environment variable is not set")
// 	}

// 	data, err := os.ReadFile(filename)
// 	if err != nil {
// 		log.Fatalf("error reading YAML file: %v", err)
// 	}

// 	var cfg MyConfig
// 	if err := yaml.Unmarshal(data, &cfg); err != nil {
// 		log.Fatalf("error unmarshalling YAML: %v", err)
// 	}
// 	log.Printf(cfg.Auth0Token)

// 	// returns just the value like: xxx
// 	return cfg.Auth0Token
// }

func GetAuth0Token() string {
    filename := os.Getenv("CONFIG_FILE")
    if filename == "" {
        log.Fatal("CONFIG_FILE environment variable is not set")
	}

    data, err := os.ReadFile(filename)
    if err != nil {
        log.Fatalf("error reading YAML file: %v", err)
    }

    var cfg MyConfig
    if err := yaml.Unmarshal(data, &cfg); err != nil {
        log.Fatalf("error unmarshalling YAML: %v", err)
    }

    // Only return the token string, not the whole struct or any timestamp
	var token string = cfg.Auth0Token

	fmt.Println(token)

    return token
}

func (f *Factory) GetAuth() map[string]string {
    if f.UserWebToken != "" {
        header := map[string]string{"X-PERSONAL-TOKEN": f.UserWebToken}
        f.WebTokenUsed = true
        return header
    }
    f.WebTokenUsed = false
	fmt.Printf("Auth0 Token: %s\n", GetAuth0Token())
    header := map[string]string{
        "Authorization": "Basic " + f.UserAuthToken,
        "X-CUSTOM-AUTH": GetAuth0Token(),    
		}
    return header
}
