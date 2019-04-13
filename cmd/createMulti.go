package cmd

import (
	sadBoy "encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	astar "github.com/beefsack/go-astar"
	"github.com/spf13/cobra"
)

var threadCount int

// createMultiCmd represents the createMulti command
var createMultiCmd = &cobra.Command{
	Use:   "createMulti",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		t1 := time.Now()
		fmt.Println("create called")
		var Counter int
		var export bool
		var pathChannel chan []string = make(chan []string)
		var endChannel chan bool = make(chan bool)
		guard := make(chan struct{}, threadCount)
		go writeToCSV(pathChannel, endChannel)
		if csv != "" && strings.HasSuffix(csv, ".csv") {
			export = true
			fmt.Println(csv)
		}

		rand.Seed(time.Now().UnixNano())

		if export {
			// file, err := os.OpenFile(csv, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
			// checkError("Cannot create file", err)
			// defer file.Close()

			// writer = sadBoy.NewWriter(file)
			// defer writer.Flush()
		} else {
			panic(fmt.Errorf("hell being mell"))
		}

		for cycle := 0; cycle < maxCycles; cycle++ {
			guard <- struct{}{}
			go func() {
				var fromPos, toPos alphaPos
				var worldStr string

				Counter++
				if anySide {
					fromPos = alphaPos{
						x: rand.Intn(cols),
						y: rand.Intn(rows),
					}

					toPos = alphaPos{
						x: rand.Intn(cols),
						y: rand.Intn(rows),
					}

					if Side := rand.Intn(2); Side == 0 {
						fromPos.x = 0
					} else {
						fromPos.y = 0
					}

					if Side := rand.Intn(2); Side == 0 {
						toPos.x = cols - 1
					} else {
						toPos.y = rows - 1
					}
				} else {
					fromPos = alphaPos{
						x: 0,
						y: rand.Intn(rows),
					}

					toPos = alphaPos{
						x: cols - 1,
						y: rand.Intn(rows),
					}
				}

				for r := 0; r < rows; r++ {
					for c := 0; c < cols; c++ {
						if c == fromPos.x && r == fromPos.y {
							worldStr = fmt.Sprintf("%sF", worldStr)
						} else if c == toPos.x && r == toPos.y {
							worldStr = fmt.Sprintf("%sT", worldStr)
						} else {
							worldStr = fmt.Sprintf("%s%v", worldStr, rand.Intn(2))
						}
					}
					worldStr = fmt.Sprintf("%s\n", worldStr)
				}
				world := ParseWorld(worldStr)
				p, dist, found := astar.Path(world.From(), world.To())
				// spew.Dump(p)
				if found {
					fmt.Printf("Found in distance = %v in %v cycle\n", dist, cycle)
					fmt.Println("------ Printing World ------")
					fmt.Println(world.RenderPath(p))
					if export {
						fmt.Println("-------- EXPORT:CSV --------")
						pathChannel <- world.PrepPath(p)

					} else {
						fmt.Println("----------------------------")
					}

				}
				<-guard
			}()

		}
		pathChannel <- []string{"kill"}
		<-endChannel
		timeTook := time.Now().Sub(t1)
		fmt.Printf("Completed %v cycles in %s\n", maxCycles, timeTook)
	},
}

func writeToCSV(pathChannel chan []string, endChannel chan bool) {
	file, err := os.OpenFile(csv, os.O_CREATE|os.O_WRONLY, 0777)
	checkError("Cannot create file", err)
	defer file.Close()

	writer := sadBoy.NewWriter(file)
	defer writer.Flush()

	for {
		pPath := <-pathChannel
		if pPath[0] == "kill" {
			break
		}
		err := writer.Write(pPath)
		checkError("Cannot create file", err)
	}

	endChannel <- true
}

func init() {
	rootCmd.AddCommand(createMultiCmd)
	createMultiCmd.Flags().IntVarP(&rows, "rows", "r", 12, "Set Rows")
	createMultiCmd.Flags().IntVarP(&cols, "cols", "c", 12, "Set Columns")
	createMultiCmd.Flags().BoolVarP(&anySide, "any", "a", false, "Let path run along any side")
	createMultiCmd.Flags().IntVarP(&maxCycles, "mcycles", "m", 500, "Max cycles to generate paths")
	createMultiCmd.Flags().StringVar(&csv, "csv", "", "Save to csv")
	createMultiCmd.Flags().IntVarP(&threadCount, "threads", "t", 1, "Choose the amount of threads")
}
