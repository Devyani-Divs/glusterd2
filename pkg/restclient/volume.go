package restclient

import (
	"fmt"
	"net/http"

	"github.com/gluster/glusterd2/pkg/api"
)

// VolumeCreate creates Gluster Volume
func (c *Client) VolumeCreate(req api.VolCreateReq) (api.Volinfo, error) {
	var vol api.Volinfo
	err := c.post("/v1/volumes", req, http.StatusCreated, &vol)
	return vol, err
}

// Volumes returns list of all volumes
func (c *Client) VolumeAll() ([]api.Volinfo, error) {
	var vols []api.Volinfo
	err := c.get("/v1/volumes", nil, http.StatusOK, &vols)
	return vols, err
}

func (c *Client) Volumes(volname string) (api.Volinfo, error) {
	var vols api.Volinfo
	url := fmt.Sprintf("/v1/volumes/%s", volname)
	err := c.get(url, nil, http.StatusOK, &vols)
	return vols, err
}

// VolumeStart starts a Gluster Volume
func (c *Client) VolumeStart(volname string) error {
	url := fmt.Sprintf("/v1/volumes/%s/start", volname)
	return c.post(url, nil, http.StatusOK, nil)
}

// VolumeStop stops a Gluster Volume
func (c *Client) VolumeStop(volname string) error {
	url := fmt.Sprintf("/v1/volumes/%s/stop", volname)
	return c.post(url, nil, http.StatusOK, nil)
}

// VolumeDelete deletes a Gluster Volume
func (c *Client) VolumeDelete(volname string) error {
	url := fmt.Sprintf("/v1/volumes/%s", volname)
	return c.del(url, nil, http.StatusOK, nil)
}
