package srvr

import (
	"net/http"

	"github.com/spaceweasel/mango"
)

//ListHolder is an interface that any persistable that contains a list must implement
type ListHolder interface {
	Persistable
	GetItems(listname string) []string
	UpdateItems(listname string, items []string)
}

//ListItem is an interface that any persistable that can be added to a list must implement
type ListItem interface {
	Persistable
	RelPath() string
}

//ListHandlerState holds state information for a ListHandler
type ListHandlerState struct {
	FileManager FileManager
}

//GetList returns the list of items
func (lh *ListHandlerState) GetList(listHolder ListHolder, c *mango.Context, listname string) {
	listUpdater := lh.FileManager.Get(listHolder.Filepath())
	err := listUpdater.Read(listHolder)
	if err != nil {
		c.RespondWith(err.Error()).WithStatus(http.StatusInternalServerError)
	} else {
		c.RespondWith(listHolder).WithStatus(http.StatusOK)
	}
}

//AddItemWithContextBind binds the itemToAdd to the object passed in the mango context and
//adds it to the specified list on the passed listholder
func (lh *ListHandlerState) AddItemWithContextBind(itemToAdd ListItem, listHolder ListHolder, c *mango.Context, listname string) {
	err := c.Bind(itemToAdd)
	if err != nil {
		c.Error(err.Error(), http.StatusBadRequest)
		return
	}
	err = lh.AddItem(itemToAdd, listHolder, c, listname)
	if err != nil {
		c.RespondWith(err.Error()).WithStatus(http.StatusInternalServerError)
	} else {
		c.RespondWith(itemToAdd).WithStatus(http.StatusCreated)
	}
}

//AddItem adds the passed item to the specified list on the passed listholder
func (lh *ListHandlerState) AddItem(itemToAdd ListItem, listHolder ListHolder, c *mango.Context, listname string) error {
	itemUpdater := lh.FileManager.Get(itemToAdd.Filepath())
	err := itemUpdater.Create(itemToAdd)
	if err != nil {
		return err
	}
	listUpdater := lh.FileManager.Get(listHolder.Filepath())

	//Retry if there is a failure
	for i := 0; i < 2; i++ {
		err = listUpdater.Read(listHolder)
		if err != nil {
			continue
		}
		newlist := append(listHolder.GetItems(listname), itemToAdd.RelPath())
		listHolder.UpdateItems(listname, newlist)
		err = listUpdater.Update(listHolder)
		if err == nil {
			break
		}
	}
	return err
}

//DeleteItem removes an item from the specified list on the passed listholder
func (lh *ListHandlerState) DeleteItem(itemToDelete ListItem, listHolder ListHolder, c *mango.Context, listname string) {
	itemUpdater := lh.FileManager.Get(itemToDelete.Filepath())
	listUpdater := lh.FileManager.Get(listHolder.Filepath())
	relpath := itemToDelete.RelPath()
	var err error
	for i := 0; i < 2; i++ {
		err = listUpdater.Read(listHolder)
		if err != nil {
			continue
		}
		items := listHolder.GetItems(listname)
		count := len(items) - 1
		if count < 0 {
			c.Error("Item not found", http.StatusNotFound)
			return
		}
		upl := make([]string, count)
		j := 0
		for _, item := range items {
			if j == count && item != relpath {
				c.Error("Item not found", http.StatusNotFound)
				return
			}
			if item != relpath {
				upl[j] = item
				j++
			}
		}
		listHolder.UpdateItems(listname, upl)
		err = listUpdater.Update(listHolder)
		if err == nil {
			break
		}
	}
	if err != nil {
		c.RespondWith(err.Error()).WithStatus(http.StatusInternalServerError)
		return
	}
	err = itemUpdater.Delete()
	if err != nil {
		c.RespondWith(err.Error()).WithStatus(http.StatusInternalServerError)
	} else {
		c.Respond().WithStatus(http.StatusOK)
	}
}
