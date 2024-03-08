/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"drexel.edu/voter-api/pkg/create"
	"drexel.edu/voter-api/pkg/delete"
	"drexel.edu/voter-api/pkg/http/rest"
	"drexel.edu/voter-api/pkg/read"
	rediscache "drexel.edu/voter-api/pkg/storage/redis"
	"drexel.edu/voter-api/pkg/update"
	"github.com/spf13/cobra"
)

const (
	defaultPort = 3000
)

var port int
var redisLocation string

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "starts the server",
	Long: `Starts the server on port 3000. The user can define
	a different port using the -p --port flag. 
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("start called")

		redisCache, err := rediscache.New()
		if err != nil {
			fmt.Sprintln("Error initializing repository: ", err)
			panic(err)
		}

		createAdapter := create.New(redisCache)

		updateAdapter := update.New(redisCache)

		readAdapter := read.New(redisCache)

		deleteAdapter := delete.New(redisCache)

		router := rest.Handler(port, createAdapter, updateAdapter, readAdapter, deleteAdapter)

		msg := fmt.Sprintf("the server is started at: http://localhost:%d", port)

		fmt.Println(msg)

		log.Fatal(router.Listen(string(fmt.Sprintf(":%d", port))))
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	startCmd.Flags().IntVarP(&port, "port", "p", defaultPort, "The port voter-api will use.")
	//startCmd.Flags().StringVarP(&redisLocation, "redis", "r", defaultRedisLocation, "The redis location to use.")
}
