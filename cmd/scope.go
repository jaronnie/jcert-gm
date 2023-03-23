/*
Copyright © 2023 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

/*
	根据场景生成 hyperchain 的所有证书文件
	1. 创建联盟链
	2. 创建命名空间
	3. 加入联盟链
	4. 增加 vp 节点
	5. 增加 nvp 节点
*/

// scopeCmd represents the scope command
var scopeCmd = &cobra.Command{
	Use:   "scope",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
