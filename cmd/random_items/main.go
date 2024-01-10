// package main

// import (
// 	"web/internal/repo"

// 	_ "web/docs"
// )

// func main() {

// 	db, _ := repo.NewRepo()
// 	db.DB.Save()
// }

package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"web/internal/models"
	"web/internal/repo"
)

func main() {
	file, err := os.Open("cmd/random_items/items.csv")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	db, _ := repo.NewRepo()
	items := make([]*models.Item, 0)
	for i, v := range records {
		if i == 0 {
			continue
		}
		// if i == 2 {
		// 	break
		// }
		item := &models.Item{
			ItemPrototype: models.ItemPrototype{
				Title:    v[0],
				Subtitle: v[2],
				Price:    v[3],
				ImgURL:   v[4],
				Info:     v[5],
				Type:     v[1],
			},
		}
		items = append(items, item)
		if len(items) == 100 {
			tx := db.DB.Create(items)
			if tx.Error != nil {
				fmt.Println(tx.Error)
			}
			items = make([]*models.Item, 0)
		}
	}

	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// fmt.Println()
	// db.DB.Save(items)
	// 	if tx.Error != nil {
	// 		fmt.Println(tx.Error)
	// 	}
	// }

}
