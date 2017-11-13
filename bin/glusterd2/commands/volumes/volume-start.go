package volumecommands

import (
	"net/http"

	"github.com/gluster/glusterd2/bin/glusterd2/gdctx"
	restutils "github.com/gluster/glusterd2/bin/glusterd2/servers/rest/utils"
	"github.com/gluster/glusterd2/bin/glusterd2/transaction"
	"github.com/gluster/glusterd2/pkg/errors"
	"github.com/gluster/glusterd2/volume"

	"github.com/gorilla/mux"
	"github.com/pborman/uuid"
	log "github.com/sirupsen/logrus"
)

func startAllBricks(c transaction.TxnCtx) error {
	var volname string
	if err := c.Get("volname", &volname); err != nil {
		return err
	}

	volinfo, err := volume.GetVolume(volname)
	if err != nil {
		return err
	}

	for _, b := range volinfo.Bricks {

		if !uuid.Equal(b.NodeID, gdctx.MyUUID) {
			continue
		}

		c.Logger().WithFields(log.Fields{
			"volume": b.VolumeName,
			"brick":  b.String(),
		}).Info("Starting brick")

		if err := startBrick(b); err != nil {
			return err
		}
	}

	return nil
}

func stopAllBricks(c transaction.TxnCtx) error {
	var volname string
	if e := c.Get("volname", &volname); e != nil {
		c.Logger().WithFields(log.Fields{
			"error": e,
			"key":   "volname",
		}).Error("failed to get value for key from context")
		return e
	}

	vol, e := volume.GetVolume(volname)
	if e != nil {
		// this shouldn't happen
		c.Logger().WithFields(log.Fields{
			"error":   e,
			"volname": volname,
		}).Error("failed to get volinfo for volume")
		return e
	}

	for _, b := range vol.Bricks {

		if !uuid.Equal(b.NodeID, gdctx.MyUUID) {
			continue
		}

		c.Logger().WithFields(log.Fields{
			"volume": b.VolumeName,
			"brick":  b.String(),
		}).Info("volume start failed, stopping brick")

		if err := stopBrick(b); err != nil {
			return err
		}
	}

	return nil
}

func registerVolStartStepFuncs() {
	transaction.RegisterStepFunc(startAllBricks, "vol-start.Commit")
	transaction.RegisterStepFunc(stopAllBricks, "vol-start.Undo")
}

func volumeStartHandler(w http.ResponseWriter, r *http.Request) {
	p := mux.Vars(r)
	volname := p["volname"]
	reqID, logger := restutils.GetReqIDandLogger(r)

	vol, e := volume.GetVolume(volname)
	if e != nil {
		restutils.SendHTTPError(w, http.StatusNotFound, errors.ErrVolNotFound.Error())
		return
	}
	if vol.Status == volume.VolStarted {
		restutils.SendHTTPError(w, http.StatusBadRequest, errors.ErrVolAlreadyStarted.Error())
		return
	}

	// A simple one-step transaction to start the brick processes
	txn := transaction.NewTxn(reqID)
	defer txn.Cleanup()
	lock, unlock, err := transaction.CreateLockSteps(volname)
	if err != nil {
		restutils.SendHTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}
	txn.Nodes = vol.Nodes()
	txn.Steps = []*transaction.Step{
		lock,
		{
			DoFunc:   "vol-start.Commit",
			UndoFunc: "vol-start.Undo",
			Nodes:    txn.Nodes,
		},
		unlock,
	}
	txn.Ctx.Set("volname", volname)

	_, e = txn.Do()
	if e != nil {
		logger.WithFields(log.Fields{
			"error":  e.Error(),
			"volume": volname,
		}).Error("failed to start volume")
		restutils.SendHTTPError(w, http.StatusInternalServerError, e.Error())
		return
	}

	vol.Status = volume.VolStarted

	e = volume.AddOrUpdateVolumeFunc(vol)
	if e != nil {
		restutils.SendHTTPError(w, http.StatusInternalServerError, e.Error())
		return
	}
	restutils.SendHTTPResponse(w, http.StatusOK, vol)
}
