package database

import (
	"encoding/json"
	"errors"
	"os"
	"sort"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type DB struct {
	path string
	mu   *sync.RWMutex
}

type DBStructure struct {
	Chirps        map[int]Chirp  `json:"chirps"`
	Users         map[int]User   `json:"users"`
	RevokedTokens map[int]string `json:"revoked_tokens"`
}

type Chirp struct {
	ID       int    `json:"id"`
	AuthorID int    `json:"author_id"`
	Body     string `json:"body"`
}

type User struct {
	ID          int    `json:"id"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mu:   &sync.RWMutex{},
	}
	err := db.ensureDB()
	return db, err
}

func (db *DB) CreateChirp(body string, authorID int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	id := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		ID:       id,
		AuthorID: authorID,
		Body:     body,
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

func (db *DB) GetChirpById(id int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := dbStructure.Chirps[id]
	if !ok {
		return Chirp{}, errors.New("chirp not found")
	}

	return chirp, nil
}

func (db *DB) GetChirpsByAuthorID(authorID int) ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0)
	for _, chirp := range dbStructure.Chirps {
		if chirp.AuthorID == authorID {
			chirps = append(chirps, chirp)
		}
	}

	return chirps, nil
}

func (db *DB) GetChirpsByIdAsc() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	sortChirps(chirps, "asc")

	return chirps, nil
}

func (db *DB) GetChirpsByIdDesc() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	sortChirps(chirps, "desc")

	return chirps, nil
}

func sortChirps(chirps []Chirp, order string) {
	if order == "asc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].ID < chirps[j].ID
		})
	} else if order == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].ID > chirps[j].ID
		})
	}
}

func (db *DB) createDB() error {
	dbStructure := DBStructure{
		Chirps:        map[int]Chirp{},
		Users:         map[int]User{},
		RevokedTokens: map[int]string{},
	}
	return db.writeDB(dbStructure)
}

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return db.createDB()
	}
	return err
}

func (db *DB) loadDB() (DBStructure, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	dbStructure := DBStructure{}
	dat, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return dbStructure, err
	}
	err = json.Unmarshal(dat, &dbStructure)
	if err != nil {
		return dbStructure, err
	}

	return dbStructure, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	dat, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, dat, 0600)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) CreateUser(email, password string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, user := range dbStructure.Users {
		if user.Email == email {
			return User{}, errors.New("user already exists")
		}
	}

	id := len(dbStructure.Users) + 1
	user := User{
		ID:          id,
		Email:       email,
		Password:    hashed(password),
		IsChirpyRed: false,
	}
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func hashed(password string) string {
	data, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func compareHash(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (db *DB) Login(email, password string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, user := range dbStructure.Users {
		if user.Email == email && compareHash(user.Password, password) {
			return user, nil
		}
	}

	return User{}, errors.New("user not found")
}

func (db *DB) UpdateUser(id int, email, password string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return errors.New("user not found")
	}

	if email == "" || password == "" {
		return errors.New("email or password cannot be empty")
	}

	user.Email = email
	user.Password = hashed(password)
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) UpdateUserMembership(id int) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return errors.New("user not found")
	}

	user.IsChirpyRed = true
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) RevokeRefreshToken(id int, token string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	_, ok := dbStructure.RevokedTokens[id]
	if ok {
		return errors.New("refresh token already revoked")
	}

	dbStructure.RevokedTokens[id] = token

	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) IsTokenRevoked(id int) (bool, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return false, err
	}

	_, ok := dbStructure.RevokedTokens[id]
	if ok {
		return true, nil
	}

	return false, nil
}

func (db *DB) DeleteChirp(id int) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	_, ok := dbStructure.Chirps[id]
	if !ok {
		return errors.New("chirp not found")
	}

	delete(dbStructure.Chirps, id)

	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}
