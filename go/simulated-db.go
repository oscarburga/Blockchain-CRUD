package main

import (
	"encoding/json"
	"fmt"
)

type SimulatedDatabase struct {
	Candidates []CANDIDATE `json:"Candidates"`
	Voters     []VOTANTE   `json:"Votantes"`
}

func (db *SimulatedDatabase) EnsureInitialized() {
	if db.Candidates == nil {
		db.Candidates = make([]CANDIDATE, 0, 10)
	}
	if db.Voters == nil {
		db.Voters = make([]VOTANTE, 0, 10)
	}
}

func (db *SimulatedDatabase) PrintDatabase() {
	fmt.Println("\n--- Database ---")
	jsonData, _ := json.MarshalIndent(db, "", "\t")
	fmt.Println(string(jsonData))
	fmt.Println("")
}

func (db *SimulatedDatabase) DataExistsInCandidates(id int, dni string) bool {
	for i := 0; i < len(db.Candidates); i++ {
		if db.Candidates[i].ID == id || db.Candidates[i].DNI == dni {
			fmt.Println("Data exists in candidate: ", id, dni)
			return true
		}
	}
	return false
}

func (db *SimulatedDatabase) DataExistsInVoters(id int, dni string) bool {
	for i := 0; i < len(db.Voters); i++ {
		if db.Voters[i].ID == id || db.Voters[i].DNI == dni {
			fmt.Println("Data exists in voters: ", id, dni)
			return true
		}
	}
	return false
}

func (db *SimulatedDatabase) CreateCandidate(cand CANDIDATE) bool {
	if cand.ID > 0 && len(cand.DNI) > 0 &&
		!db.DataExistsInCandidates(cand.ID, cand.DNI) {
		db.Candidates = append(db.Candidates, cand)
		return true
	}
	return false
}

// Update con campos vacíos se usa para delete
func (db *SimulatedDatabase) UpdateCandidate(cand CANDIDATE) bool {
	if cand.ID > 0 && db.DataExistsInCandidates(cand.ID, cand.DNI) {
		for i := 0; i < len(db.Candidates); i++ {
			if db.Candidates[i].ID == cand.ID {
				if cand.Nombre == "" && cand.Apellido == "" && cand.DNI == "" && cand.Numero_votacion == "" {
					db.Candidates[i] = db.Candidates[len(db.Candidates)-1]
					db.Candidates[len(db.Candidates)-1] = CANDIDATE{}
					db.Candidates = db.Candidates[:len(db.Candidates)-1]
					fmt.Println("Deleted candidate with data: ", cand.ID)
				} else {
					db.Candidates[i] = cand
					fmt.Println("Updated candidate with data: ", cand.ID)
				}
				return true
			}
		}
	}
	return false
}

func (db *SimulatedDatabase) CreateVotante(voter VOTANTE) bool {
	if voter.ID > 0 && len(voter.DNI) > 0 &&
		!db.DataExistsInVoters(voter.ID, voter.DNI) {
		db.Voters = append(db.Voters, voter)
		return true
	}
	return false
}

// Update con campos vacíos se usa para delete
func (db *SimulatedDatabase) UpdateVotante(voter VOTANTE) bool {
	if voter.ID > 0 && db.DataExistsInVoters(voter.ID, voter.DNI) {
		for i := 0; i < len(db.Voters); i++ {
			if db.Voters[i].ID == voter.ID {
				if voter.Nombre == "" && voter.Apellido == "" && voter.DNI == "" && voter.Candidato_Voto == "" && voter.Lugar_Votacion == "" {
					db.Voters[i] = db.Voters[len(db.Voters)-1]
					db.Voters[len(db.Voters)-1] = VOTANTE{}
					db.Voters = db.Voters[:len(db.Voters)-1]
					fmt.Println("Deleted voter with data: ", voter.ID)
				} else {
					db.Voters[i] = voter
					fmt.Println("Updated voter with data: ", voter.ID)
				}
				return true
			}
		}
	}
	return false
}
