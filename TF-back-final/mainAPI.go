package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const TransCmdCreateCandidate = "CREATE CANDIDATE"
const TransCmdCreateVotante = "CREATE VOTANTE"
const TransCmdUpdateCandidate = "UPDATE CANDIDATE"
const TransCmdUpdateVotante = "UPDATE VOTANTE"
const TransCmdGetCandidates = "GET CANDIDATES"
const TransCmdGetVotantes = "GET VOTANTES"

var candidatos []CANDIDATE
var votantes []VOTANTE
var maxCandId int = -1
var maxVotanteId int = -1
var chRemotes chan []string
var host = "localhost:8080"
var remote = "localhost:8000"

type Frame struct {
	Cmd    string   `json:"cmd"`
	Sender string   `json:"sender"`
	Data   []string `json:"data"`
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func updateMaxCandId() {
	for _, cand := range candidatos {
		if cand.ID > maxCandId {
			maxCandId = cand.ID
		}
	}
}

func updateMaxVoterId() {
	for _, voter := range votantes {
		if voter.ID > maxVotanteId {
			maxVotanteId = voter.ID
		}
	}
}

//Para el candidato
func GetCandidatos(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if !send(remote, Frame{"server", host, []string{TransCmdGetCandidates}}, func(cn net.Conn) {
		log.Printf("Get Candidates")
		dec := json.NewDecoder(cn)
		var cands []CANDIDATE
		dec.Decode(&cands)
		candidatos = cands
		updateMaxCandId()
		json.NewEncoder(w).Encode(cands)
		log.Printf("%s: candidatos: %s\n", host, cands)
	}) {
		log.Printf("%s: unable to connect to %s\n", host, remote)
	}
}

//CREATE
func CreateCandidato(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	var NewCandidato CANDIDATE
	reqbody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Invalid Data")
	}
	json.Unmarshal(reqbody, &NewCandidato)
	NewCandidato.ID = maxCandId + 1
	jsonData, err := json.MarshalIndent(NewCandidato, "", "\t")
	if !send(remote, Frame{"server", host, []string{TransCmdCreateCandidate, string(jsonData)}}, func(cn net.Conn) {
		log.Printf("Create Candidate")
		dec := json.NewDecoder(cn)
		var cands []CANDIDATE
		dec.Decode(&cands)
		candidatos = cands
		updateMaxCandId()
		json.NewEncoder(w).Encode(cands)
		log.Printf("%s: candidatos: %s\n", host, cands)
	}) {
		log.Printf("%s: unable to connect to %s\n", host, remote)
	}
}

//READ
func ReadCandidato(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	params := mux.Vars(r)
	CandidatoID, err := strconv.Atoi(params["id"])
	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
		return
	}

	for _, _candidato := range candidatos {
		if _candidato.ID == CandidatoID {
			json.NewEncoder(w).Encode(_candidato)
		}
	}
}

//UPDATE
func UpdateCandidato(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	//params := mux.Vars(r)
	//CandidatoID, err := strconv.Atoi(params["id"]) // asumo que el id se encuentra en el body
	var updateCandidato CANDIDATE

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Invalid Data")

	}
	json.Unmarshal(reqBody, &updateCandidato)
	jsonData, err := json.MarshalIndent(updateCandidato, "", "\t")
	if !send(remote, Frame{"server", host, []string{TransCmdUpdateCandidate, string(jsonData)}}, func(cn net.Conn) {
		log.Printf("Update Candidate")
		dec := json.NewDecoder(cn)
		var cands []CANDIDATE
		dec.Decode(&cands)
		candidatos = cands
		updateMaxCandId()
		json.NewEncoder(w).Encode(cands)
		log.Printf("%s: candidatos: %s\n", host, cands)
	}) {
		log.Printf("%s: unable to connect to %s\n", host, remote)
	}

}

//DELETED
func DeletedCandidato(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	params := mux.Vars(r)
	CandidatoID, err := strconv.Atoi(params["id"])
	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
		return
	}
	deleteCandidate := CANDIDATE{}
	deleteCandidate.ID = CandidatoID
	jsonData, err := json.MarshalIndent(deleteCandidate, "", "\t")
	if !send(remote, Frame{"server", host, []string{TransCmdUpdateCandidate, string(jsonData)}}, func(cn net.Conn) {
		log.Printf("Update Candidate")
		dec := json.NewDecoder(cn)
		var cands []CANDIDATE
		dec.Decode(&cands)
		candidatos = cands
		updateMaxCandId()
		json.NewEncoder(w).Encode(cands)
		log.Printf("%s: candidatos: %s\n", host, cands)
	}) {
		log.Printf("%s: unable to connect to %s\n", host, remote)
	}

}

//Para el votante
func GetVotantes(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if !send(remote, Frame{"server", host, []string{TransCmdGetVotantes}}, func(cn net.Conn) {
		log.Printf("Get Votantes")
		dec := json.NewDecoder(cn)
		var voters []VOTANTE
		dec.Decode(&voters)
		votantes = voters
		updateMaxVoterId()
		json.NewEncoder(w).Encode(voters)
		log.Printf("%s: votantes: %s\n", host, voters)
	}) {
		log.Printf("%s: unable to connect to %s\n", host, remote)
	}
	log.Printf("HOLAAAAAAA")
}

//CREATE
func CreateVotante(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	var NewVotante VOTANTE
	reqbody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Invalid Data")
	}

	json.Unmarshal(reqbody, &NewVotante)

	NewVotante.ID = maxVotanteId + 1
	jsonData, err := json.MarshalIndent(NewVotante, "", "\t")
	if !send(remote, Frame{"server", host, []string{TransCmdCreateVotante, string(jsonData)}}, func(cn net.Conn) {
		log.Printf("Create Votante")
		dec := json.NewDecoder(cn)
		var voters []VOTANTE
		dec.Decode(&voters)
		votantes = voters
		updateMaxVoterId()
		json.NewEncoder(w).Encode(voters)
		log.Printf("%s: votantes: %s\n", host, voters)
	}) {
		log.Printf("%s: unable to connect to %s\n", host, remote)
	}

}

//READ
func ReadVotante(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	params := mux.Vars(r)
	votanteID, err := strconv.Atoi(params["id"])
	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
		return
	}

	for _, _votante := range votantes {
		if _votante.ID == votanteID {
			json.NewEncoder(w).Encode(_votante)
		}
	}
}

//UPDATE
func UpdateVotante(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	//params := mux.Vars(r)
	//votanteID, err := strconv.Atoi(params["id"]) asumo que el id se encuentra el body
	var updateVotante VOTANTE

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Invalid Data")

	}
	json.Unmarshal(reqBody, &updateVotante)
	jsonData, err := json.MarshalIndent(updateVotante, "", "\t")
	if !send(remote, Frame{"server", host, []string{TransCmdUpdateVotante, string(jsonData)}}, func(cn net.Conn) {
		log.Printf("Update Votante")
		dec := json.NewDecoder(cn)
		var voters []VOTANTE
		dec.Decode(&voters)
		votantes = voters
		updateMaxVoterId()
		json.NewEncoder(w).Encode(voters)
		log.Printf("%s: votantes: %s\n", host, voters)
	}) {
		log.Printf("%s: unable to connect to %s\n", host, remote)
	}

}

//DELETED
func DeletedVotante(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	params := mux.Vars(r)
	votanteID, err := strconv.Atoi(params["id"])
	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
		return
	}
	deleteVotante := VOTANTE{}
	deleteVotante.ID = votanteID
	jsonData, err := json.MarshalIndent(deleteVotante, "", "\t")
	if !send(remote, Frame{"server", host, []string{TransCmdUpdateVotante, string(jsonData)}}, func(cn net.Conn) {
		log.Printf("Delete Votante")
		dec := json.NewDecoder(cn)
		var voters []VOTANTE
		dec.Decode(&voters)
		votantes = voters
		updateMaxVoterId()
		json.NewEncoder(w).Encode(voters)
		log.Printf("%s: votantes: %s\n", host, voters)
	}) {
		log.Printf("%s: unable to connect to %s\n", host, remote)
	}

}

type CANDIDATE struct {
	ID              int    `json:"id,omitempty"`
	Nombre          string `json:"nombre,omitempty"`
	Apellido        string `json:"apellido,omitempty"`
	DNI             string `json:"dni,omitempty"`
	Numero_votacion string `json:"numero_votacion,omitempty"`
}

//CREA votante
//READ botante
//UP datos
//delete
type VOTANTE struct {
	ID             int    `json:"id,omitempty"`
	DNI            string `json:"dni,omitempty"`
	Nombre         string `json:"nombre,omitempty"`
	Apellido       string `json:"apellido,omitempty"`
	Lugar_Votacion string `json:"lugar_votacion,omitempty"`
	Candidato_Voto string `json:"candidato_voto,omitempty"`
}

func main_api() {

	log.Printf("HOLAAAAAAA")
	chRemotes = make(chan []string, 1)
	chRemotes <- []string{}
	router := mux.NewRouter()

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})

	router.HandleFunc("/candidatos", GetCandidatos).Methods(http.MethodGet)
	router.HandleFunc("/candidatos/create", CreateCandidato).Methods(http.MethodPost)
	router.HandleFunc("/candidatos/find/{id}", ReadCandidato).Methods(http.MethodGet)
	router.HandleFunc("/candidatos/update/{id}", UpdateCandidato).Methods(http.MethodPut)
	router.HandleFunc("/candidatos/delete/{id}", DeletedCandidato).Methods(http.MethodDelete)

	router.HandleFunc("/votantes", GetVotantes).Methods("GET")
	router.HandleFunc("/votantes/create", CreateVotante).Methods("POST")
	router.HandleFunc("/votantes/find/{id}", ReadVotante).Methods("GET")
	router.HandleFunc("/votantes/update/{id}", UpdateVotante).Methods("PUT")
	router.HandleFunc("/votantes/delete/{id}", DeletedVotante).Methods("DELETE")
	getRemotes("localhost:8000")

	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(headers, methods, origins)(router)))

}

func getRemotes(remote string) {
	if !send(remote, Frame{"server", host, []string{}}, func(cn net.Conn) {
		log.Printf("a")
		dec := json.NewDecoder(cn)
		var frame Frame
		log.Printf("a2")
		dec.Decode(&frame)
		remotes := <-chRemotes
		log.Printf("a3")
		remotes = append(remotes, frame.Data...)
		chRemotes <- remotes
		log.Printf("%s: friends: %s\n", host, remotes)
	}) {
		log.Printf("%s: unable to connect to %s\n", host, remote)
	}
	log.Printf("HOLAAAAAAA")
}

func send(remote string, frame Frame, callback func(net.Conn)) bool {
	if cn, err := net.Dial("tcp", remote); err == nil {
		defer cn.Close()
		enc := json.NewEncoder(cn)
		enc.Encode(frame)
		if callback != nil {
			callback(cn)
		}
		return true
	} else {
		log.Printf("%s: can't connect to %s\n", host, remote)
		idx := -1
		remotes := <-chRemotes
		for i, rem := range remotes {
			if remote == rem {
				idx = i
				break
			}
		}
		if idx >= 0 {
			remotes[idx] = remotes[len(remotes)-1]
			remotes = remotes[:len(remotes)-1]
		}
		chRemotes <- remotes
		return false
	}
}
