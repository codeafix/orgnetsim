import React, { useState, useEffect } from 'react';
import API from '../api';
import SimList from './SimList';
import Modal from 'react-bootstrap/Modal'

const Home = () => {
    const [simlist, setSimlist] = useState([]);
    const [notes, setNotes] = useState("");
    const [showmodal, setShowmodal] = useState(false);
    const handleClose = () => setShowmodal(false);
    const handleShow = () => setShowmodal(true);  

    useEffect(() => {
        API.sims()
            .then(response => {
                setSimlist(response.simulations);
                setNotes(response.notes)
                API.simCount = response.simulations.length;
            })
      },[]);

    const addSimulation = () => {
        API.add().then(response => {
            setSimlist(simlist.concat(response));
        });
        handleClose();
    }
    
    return(
        <div>
            <h1>Simulation Set</h1>
            <p>{notes}</p>
            <h2>List of Simulations</h2>
            <SimList sims={simlist}/>
            <button onClick={handleShow}>Add</button>
            <Modal
                show={showmodal}
                onHide={handleClose}
                backdrop="static"
                keyboard={false}
            >
                <Modal.Header closeButton>
                    <Modal.Title>Add Simulation</Modal.Title>
                </Modal.Header>
                <Modal.Body>
                    Add a new simulation
                </Modal.Body>
                <Modal.Footer>
                    <button variant="secondary" onClick={handleClose}>Close</button>
                    <button variant="primary" onClick={addSimulation}>Add</button>
                </Modal.Footer>
            </Modal>
        </div>
    )
}

export default Home