// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"log"
	"net"
	"github.com/melotusme/middleman"
)

var typ string
// proxyCmd represents the proxy command
var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(typ, " proxy called")

		log.SetFlags(log.LstdFlags | log.Lshortfile)
		l, err := net.Listen("tcp", ":8080")
		if err != nil {
			log.Panic(err)
		}
		for {
			client, err := l.Accept()
			if err != nil {
				log.Panic(err)
			}
			handleFunc := middleman.RequestHandleManager[typ]
			go handleFunc(client)
		}
	},
}

func init() {
	rootCmd.AddCommand(proxyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	proxyCmd.PersistentFlags().StringVar(&typ, "typ", "http", "proxy type")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	proxyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
