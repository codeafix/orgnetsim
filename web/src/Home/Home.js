import React, { useState, useEffect } from 'react';
import API from '../api';
import SimList from './SimList';
import Modal from 'react-bootstrap/Modal'
import InputGroup from 'react-bootstrap/InputGroup'
import FormControl from 'react-bootstrap/FormControl'
import Button from 'react-bootstrap/Button'
import CardDeck from 'react-bootstrap/CardDeck';

const Home = () => {
    const [simlist, setSimlist] = useState([]);
    const [notes, setNotes] = useState("");
    const [showmodal, setShowmodal] = useState(false);
    const [addsimname, setAddsimname] = useState("");
    const [addsimdescription, setAddsimdescription] = useState("");

    const handleClose = () => {
        setShowmodal(false);
        setAddsimname("");
        setAddsimdescription("");
    };
    const handleShow = () => setShowmodal(true);  

    useEffect(() => {
        API.sims()
            .then(response => {
                setSimlist(response.simulations);
                setNotes(response.notes)
            })
      },[]);

    const addSimulation = () => {
        API.add(addsimname, addsimdescription).then(response => {
            setSimlist(simlist.concat(response));
        });
        handleClose();
    }
    
    const deleteSimulation = (id) => {
        API.delete(id).then(() => {
            setSimlist(simlist.filter(item => item.id !== id));
        });
    }

    return(
        <div>
            <h1>Simulation Set</h1>
            <p>{notes}</p>
            <h2>List of Simulations<button class="btn btn-primary float-right" onClick={handleShow}>Add</button></h2>
            <CardDeck>
                <SimList sims={simlist} deleteFunc={deleteSimulation}/>
            </CardDeck>
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
                    <InputGroup className="mb-3">
                        <InputGroup.Prepend>
                            <InputGroup.Text id="basic-addon1">Name</InputGroup.Text>
                        </InputGroup.Prepend>
                        <FormControl
                        placeholder="Simulation name"
                        aria-label="Simulation name"
                        aria-describedby="basic-addon1"
                        onChange={e => setAddsimname(e.target.value)}
                        />
                    </InputGroup>
                    <InputGroup>
                        <InputGroup.Prepend>
                            <InputGroup.Text>Description</InputGroup.Text>
                        </InputGroup.Prepend>
                        <FormControl as="textarea" aria-label="Description" onChange={e => setAddsimdescription(e.target.value)}/>
                    </InputGroup>
                </Modal.Body>
                <Modal.Footer>
                    <Button variant="primary" onClick={addSimulation}>Add</Button>
                    <Button variant="secondary" onClick={handleClose}>Cancel</Button>
                </Modal.Footer>
            </Modal>
        </div>
    )
}

export default Home