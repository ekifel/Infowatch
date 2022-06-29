package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ekifel/Infowatch/internal/file_generator"
	wp "github.com/ekifel/Infowatch/internal/workerpool"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func main() {
	numberOfFiles := 500
	numberOfWorkers := 50
	path := "demo_files"
	totalResult := map[string]int{}

	fmt.Printf("Starting to generate files...\n")

	g := file_generator.NewGenerator(path)
	err := g.CreateFiles(numberOfFiles)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Done\n")

	jobBatch := []wp.Job{}

	files, _ := os.ReadDir(path)
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			log.Printf("Can't read info from file  %s: %v", path, err)
		}

		jobBatch = append(jobBatch, wp.Job{
			DirPath:  path,
			FileInfo: info,
		})
	}

	fmt.Printf("Starting to check files...\n")

	pool := wp.NewWorkerPool(numberOfWorkers, numberOfFiles)
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	go pool.Run(ctx)
	pool.GenerateFrom(jobBatch)

	for {
		select {
		case r, ok := <-pool.Results():
			if !ok {
				continue
			}

			for k, v := range r {
				if num, ok := totalResult[k]; !ok {
					totalResult[k] = 1
				} else {
					totalResult[k] = num + v
				}
			}

		case <-pool.Done:
			g.CleanFiles()

			if err := ui.Init(); err != nil {
				log.Fatalf("failed to initialize termui: %v", err)
			}
			defer ui.Close()

			labels := []string{}
			data := []float64{}
			for k, v := range totalResult {
				if int([]rune(k)[0]) >= 32 {
					if k != " " {
						labels = append(labels, k)
					} else {
						labels = append(labels, "Space")
					}

					data = append(data, float64(v))
				}
			}

			bc := widgets.NewBarChart()
			bc.Title = "Bar Chart"
			bc.SetRect(0, 0, 125, 10)
			bc.BarColors = []ui.Color{ui.ColorGreen}
			bc.NumStyles = []ui.Style{ui.NewStyle(ui.ColorBlack)}
			bc.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorRed)}

			draw := func(count int) {
				bc.Data = data[count/2%10:]
				bc.Labels = labels[count/2%10:]

				ui.Render(bc)
			}

			tickerCount := 1
			draw(tickerCount)
			tickerCount++
			uiEvents := ui.PollEvents()
			ticker := time.NewTicker(time.Second).C
			for {
				select {
				case e := <-uiEvents:
					switch e.ID {
					case "q", "<C-c>":
						return
					}
				case <-ticker:
					draw(tickerCount)
					tickerCount++
				}
			}

		}
	}

}
