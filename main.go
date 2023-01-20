package main

import (
	"github.com/rivo/tview"
	"net"
)

var running = false
var quit = make(chan bool)

func main() {
	veil := tview.NewApplication().EnableMouse(true)
	pages := tview.NewPages()

	// create the intercept page
	intercept_page := tview.NewFlex().SetDirection(tview.FlexColumn)

	// -----------------------------------------------------------------
	// -----   create the different columns -----------
	// -----------------------------------------------------------------
	first_column := tview.NewFlex().SetDirection(tview.FlexRow)
	intercept_column := tview.NewFlex().SetDirection(tview.FlexColumn)
	// intercept_features := tview.NewFlex().SetDirection(tview.FlexRow)

	//

	// -----------------------------------------------------------------
	// ---- add request and response to the intercept column ---------
	// -----------------------------------------------------------------
	captured_request := tview.NewTextArea().SetTitle(" Request ").SetTitleAlign(tview.AlignLeft).SetBorder(true)
	server_response := tview.NewTextArea().SetTitle(" Response ").SetTitleAlign(tview.AlignLeft).SetBorder(true)
	intercept_column.AddItem(captured_request, 0, 1, true)
	intercept_column.AddItem(server_response, 0, 1, true)

	// -----------------------------------------------------------------
	// --- add features to the first column
	// -----------------------------------------------------------------

	set_options := tview.NewForm()
	set_options.AddDropDown("Intercept", []string{"ON", "OFF"}, 0, nil)
	set_options.AddDropDown("Request Method", []string{"GET", "POST", "HEAD", "PUT", "DELETE", "TRACE"}, 0, nil)

	first_column.AddItem(set_options, 0, 1, true)

	// -----------------------------------------------------------------
	// --- add the columns to the intercept page --------
	// -----------------------------------------------------------------
	intercept_page.AddItem(first_column, 0, 1, true)
	intercept_page.AddItem(intercept_column, 0, 6, true)
	// intercept_page.AddItem(intercept_features, 0, 1, true)
	intercept_page.SetBorder(true).SetTitle(" Veil Proxy ")

	pages.AddPage("intercept_page", intercept_page, true, true)

	Intercept("on")
	veil.SetRoot(pages, true).Run()
}

func Intercept(intercept string) {
	if intercept == "on" {
		if !running {
			go func() {
				proxy, _ := net.Listen("tcp", ":8081")
				for {
					stop := <-quit
					if stop == true {
						return
					}
					proxy.Accept()
				}
			}()
			running = true
		}
	}

	if intercept == "on" {
		if running {
			quit <- true
		}
		running = false
	}
}

func Update_Request_box(){

}
