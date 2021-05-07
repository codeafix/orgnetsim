import React, { useState, useEffect } from 'react';
import {Card, Form, Modal, FormControl, Button} from 'react-bootstrap'
import API from '../api';

const NetworkCard = (props) => {
    const [idcol, setidcol] = useState(0);
    const [pcol, setpcol] = useState(2);
    const [delim, setdelim] = useState(",");
    const [filetoupload, setfiletoupload] = useState();
    const [showimpmodal, setshowimpmodal] = useState(false);
    const [hasstep, sethasstep] = useState(false);

    useEffect(() => {
        sethasstep((props.sim.steps || []).length > 0)
    },[props.sim]);

    const handleimpclose = () => {
        setshowimpmodal(false);
        setidcol(0);
        setpcol(2);
        setdelim(",");
        setfiletoupload();
    };
    
    const handleimport = () => {
        const fr = new FileReader();

        fr.readAsText(filetoupload);
        fr.onload = function() {
            const buff = new Buffer(fr.result);
            const base64data = buff.toString('base64');
        
            const pdata = {
                "identifier": idcol,
                "parent": pcol,
                "regex": {
                    "0":"^\\d+_(.*)$",
                    "2":"^\\d+_(.*)$",
                    "3":"\\S+"
                  },
                "delimiter": delim,
                "Payload": base64data
            };
            API.parse(props.sim, pdata).then(response => {
                API.get(response.parent).then(sim => {
                    props.setsim(sim);
                })
            })
        };
        handleimpclose();
    };
    
    return(
        <Card className="mb-2 mx-n2">
            <Card.Header><Card.Title>Network<Button size="sm" className="btn btn-primary float-right" onClick={() => setshowimpmodal(true)} disabled={hasstep}>Import</Button></Card.Title></Card.Header>
            <Modal
                show={showimpmodal}
                onHide={handleimpclose}
                backdrop="static"
                keyboard={false}
            >
                <Modal.Header closeButton>
                    <Modal.Title>Import Network</Modal.Title>
                </Modal.Header>
                <Modal.Body>
                    <Form.Group controlId="form-file">
                        <Form.File label="Select File" onChange={(e) => {
                            if(e.target.files.length) setfiletoupload(e.target.files[0]);
                        }}/>
                    </Form.Group>
                    <Form.Group controlId="form-identifier">
                        <Form.Label>Identifier column</Form.Label>
                        <Form.Control type="number" value={idcol} onChange={e => setidcol(e.target.valueAsNumber)}/>
                        <Form.Text className="text-muted">
                            The column that has the unique identifier in it
                        </Form.Text>
                    </Form.Group>
                    <Form.Group controlId="form-parent">
                        <Form.Label>Parent Identifier Column</Form.Label>
                        <Form.Control type="number" value={pcol} onChange={e => setpcol(e.target.valueAsNumber)}/>
                        <Form.Text className="text-muted">
                            The column that has a parent identifier in it to show hierarchy in the network
                        </Form.Text>
                    </Form.Group>
                    <Form.Group className="mb-3" controlId="form-delimiter">
                        <Form.Label>Delimiter</Form.Label>
                        <FormControl type="string" value={delim} onChange={e => setdelim(e.target.value)}/>
                        <Form.Text className="text-muted">
                            The delimiter used to separate columns in rows of the imported file
                        </Form.Text>
                    </Form.Group>
                </Modal.Body>
                <Modal.Footer>
                    <Button variant="success" onClick={handleimport}>Save</Button>
                    <Button variant="secondary" onClick={handleimpclose}>Cancel</Button>
                </Modal.Footer>
            </Modal>
        </Card>
    )
}

export default NetworkCard