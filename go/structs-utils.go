package main

import (
	"crypto/sha256"
	"fmt"
)

func GetHashOfString(strToHash string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(strToHash)))
}

func GetHashOfBytes(bytesToHash []byte) string {
	return fmt.Sprintf("%x", sha256.Sum256(bytesToHash))
}

type CANDIDATE struct {
	ID              int    `json:"id,omitempty"`
	Nombre          string `json:"nombre,omitempty"`
	Apellido        string `json:"apellido,omitempty"`
	DNI             string `json:"dni,omitempty"`
	Numero_votacion string `json:"numero_votacion,omitempty"`
}

type VOTANTE struct {
	ID             int    `json:"id,omitempty"`
	DNI            string `json:"dni,omitempty"`
	Nombre         string `json:"nombre,omitempty"`
	Apellido       string `json:"apellido,omitempty"`
	Lugar_Votacion string `json:"lugar_votacion,omitempty"`
	Candidato_Voto string `json:"candidato_voto,omitempty"`
}

type Transaction struct {
	Cmd           string    `json:"Cmd,omitempty"`
	Sender        string    `json:"Sender,omitempty"`
	CandidateData CANDIDATE `json:"CandidateData,omitempty"`
	VotanteData   VOTANTE   `json:"VotanteData,omitempty"`
}

const TransCmdGetCandidates = "GET CANDIDATES"
const TransCmdGetVotantes = "GET VOTANTES"
const TransCmdCreateCandidate = "CREATE CANDIDATE"
const TransCmdCreateVotante = "CREATE VOTANTE"
const TransCmdUpdateCandidate = "UPDATE CANDIDATE"
const TransCmdUpdateVotante = "UPDATE VOTANTE"
