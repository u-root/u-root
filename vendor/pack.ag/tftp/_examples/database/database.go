// Copyright (C) 2017 Kale Blankenship. All rights reserved.
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details

package main

import (
	"database/sql"
	"io/ioutil"
	"log"

	"pack.ag/tftp"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Create or open a sqlite database
	db, err := sql.Open("sqlite3", "tftp.db")
	if err != nil {
		log.Fatal(err)
	}

	// Create a simple table to hold the ip and sent log data from
	// the client.
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tftplogs (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        ip TEXT,
        log TEXT
    );`)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new server listening on port 6900, all interfaces
	server, err := tftp.NewServer(":6900")
	if err != nil {
		log.Fatal(err)
	}

	// Set the server's write handler, read requests will be rejeccted
	server.WriteHandler(&tftpDB{db})

	// Start the server, if it fails error will be printed by log.Fatal
	log.Fatal(server.ListenAndServe())
}

// tftpDB embeds a *sql.DB and implements the tftp.ReadHandler
// interface.
type tftpDB struct {
	*sql.DB
}

func (db *tftpDB) ReceiveTFTP(w tftp.WriteRequest) {
	// Get the file size
	size, err := w.Size()

	// We're choosing to only store logs that are less than 1MB.
	// An error indicates no size was received.
	if err != nil || size > 1024*1024 {
		// Send a "disk full" error.
		w.WriteError(tftp.ErrCodeDiskFull, "File too large or no size sent")
		return
	}

	// Note: The size value is sent by the client, the client could send more data than
	// it indicated in the size option. To be safe we'd want to allocate a buffer
	// with the size we're expecting and use w.Read(buf) rather than ioutil.ReadAll.

	// Read the data from the client into memory
	data, err := ioutil.ReadAll(w)
	if err != nil {
		log.Println(err)
		return
	}

	// Insert the IP address of the client and the data into the database
	res, err := db.Exec("INSERT INTO tftplogs (ip, log) VALUES (?, ?)", w.Addr().IP.String(), string(data))
	if err != nil {
		log.Println(err)
		return
	}

	// Log a message with the details
	id, _ := res.LastInsertId()
	log.Printf("Inserted %d bytes of data from %s. (ID=%d)", len(data), w.Addr().IP, id)
}
