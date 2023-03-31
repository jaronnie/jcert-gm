/*
Copyright © 2023 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

/*
	根据场景生成所有证书文件
*/

// scopeCmd represents the scope command
var scopeCmd = &cobra.Command{
	Use:   "scope",
	Short: "build scope",
	Long:  `build scope`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("scope called")
	},
}

func init() {
	rootCmd.AddCommand(scopeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scopeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scopeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
