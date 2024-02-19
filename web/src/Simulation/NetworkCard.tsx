import React, { useState, useEffect } from 'react';
import {Card, Form, Modal, FormControl, Button} from 'react-bootstrap'
import API from '../API/api';
import NetworkGraph from './NetworkGraph';

type NetworkCardProps = {
    sim:SimInfo;
    steps:Array<Step>;
    readsim(id:string): void;
}

const NetworkCard = (props:NetworkCardProps) => {
    const [idcol, setidcol] = useState<number>(0);
    const [pcol, setpcol] = useState<number>(2);
    const [delim, setdelim] = useState<string>(",");
    const [filetoupload, setfiletoupload] = useState<Blob>();
    const [showimpmodal, setshowimpmodal] = useState<boolean>(false);
    const [hasstep, sethasstep] = useState<boolean>(false);

    useEffect(() => {
        sethasstep((props.steps || []).length > 0)
    },[props.sim, props.steps]);

    const handleimpclose = () => {
        setshowimpmodal(false);
        setidcol(0);
        setpcol(2);
        setdelim(",");
        setfiletoupload({} as Blob);
    };
    
    const handleimport = () => {
        const fr = new FileReader();
        if(!filetoupload) return;

        fr.readAsText(filetoupload);
        fr.onload = function() {
            if(!fr.result) return;
            const base64data = btoa(fr.result.toString());
        
            const pdata:ParseOptions = {
                "identifier": idcol,
                "parent": pcol,
                "regex": {
                    "0":"^\\d+_(.*)$",
                    "2":"^\\d+_(.*)$",
                    "3":"\\S+"
                  },
                "delimiter": delim,
                "payload": base64data
            };
            API.parse(props.sim, pdata).then(response => {
                props.readsim(response.parent);
            });
        };
        handleimpclose();
    };
    
    return(
        <Card className="mb-2 mx-n2">
            <Card.Header><Card.Title>Network<Button size="sm" className="btn btn-primary float-right" onClick={() => setshowimpmodal(true)} disabled={hasstep}>Import</Button></Card.Title></Card.Header>
            <Card.Body>
                <NetworkGraph sim={props.sim} steps={props.steps}/>
            </Card.Body>
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
                        <Form.File label="Select File" onChange={(e:any) => {
                            if(e.target.files.length) setfiletoupload(e.target.files[0]);
                        }}/>
                    </Form.Group>
                    <Form.Group controlId="form-identifier">
                        <Form.Label>Identifier column</Form.Label>
                        <Form.Control type="number" value={idcol} onChange={e => setidcol((e.target as HTMLInputElement).valueAsNumber)}/>
                        <Form.Text className="text-muted">
                            The column that has the unique identifier in it
                        </Form.Text>
                    </Form.Group>
                    <Form.Group controlId="form-parent">
                        <Form.Label>Parent Identifier Column</Form.Label>
                        <Form.Control type="number" value={pcol} onChange={e => setpcol((e.target as HTMLInputElement).valueAsNumber)}/>
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