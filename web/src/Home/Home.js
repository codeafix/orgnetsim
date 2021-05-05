import React, { useState, useEffect } from 'react';
import API from '../api';
import SimList from './SimList';
import {Modal} from 'react-bootstrap'
import EditNameDescModal from './EditNameDescModal'
import {Button} from 'react-bootstrap'
import {CardDeck} from 'react-bootstrap';

const Home = () => {
    const [simlist, setSimlist] = useState([]);
    const [notes, setNotes] = useState("");
    const [showaddmodal, setShowaddmodal] = useState(false);
    const [showdelmodal, setShowdelmodal] = useState(false);
    const [delsimid, setDelsimid] = useState("");
    const [delsimname, setDelsimname] = useState("");

    const handleAddClose = () => setShowaddmodal(false);
    const handleAddShow = () => setShowaddmodal(true);

    const handleDelClose = () => {
        setShowdelmodal(false);
        setDelsimid("");
        setDelsimname("");
    }
    const handleDelShow = (id, simname) => {
        setShowdelmodal(true);
        setDelsimid(id);
        setDelsimname(simname);
    }

    useEffect(() => {
        API.sims()
            .then(response => {
                setSimlist(response.simulations);
                setNotes(response.notes)
            })
      },[]);

    const addSimulation = (sim) => {
        API.add(sim.name, sim.description).then(response => {
            setSimlist(simlist.concat(response));
        });
    };
    
    const deleteSimulation = (id) => {
        API.delete(id).then(() => {
            setSimlist(simlist.filter(item => item.id !== id));
        });
        handleDelClose();
    };

    return(
        <div className="container-fluid">
            <h1>Simulation Set</h1>
            <p>{notes}</p>
            <h2>List of Simulations<Button className="btn btn-primary float-right" onClick={handleAddShow}>Add</Button></h2>
            <CardDeck>
                <SimList sims={simlist} deleteFunc={handleDelShow}/>
            </CardDeck>
            <EditNameDescModal show={showaddmodal} saveFunc={addSimulation} closeFunc={handleAddClose}/>
            <Modal
                show={showdelmodal}
                onHide={handleDelClose}
                backdrop="static"
                keyboard={false}
            >
                <Modal.Header closeButton>
                    <Modal.Title className="text-danger">Delete "{delsimname}"</Modal.Title>
                </Modal.Header>
                <Modal.Body>
                    Are you sure you want to permanently delete this simulation?
                </Modal.Body>
                <Modal.Footer>
                    <Button variant="warning" onClick={() => deleteSimulation(delsimid)}>Delete</Button>
                    <Button variant="success" onClick={handleDelClose}>Cancel</Button>
                </Modal.Footer>
            </Modal>
        </div>
    )
}

export default Home