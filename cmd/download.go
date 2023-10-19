package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/DedicatedDev/merkle"
	"github.com/DedicatedDev/merkleclient/config"
	"github.com/spf13/cobra"
)

type DownloadResponse struct {
	Content     []byte       `json:"content"`
	MerkleProof merkle.Proof `json:"merkleProof"`
}

var fileID string

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download a file from the server",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := http.Get(fmt.Sprintf("%s/%s", config.ServerDownloadEndpoint, fileID))
		if err != nil {
			fmt.Println("Error downloading file:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading response body:", err)
				return
			}

			var downloadResponse DownloadResponse
			err = json.Unmarshal(bodyBytes, &downloadResponse)
			if err != nil {
				fmt.Println("Error unmarshalling response:", err)
				return
			}

			// Validate the downloaded file with Merkle proof
			rootHashBytes, err := ioutil.ReadFile("rootHash.txt")
			if err != nil {
				fmt.Println("Error reading saved root hash:", err)
				return
			}
			rootHash := string(rootHashBytes)
			fmt.Println("download Res:", downloadResponse.MerkleProof)

			isValid := downloadResponse.MerkleProof.Validate(rootHash)
			if !isValid {
				fmt.Println("The downloaded file failed Merkle proof validation!")
				return
			}
			// Save the file to disk
			err = ioutil.WriteFile(fileID, downloadResponse.Content, 0644)
			if err != nil {
				fmt.Println("Error saving file to disk:", err)
				return
			}

			fmt.Println("File downloaded and saved successfully!")
		} else {
			bodyBytes, _ := ioutil.ReadAll(resp.Body)
			fmt.Printf("Error: %s\n", bodyBytes)
		}
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().StringVarP(&fileID, "id", "i", "", "File ID to download (required)")
	downloadCmd.MarkFlagRequired("id")
}
