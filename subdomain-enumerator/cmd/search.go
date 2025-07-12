package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/cheynewallace/tabby"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"log"
	"net"
	"os"
	"subdomain-enumerator/dnsStuff"
	"sync"
)

var showProgress bool

var searchCmd = &cobra.Command{
	Use:   "search <domain>",
	Short: "perform a search for a domain",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		domain := args[0]

		if !viper.IsSet("wordlist") {
			return fmt.Errorf("wordlist missing")
		}

		displayInfo()

		wordlistFile, err := openWordlist()
		if err != nil {
			return err
		}
		defer func() { _ = wordlistFile.Close() }()

		//_, err = countLines(wordlistFile)
		if err != nil {
			return fmt.Errorf("error counting wordlist entries: %v", err)
		}

		resultsChan := make(chan []dnsStuff.Result)
		performSearch(domain, wordlistFile, resultsChan)

		// output result
		if showProgress {
			displayResultWithProgress(wordlistFile, resultsChan)
			return nil
		}

		displayResult(resultsChan)
		return nil
	},
}

func init() {
	searchCmd.Flags().StringP("wordlist", "w", "", "wordlist to use")
	searchCmd.Flags().StringP("server", "s", "", "dns server to use")
	searchCmd.Flags().IntP("workers", "n", 0, "amount of workers to use")
	searchCmd.Flags().BoolVarP(&showProgress, "progress", "p", false, "show progress")

	cobra.CheckErr(
		viper.BindPFlag("wordlist", searchCmd.Flags().Lookup("wordlist")))
	cobra.CheckErr(
		viper.BindPFlag("server", searchCmd.Flags().Lookup("server")))
	cobra.CheckErr(
		viper.BindPFlag("workers", searchCmd.Flags().Lookup("workers")))
}

func displayInfo() {
	searchInfo := tabby.New()
	searchInfo.AddHeader(fmt.Sprintf("Performing search on: %s", domain))
	for key, value := range viper.AllSettings() {
		searchInfo.AddLine(key, value)
	}
	searchInfo.Print()
}

func openWordlist() (*os.File, error) {
	wordList := viper.GetString("wordlist")
	wordlistFile, err := os.Open(wordList)
	if err != nil {
		return nil, fmt.Errorf("error opening wordlist: %v", err)
	}
	return wordlistFile, nil
}

func performSearch(domain string, wordlistFile *os.File, resultsChan chan []dnsStuff.Result) {
	workerCount := viper.GetInt("workers")
	serverAddr := viper.GetString("server")
	fqdnChan := make(chan string)
	var wg sync.WaitGroup

	for w := 0; w < workerCount; w++ {
		wg.Add(1)

		conn, err := net.Dial("udp", serverAddr)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error setting up connection: %v", err)
		}

		go worker(fqdnChan, resultsChan, &wg, conn)
		log.Printf("worker %d using %s", w, serverAddr) // debug
	}

	// add fqdns to test
	go func() {
		defer close(fqdnChan)

		scanner := bufio.NewScanner(wordlistFile)
		for scanner.Scan() {
			fqdnChan <- fmt.Sprintf("%s.%s", scanner.Text(), domain)
		}
	}()

	go func() {
		wg.Wait()
		close(resultsChan)
	}()
}

func worker(fqdnChan <-chan string, resultsChan chan<- []dnsStuff.Result, wg *sync.WaitGroup, conn net.Conn) {
	for fqdn := range fqdnChan {
		results, err := dnsStuff.Lookup(fqdn, conn)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		}

		resultsChan <- results
	}
	wg.Done()
}

func displayResultWithProgress(wordlistFile *os.File, resultsChan chan []dnsStuff.Result) {
	fmt.Printf("\nResults\n----------------\n")
	totalLines, err := countLines(wordlistFile)
	if err != nil {
		fmt.Printf("Unable to count total lines")
	}

	currentLine := 0
	for results := range resultsChan {
		currentLine++
		fmt.Printf("\r")

		for _, result := range results {
			fmt.Printf("%-15s  <-  %s\n", result.IPAddress, result.Hostname)
		}
		fmt.Printf("%d/%d", currentLine, totalLines)
	}
}

func displayResult(resultsChan chan []dnsStuff.Result) {
	fmt.Printf("\nResults\n----------------\n")
	for results := range resultsChan {
		for _, result := range results {
			fmt.Printf("%-15s  <-  %s\n", result.IPAddress, result.Hostname)
		}
	}
}

func countLines(file *os.File) (int, error) {
	buf := make([]byte, 32*1024) // 32KB buffer
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := file.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			_, _ = file.Seek(0, 0)
			return count, nil
		case err != nil:
			return 0, err
		}
	}
}
