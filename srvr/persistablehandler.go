package srvr

import (
	"net/http"

	"github.com/spaceweasel/mango"
)

//Updateable must be implemented by any persistable object that can be updated
type Updateable interface {
	Persistable
	CopyValues(objToCopy Persistable) error
}

//PersistableHandlerState holds state information for a PersistableHandlerState
type PersistableHandlerState struct {
	FileManager FileManager
}

//GetObject returns the object
func (ph *PersistableHandlerState) GetObject(obj Persistable, c *mango.Context) {
	objUpdater := ph.FileManager.Get(obj.Filepath())
	err := objUpdater.Read(obj)
	if err != nil {
		c.Error(err.Error(), http.StatusInternalServerError)
	} else {
		c.RespondWith(obj).WithStatus(http.StatusOK)
	}
}

//UpdateObjectWithContextBind updates the saved copy of the passed object
func (ph *PersistableHandlerState) UpdateObjectWithContextBind(obj Persistable, savedObj Updateable, c *mango.Context) {
	err := c.Bind(obj)
	if err != nil {
		c.Error(err.Error(), http.StatusBadRequest)
		return
	}
	err = ph.UpdateObject(obj, savedObj, c)
	if err != nil {
		c.Error(err.Error(), http.StatusInternalServerError)
	} else {
		c.RespondWith(savedObj).WithStatus(http.StatusOK)
	}
}

//UpdateObject updates the saved copy of the passed object
func (ph *PersistableHandlerState) UpdateObject(obj Persistable, savedObj Updateable, c *mango.Context) error {
	objUpdater := ph.FileManager.Get(obj.Filepath())
	var err error

	//Retry if there is a failure
	for i := 0; i < 2; i++ {
		err = objUpdater.Read(savedObj)
		if err != nil {
			continue
		}
		err = savedObj.CopyValues(obj)
		if err != nil {
			break
		}
		err = objUpdater.Update(savedObj)
		if err == nil {
			break
		}
	}
	return err
}
