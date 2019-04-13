package cmd

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	sadBoy "encoding/csv"

	astar "github.com/beefsack/go-astar"
	"github.com/spf13/cobra"
)

type alphaPos struct {
	x int
	y int
}

var rows int
var cols int
var anySide bool
var csv string
var maxCycles int

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create train/test data",
	Long:  `Create training data and split into train/test sets`,
	Run: func(cmd *cobra.Command, args []string) {
		t1 := time.Now()
		fmt.Println("create called")
		var Counter int
		var export bool
		var writer *sadBoy.Writer
		if csv != "" && strings.HasSuffix(csv, ".csv") {
			export = true
			fmt.Println(csv)
		}

		rand.Seed(time.Now().UnixNano())

		if export {
			file, err := os.OpenFile(csv, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
			checkError("Cannot create file", err)
			defer file.Close()

			writer = sadBoy.NewWriter(file)
			defer writer.Flush()
		}

		for cycle := 0; cycle < maxCycles; cycle++ {
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
				fmt.Printf("Found in distance = %v in %v cycle\n", dist, Counter)
				fmt.Println("------ Printing World ------")
				fmt.Println(world.RenderPath(p))
				if export {
					fmt.Println("-------- EXPORT:CSV --------")
					err := writer.Write(world.PrepPath(p))
					checkError("Cannot create file", err)
				} else {
					fmt.Println("----------------------------")
				}

			}
		}

		timeTook := time.Now().Sub(t1)
		fmt.Printf("Completed %v cycles in %s\n", maxCycles, timeTook)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().IntVarP(&rows, "rows", "r", 12, "Set Rows")
	createCmd.Flags().IntVarP(&cols, "cols", "c", 12, "Set Columns")
	createCmd.Flags().BoolVarP(&anySide, "any", "a", false, "Let path run along any side")
	createCmd.Flags().IntVarP(&maxCycles, "mcycles", "m", 500, "Max cycles to generate paths")
	createCmd.Flags().StringVar(&csv, "csv", "", "Save to csv")
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
