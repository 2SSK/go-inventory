package main

import (
	"encoding/json"
	"fmt"
	"github.com/rivo/tview"
	"log"
	"os"
	"strconv"
)

// Define a item structure that will hold the stock information
type Item struct {
	Name  string `json:"name"`
	Stock int    `json:"stock"`
}

var (
	inventory     = []Item{}
	inventoryFile = "inventory.json"
)

// (function):1 - Load inventory from json file
func loadInventory() {
	if _, err := os.Stat(inventoryFile); err != nil {
		// Checking if the file exists
		data, err := os.ReadFile(inventoryFile)
		if err != nil {
			log.Fatal("Error reading inventory  file! - ", err)
		}
		// decode JSON into the inventory slice
		json.Unmarshal(data, &inventory)
	}
}

// (function):2 - Save inventory
func saveInventory() {
	data, err := json.MarshalIndent(inventory, "", " ")
	if err != nil {
		log.Fatal("Error saving inventory! - ", err)
	}
	// save the data in JSON file (`inventory.json`)
	os.WriteFile(inventoryFile, data, 0644)
}

// (function):3 - Delete item function
func deleteItem(index int) {
	// verify the item index
	if index < 0 || index >= len(inventory) {
		fmt.Println("Invalid item index!")
		return
	}
	// delete target item : (join all the previous elements of target with all the next elements)
	inventory = append(inventory[:index], inventory[index+1:]...)
	saveInventory() // save the changes
}

func main() {
	// Create a new TUI app
	app := tview.NewApplication()
	loadInventory()
	inventoryList := tview.NewTextView().SetDynamicColors(true).SetWordWrap(true)
	inventoryList.SetBorder(true).SetTitle("No items in inventory.")

	refreshInventory := func() {
		inventoryList.Clear()
		if len(inventory) == 0 {
			fmt.Fprintln(inventoryList, "No items in inventory.")
		} else {
			for i, item := range inventory {
				fmt.Fprintf(inventoryList, "[%d] %s (Stock: %d)\n", i+1, item.Name, item.Stock)
			}
		}
	}

	// Creating three input fields
	itemNameInput := tview.NewInputField().SetLabel("Item Name: ")
	itemStockInput := tview.NewInputField().SetLabel("Stock: ")
	itemIDInput := tview.NewInputField().SetLabel("Item ID to delete: ")

	form := tview.NewForm().AddFormItem(itemNameInput).AddFormItem(itemStockInput).AddFormItem(itemIDInput).AddButton("Add Item", func() {
		name := itemNameInput.GetText()
		stock := itemStockInput.GetText()
		if name != "" && stock != "" {
			quantity, err := strconv.Atoi(stock)
			if err != nil {
				fmt.Fprintf(inventoryList, "Invalid stock value.")
				return
			}
			inventory = append(inventory, Item{Name: name, Stock: quantity})
			saveInventory()
			refreshInventory()
			itemNameInput.SetText("")
			itemStockInput.SetText("")
		}
	}).AddButton("Delete Item", func() {
		idStr := itemIDInput.GetText()
		if idStr == "" {
			fmt.Fprintln(inventoryList, "Please enter an item ID to delete.")
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil || id < 1 || id > len(inventory) {
			fmt.Fprintln(inventoryList, "Invalid item ID.")
			return
		}
		deleteItem(id - 1)
		fmt.Fprintf(inventoryList, "Item [%d] deleted.\n", id)
		refreshInventory()
		itemIDInput.SetText("")
	}).AddButton("Exit", func() {
		app.Stop()
	})

	form.SetBorder(true).SetTitle("Manage Inventory").SetTitleAlign(tview.AlignLeft)

	flex := tview.NewFlex().AddItem(inventoryList, 0, 1, false).AddItem(form, 0, 1, true)

	refreshInventory()

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
