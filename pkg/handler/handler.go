package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	ts "google.golang.org/protobuf/types/known/timestamppb"

	"github.com/amodkala/raft/pkg/common"
	"github.com/amodkala/raft/pkg/raft"
	tx "github.com/amodkala/raft/pkg/transaction"
)

const (
	contentTypeJSON = "application/json"
)

type ResponseBody struct {
	TxID string `json:"tx_id"`
}

func NewWriteHandler(cm *raft.CM) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		var entries []*common.Entry

		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()

		if err := dec.Decode(&entries); err != nil {
			http.Error(w, fmt.Sprintf("error decoding request body: %v", err), http.StatusBadRequest)
			return
		}

		now := ts.Now()
		for _, entry := range entries {
			entry.CreateTime = now
		}

		writeTx := tx.New(entries...)

		if leader, err := cm.Replicate(writeTx); err != nil {
			switch leader {
			case "":
				http.Error(w, fmt.Sprintf("error replicating entries: %v", err), http.StatusInternalServerError)
				return
			default:
				leader := cm.Leader()
				leaderResponse, err := http.Post(leader+"/write", contentTypeJSON, r.Body)
				if err != nil || leaderResponse.StatusCode != http.StatusOK {
					http.Error(w, fmt.Sprintf("error forwarding request to cluster leader: %v", err), http.StatusInternalServerError)
					return
				}
				defer leaderResponse.Body.Close()
				resBytes, err := io.ReadAll(leaderResponse.Body)
				if err != nil {
					http.Error(w, fmt.Sprintf("error reading response from leader: %v", err), http.StatusInternalServerError)
					return
				}
				w.Write(resBytes)
				return
			}
		}

		resBytes, err := json.Marshal(ResponseBody{
			TxID: writeTx.ID(),
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("error marshalling response bytes: %v", err), http.StatusInternalServerError)
			return
		}

		w.Write(resBytes)
	}
}
