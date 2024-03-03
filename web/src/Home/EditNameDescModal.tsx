import React, { useState, useEffect } from 'react';
import {Modal, InputGroup, FormControl, Button} from 'react-bootstrap'
import { SimInfo } from '../API/SimInfo';

type EditNameProps = {
    show: boolean;
    sim: SimInfo;
    closeFunc: () => void;
    saveFunc: (sim:SimInfo) => void;
}

const EditNameDescModal = (props:EditNameProps) => {
    const [sim, setsim] = useState<SimInfo>(props.sim);
    const [simname, setsimname] = useState<string>("");
    const [simdescription, setsimdescription] = useState<string>("");

    const setfields = (ps:SimInfo) => {
        const s = ps;
        setsimname(s['name']);
        setsimdescription(s['description']);
    };

    const handleClose = () => {
        setfields(sim);
        props.closeFunc();
    };

    useEffect(() => {
        setsim(props.sim);
        setfields(props.sim);
      },[props.sim]);

    const saveSimulation = () => {
        var simtosave:SimInfo = sim;
        simtosave.name = simname;
        simtosave.description = simdescription;
        props.saveFunc(simtosave);
        handleClose();
    }

    return(
        <Modal
            show={props.show}
            onHide={handleClose}
            backdrop="static"
            keyboard={false}
        >
            <Modal.Header closeButton>
                <Modal.Title>Simulation</Modal.Title>
            </Modal.Header>
            <Modal.Body>
                <InputGroup className="mb-3">
                    <InputGroup.Prepend>
                        <InputGroup.Text id="basic-on1">Name</InputGroup.Text>
                    </InputGroup.Prepend>
                    <FormControl
                    placeholder="Simulation name"
                    aria-label="Simulation name"
                    aria-describedby="basic-on1"
                    defaultValue={simname}
                    onChange={e => setsimname(e.target.value)}
                    />
                </InputGroup>
                <InputGroup>
                    <InputGroup.Prepend>
                        <InputGroup.Text>Description</InputGroup.Text>
                    </InputGroup.Prepend>
                    <FormControl as="textarea" aria-label="Description" defaultValue={simdescription} onChange={e => setsimdescription(e.target.value)}/>
                </InputGroup>
            </Modal.Body>
            <Modal.Footer>
                <Button variant="primary" onClick={saveSimulation}>Save</Button>
                <Button variant="secondary" onClick={handleClose}>Cancel</Button>
            </Modal.Footer>
        </Modal>
    )
}

export default EditNameDescModal