package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	//"os"
	"github.com/DedicatedDev/merkle"
	"github.com/DedicatedDev/merkleclient/config"
	"github.com/spf13/cobra"
)

var filePaths []string

type UploadResponse struct {
	FileIds []string `json:"fileIds"`
}

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload a file to the server",
	Run: func(cmd *cobra.Command, args []string) {
		var fileContents [][]byte
		for _, path := range filePaths {
			content, err := ioutil.ReadFile(path)
			if err != nil {
				fmt.Printf("Error reading file %s: %v\n", path, err)
				return
			}
			fileContents = append(fileContents, content)
		}

		tree := merkle.NewMerkleTree(fileContents)
		rootHash := tree.Root.Hash

		// Save the rootHash locally. This could be saved in a file or a database.
		ioutil.WriteFile("rootHash.txt", []byte(rootHash), 0644)

		requestBody, err := json.Marshal(map[string][][]byte{
			"files": fileContents,
		})
		if err != nil {
			fmt.Printf("Error marshaling JSON: %v\n", err)
			return
		}

		// Send the request
		resp, err := http.Post(config.ServerUploadEndpoint, "application/json", bytes.NewReader(requestBody))
		if err != nil {
			fmt.Println("Error uploading file:", err)
			return
		}
		defer resp.Body.Close()
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Error: %s\n", bodyBytes)
		} else {
			var uploadResponse UploadResponse
			err = json.Unmarshal(bodyBytes, &uploadResponse)
			if err != nil {
				fmt.Printf("Error unmarshalling response: %v\n", err)
				return
			}
			fmt.Println(string(bodyBytes))
			fmt.Printf("File uploaded successfully: %s\n", uploadResponse.FileIds)
		}

		fmt.Println("All files uploaded successfully!")
		// Optionally delete local files
		for _, path := range filePaths {
			os.Remove(path)
		}
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
	uploadCmd.Flags().StringSliceVarP(&filePaths, "files", "f", []string{}, "List of file paths to upload (required)")
	uploadCmd.MarkFlagRequired("files")
}
