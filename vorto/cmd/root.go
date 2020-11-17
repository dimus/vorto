/*
Copyright © 2020 Dmitry Mozzherin <dmozzherin@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/dimus/vorto"
	"github.com/dimus/vorto/config"
	"github.com/dimus/vorto/domain/entity"
	"github.com/spf13/cobra"

	"github.com/gnames/gnlib/sys"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const configTemplate = `# Datadir is a path to vorto's flash cards--
# words, terms, definitions.
DataDir: %s

# DefaultSet is a set of flash cards that will be used by default.
DefaultSet: default

# WordsBatch is the number of terms used in a training session.
WordsBatch: 25

# Sets is a list of flash card sets created by the user
# Supported types: "general", "esperanto"
Sets:
- default
`

var (
	opts []config.Option
)

// cnf purpose is to achieve automatic import of data from the
// configuration file, if it exists.
type configData struct {
	DataDir    string
	DefaultSet string
	WordsBatch int
	Sets       []string
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "vorto",
	Short: "A flash cards app to learn new words, facts or terms",
	Run: func(cmd *cobra.Command, args []string) {
		if showVersionFlag(cmd) {
			os.Exit(0)
		}
		var cardsStack *entity.CardStack
		c := config.NewConfig(opts...)
		vrt := vorto.NewVorto(c)
		defer func() {
			err := vrt.Save(cardsStack)
			if err != nil {
				fmt.Printf("Could not save progress: %s\n", err)
			}
		}()

		err := vrt.Init()
		if err != nil {
			log.Fatalf("Cannot initiate vorto: %s.", err)
		}

		cardsStack, err = vrt.Load()
		if err != nil {
			log.Fatalf("Cannot load cards: %s.", err)
		}

		vrt.Run(cardsStack)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().BoolP("version", "V", false, "show app's version")
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.vorto.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("vocabulary check", "v", false, "After learning words do vocabulary too")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	configFile := "vorto"

	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		log.Fatalf("Cannot find home directory: %s.", err)
	}
	configDir := filepath.Join(home, ".config")
	dataDir := filepath.Join(home, ".local", "share", "vorto")

	// Search config in home directory with name ".gnames" (without extension).
	viper.AddConfigPath(configDir)
	viper.SetConfigName(configFile)

	configPath := filepath.Join(configDir, fmt.Sprintf("%s.yaml", configFile))
	touchConfigFile(configPath, configFile, dataDir)

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Printf("Using config file: %s.\n", viper.ConfigFileUsed())
	}
	getOpts()
}

func getOpts() {
	cfg := &configData{}
	err := viper.Unmarshal(cfg)
	if err != nil {
		log.Fatalf("Cannot deserialize config data: %s.", err)
	}

	if cfg.DataDir != "" {
		opts = append(opts, config.OptDataDir(cfg.DataDir))
	}

	if cfg.DefaultSet != "" {
		opts = append(opts, config.OptDefaultSet(cfg.DefaultSet))
	}

	if cfg.WordsBatch != 0 {
		opts = append(opts, config.OptWordsBatch(cfg.WordsBatch))
	}

	if len(cfg.Sets) > 0 {
		opts = append(opts, config.OptSets(cfg.Sets))
	}
}

// touchConfigFile checks if config file exists, and if not, it gets created.
func touchConfigFile(configPath, configFile, dataDir string) {
	if sys.FileExists(configPath) {
		return
	}

	log.Printf("Creating config file: %s.", configPath)
	configText := fmt.Sprintf(configTemplate, dataDir)
	createConfig(configPath, configFile, configText)
}

// createConfig creates config file.
func createConfig(path, file, configText string) {
	err := sys.MakeDir(filepath.Dir(path))
	if err != nil {
		log.Fatalf("Cannot create dir %s: %s.", path, err)
	}

	err = ioutil.WriteFile(path, []byte(configText), 0644)
	if err != nil {
		log.Fatalf("Cannot write to file %s: %s.", path, err)
	}
}

// showVersionFlag provides version and the build timestamp. If it returns
// true, it means that version flag was given.
func showVersionFlag(cmd *cobra.Command) bool {
	hasVersionFlag, err := cmd.Flags().GetBool("version")
	if err != nil {
		log.Fatalf("Cannot get version flag: %s.", err)
	}

	if hasVersionFlag {
		fmt.Printf("\nversion: %s\nbuild: %s\n\n", vorto.Version, vorto.Build)
	}
	return hasVersionFlag
}
