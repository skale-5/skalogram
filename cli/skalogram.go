package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/qeesung/image2ascii/convert"
)

type Database struct {
	dbPath       string
	ImageUpvotes map[string]int
}

func (db *Database) Upvote(imagePath string) int {
	votes := db.ImageUpvotes[imagePath]
	db.ImageUpvotes[imagePath] = votes + 1
	return db.ImageUpvotes[imagePath]
}

func (db *Database) Downvote(imagePath string) int {
	votes := db.ImageUpvotes[imagePath]
	db.ImageUpvotes[imagePath] = votes - 1
	return db.ImageUpvotes[imagePath]
}

func (db *Database) Save() error {
	// Transform in memory data to json representation
	data, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		return err
	}

	// write data into file
	err = os.WriteFile(db.dbPath, data, 0664)
	if err != nil {
		return err
	}
	return nil
}

func NewDatabase(dbPath string) *Database {
	// Open database file
	dbFile, err := os.OpenFile(dbPath, os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer dbFile.Close()

	// Create empty in memory database
	db := &Database{
		dbPath:       dbPath,
		ImageUpvotes: map[string]int{},
	}

	// Load db file into memory database
	dec := json.NewDecoder(dbFile)
	err = dec.Decode(db)
	if err != nil && err != io.EOF {
		panic(err)
	}

	return db
}

func main() {
	// Parse command line arguments
	imagePath := flag.String("image", "", "path to the image you want to print (required)")
	dbPath := flag.String("db-path", "./.db.json", "path to the skalogram database")
	flag.Parse()

	if *imagePath == "" {
		flag.PrintDefaults()
		return
	}

	// Load database into memory
	db := NewDatabase(*dbPath)

	// Create Image Converter and open image
	converter := convert.NewImageConverter()
	img, err := convert.OpenImageFile(*imagePath)
	if err != nil {
		panic(err)
	}

	// Convert image to ASCII and Print it
	ascii := converter.Image2ASCIIString(img, &convert.DefaultOptions)
	fmt.Print(ascii)

	// Ask user if he likes the image
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Do you like this image ? (Y/N)")
	var imageVotes int
	for scanner.Scan() {
		switch scanner.Text() {
		case "Y":
			imageVotes = db.Upvote(*imagePath)
		case "N":
			imageVotes = db.Downvote(*imagePath)
		default:
			fmt.Println("Wrong choice. Do you like this image ? (Y/N)")
			continue
		}
		break
	}

	// Display image total votes
	fmt.Println("The image", *imagePath, "counts", imageVotes, "votes !")

	// Save in memory data to db file
	err = db.Save()
	if err != nil {
		panic(err)
	}
}
