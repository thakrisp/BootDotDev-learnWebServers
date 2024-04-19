package database

import "errors"

type Chirp struct {
	ID       int    `json:"id"`
	Body     string `json:"body"`
	AuthorID int    `json:"author_id"`
}

func (db *DB) CreateChirp(body string, authorID int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	id := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		ID:       id,
		Body:     body,
		AuthorID: authorID,
	}
	dbStructure.Chirps[id] = chirp

	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

func (db *DB) GetChirp(id int) (Chirp, error) {
	dbStructure, err := db.loadDB()

	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := dbStructure.Chirps[id]

	if !ok {
		return Chirp{}, ErrNotExist
	}

	return chirp, nil
}

func (db *DB) DeleteChirp(chirpID, userID int) (bool, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return false, err
	}

	chrips := dbStructure.Chirps

	chirp, ok := chrips[chirpID]
	if !ok {
		return false, nil
	}

	if chirp.AuthorID != userID {
		return false, errors.New("chirp ID doesnt match user ID")
	}

	delete(chrips, chirpID)
	dbStructure.Chirps = chrips

	err = db.writeDB(dbStructure)
	if err != nil {
		return false, errors.New("couldn't write to DB")
	}

	return true, nil
}
