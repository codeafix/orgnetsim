import React, { useState, useEffect } from 'react';
import API from '../API/api';
import SimList from './SimList';
import {Modal} from 'react-bootstrap'
import EditNameDescModal from './EditNameDescModal'
import {Button} from 'react-bootstrap'
import {CardDeck} from 'react-bootstrap';
import Logo from '../logo.svg';

const Home = () => {
    const [simlist, setSimlist] = useState<Array<SimInfo>>([]);
    const [notes, setNotes] = useState<string>("");
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
    const handleDelShow = (id:string, simname:string) => {
        setShowdelmodal(true);
        setDelsimid(id);
        setDelsimname(simname);
    }

    useEffect(() => {
        document.title = "orgnetsim";
        API.sims()
            .then(response => {
                setSimlist(response.simulations);
                setNotes(response.notes)
            })
      },[]);

    const addSimulation = (sim:SimInfo) => {
        API.add(sim.name, sim.description).then(response => {
            setSimlist(simlist.concat(response));
        });
    };
    
    const deleteSimulation = (id:string) => {
        API.delete(id).then(() => {
            setSimlist(simlist.filter(item => item.id !== id));
        });
        handleDelClose();
    };

    return(
        <div className="container-fluid">
            <h1><img src={Logo} style={{ height: 60, width: 60}} alt="logo"/>Simulation Set</h1>
            <p>{notes}</p>
            <h2>List of Simulations<Button className="btn btn-primary float-right" onClick={handleAddShow}>Add</Button></h2>
            <CardDeck>
                <SimList sims={simlist} deleteFunc={handleDelShow}/>
            </CardDeck>
            <EditNameDescModal sim={API.emptySim()} show={showaddmodal} saveFunc={addSimulation} closeFunc={handleAddClose}/>
            <Modal
                show={showdelmodal}
                onHide={handleDelClose}
                backdrop="static"
                keyboard={false}
            >
                <Modal.Header closeButton={true}>
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