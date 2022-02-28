/*
Copyright © 2022 Noueman KHALIKINE <noueman.khal@gmail.com>

*/
package cmd

import (
	"github.com/NouemanKHAL/sugoku/pkg/config"
	"github.com/NouemanKHAL/sugoku/pkg/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	port     int
	logLevel string
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a sudoku server",
	Long: `
	Start a sudoku server listening on the given port (-p) or 7007 by default.
	Logging levels can be controlled with the (-l) flag. 
	List of supported logging levels (TRACE, DEBUG, INFO, WARN, ERROR).`,
	Run: func(cmd *cobra.Command, args []string) {
		configuration := config.Config{
			Port:     port,
			LogLevel: logLevel,
		}
		server.StartServer(configuration)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().IntVarP(&port, "port", "p", 7007, "The port number used by the server")
	viper.BindPFlag("port", startCmd.Flags().Lookup("port"))

	startCmd.Flags().StringVarP(&logLevel, "log-level", "ll", "INFO", "Logging level")
	viper.BindPFlag("log-level", startCmd.Flags().Lookup("log-level"))

	// TODO: add support for log file
}
