package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"little-auth/models"
	"little-auth/utils"
	"little-auth/vault"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove an authenticating entry.",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			chosenIndexer *models.Indexer
			mapStrIndexer map[string]*models.Indexer
			secretVault   vault.Vault
			err           error
		)
		if secretVault, err = vault.New(); err != nil {
			fmt.Printf("Failed to initialize vault: %v\n", err)
			panic(err)
		}

		indexes, err := secretVault.GetAllIndexes()
		if err != nil {
			fmt.Printf("Failed to get list of authenticate info: %v\n", err)
			panic(err)
		}

		if len(indexes) == 0 {
			fmt.Println("No entries found.")
			return
		}

		var indexStrs []string
		mapStrIndexer = make(map[string]*models.Indexer)
		for _, index := range indexes {
			name := ""
			if index.Issuer != "" && index.Path != "" {
				name = fmt.Sprintf("Path: %s - Issuer: %s", index.Path, index.Issuer)
			} else if index.Issuer == "" {
				name = fmt.Sprintf("Path: %s", index.Path)
			} else if index.Path == "" {
				name = fmt.Sprintf("Issuer: %s", index.Issuer)
			}
			if name != "" {
				indexStrs = append(indexStrs, name)
				mapStrIndexer[name] = index
			}
		}
		chosenIndexerStr := utils.PromptGetSelect(utils.PromptContent{
			ErrorMsg: "Please choose one to get OTP.",
			Label:    "Select an entry.",
		}, indexStrs)
		if chosenIndexer = mapStrIndexer[chosenIndexerStr]; chosenIndexer == nil {
			fmt.Println("Failed to get chosen indexer.")
			panic("failed to get chosen indexer")
		}

		if err := secretVault.Delete(chosenIndexer); err != nil {
			fmt.Printf("Failed to delete entry. Please try again. Err: %v\n", err)
			panic(err)
		}

		fmt.Printf("Successfully deleted entry:\n\tPath: %s\n\tIssuer: %s\n", chosenIndexer.Path, chosenIndexer.Issuer)
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
