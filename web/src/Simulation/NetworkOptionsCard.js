import React, { useState, useEffect } from 'react';
import API from '../api';
import {Card, Form, Button, OverlayTrigger, Tooltip, Modal} from 'react-bootstrap';

const NetworkOptionsCard = (props) => {
    const [sim, setSim] = useState(props.sim);
    const [showoptmodal, setShowoptmodal] = useState(false);
    
    const [awm, setawm] = useState(false);
    const [el, setel] = useState([]);
    const [ic, setic] = useState([0]);
    const [ltp, setltp] = useState(false);
    const [ltl, setltl] = useState([]);
    const [le, setle] = useState("");
    const [mc, setmc] = useState(2);
 
    useEffect(() => {
        setSim(props.sim);
        setOptions(props.sim.options);
      }, [props.sim]);

    const setOptions = (options) => {
        setawm(options['agentsWithMemory'] === true);
        setel(options['evangelistList'] || []);
        setic(options['initColors'] || []);
        setltp(options['linkTeamPeers'] === true);
        setltl(options['linkedTeamList'] || []);
        setle(options['loneEvangelist']);
        setmc(options['maxColors']);
    };

    const colorFromVal = (color) => {
        switch(color) {
            case 0:
                return "Grey";
            case 1:
                return "Blue";
            case 2:
                return "Red";
            case 3:
                return "Green";
            case 4:
                return "Yellow";
            case 5:
                return "Orange";
            case 6:
                return "Purple";
           default:
              return "Invalid Color";
          }
    };

    const handleoptshow = () => setShowoptmodal(true);

    const handleoptclose = () => {
        setShowoptmodal(false);
        setOptions(sim.options);
    };

    const handlesaveopt = () => {
        setShowoptmodal(false);
        sim.options['agentsWithMemory'] = awm;
        sim.options['evangelistList'] = el;
        sim.options['initColors'] = ic;
        sim.options['linkTeamPeers'] = ltp;
        sim.options['linkedTeamList'] = ltl;
        sim.options['loneEvangelist'] = le;
        sim.options['maxColors'] = mc;
        API.update(sim);
    }   

    return (
        <Card className="mt-3">
            <Card.Header><Card.Title>Network Options<Button size="sm" className="btn btn-primary float-right" onClick={handleoptshow}>Edit</Button></Card.Title></Card.Header>
            <Card.Body className="small">
                <fieldset disabled>
                    <OverlayTrigger overlay={<Tooltip>Use agents that have memory in the network simulation</Tooltip>}>
                        <span>
                            <Form.Check type="checkbox" label="Use agents with memory" checked={awm}/>
                        </span>
                    </OverlayTrigger>
                    <OverlayTrigger overlay={<Tooltip>List the agents that are evangelists for a new idea</Tooltip>}>
                        <span>
                            <Form.Label className="pt-2">Evangelist list</Form.Label>
                            <Form.Control size="sm" as="textarea" value={el.join("; ")}/>
                        </span>
                    </OverlayTrigger>
                    <OverlayTrigger overlay={<Tooltip>Select the colors that the agents will be randomly assigned to</Tooltip>}>
                        <span>
                            <Form.Label className="pt-2">Initial colors</Form.Label>
                            <Form.Control size="sm" type="text" value={ic.map(culr => colorFromVal(culr)).join("; ")}/>
                        </span>
                    </OverlayTrigger>
                    <OverlayTrigger overlay={<Tooltip>Generate links between all members of a team</Tooltip>}>
                        <span>
                            <Form.Check className="pt-2" type="checkbox" label="Link team peers" checked={ltp}/>
                        </span>
                    </OverlayTrigger>
                    <OverlayTrigger overlay={<Tooltip>Generate links between the specified teams</Tooltip>}>
                        <span>
                            <Form.Label className="pt-2">Linked team list</Form.Label>
                            <Form.Control size="sm" as="textarea" value={ltl.join("; ")}/>
                        </span>
                    </OverlayTrigger>
                    <OverlayTrigger overlay={<Tooltip>An agent that will act as an evangelist for a new idea</Tooltip>}>
                        <span>
                            <Form.Label className="pt-2">Lone evangelist</Form.Label>
                            <Form.Control size="sm" type="text" value={le}/>
                        </span>
                    </OverlayTrigger>
                    <OverlayTrigger overlay={<Tooltip>Set the maximum number of colors representing competing ideas in the network</Tooltip>}>
                        <span>
                            <Form.Label className="pt-2">Maximum colors</Form.Label>
                            <Form.Control size="sm" type="text" value={mc}/>
                        </span>
                    </OverlayTrigger>
                </fieldset>
            </Card.Body>
            <Modal
                show={showoptmodal}
                onHide={handleoptclose}
                backdrop="static"
                keyboard={false}
            >
                <Modal.Header closeButton>
                    <Modal.Title>Edit network options</Modal.Title>
                </Modal.Header>
                <Modal.Body>
                    <Form.Group controlId="form-awm">
                        <Form.Check type="checkbox" label="Use agents with memory" checked={awm} onChange={
                            e => setawm(e.target.checked)
                            }/>
                        <Form.Text className="text-muted">
                            Use agents that have memory in the network simulation
                        </Form.Text>
                    </Form.Group>
                    <Form.Group controlId="form-evangelist">
                        <Form.Label>Evangelist list</Form.Label>
                        <Form.Control as="select" value={el} onChange={e => setel(Array.from(e.target.selectedOptions).filter(sel => sel.value))} multiple/>
                        <Form.Text className="text-muted">
                            List the agents that are evangelists for a new idea
                        </Form.Text>
                    </Form.Group>
                    <Form.Group controlId="form-initcolors">
                        <Form.Label>Initialisation colors</Form.Label>
                        <Form.Control as="select" value={ic} onChange={e => setic(Array.from(e.target.selectedOptions).filter(sel => sel.value).map(sel => parseInt(sel.value)))} multiple>
                            <option></option>
                            <option value={0}>{colorFromVal(0)}</option>
                            <option value={1}>{colorFromVal(1)}</option>
                            <option value={2}>{colorFromVal(2)}</option>
                            <option value={3}>{colorFromVal(3)}</option>
                            <option value={4}>{colorFromVal(4)}</option>
                            <option value={5}>{colorFromVal(5)}</option>
                            <option value={6}>{colorFromVal(6)}</option>
                        </Form.Control>
                        <Form.Text className="text-muted">
                            Select the colors that the agents will be randomly assigned to
                        </Form.Text>
                    </Form.Group>
                    <Form.Group controlId="form-ltp">
                        <Form.Check type="checkbox" label="Link team peers" checked={ltp} onChange={e => setltp(e.target.checked)}/>
                        <Form.Text className="text-muted">
                            Generate links between all members of a team
                        </Form.Text>
                    </Form.Group>
                    <Form.Group controlId="form-linkedteam">
                        <Form.Label>Linked team list</Form.Label>
                        <Form.Control as="select" value={ltl} onChange={e => setltl(Array.from(e.target.selectedOptions).filter(sel => sel.value))} multiple/>
                        <Form.Text className="text-muted">
                            Generate links between the specified teams
                        </Form.Text>
                    </Form.Group>
                        <Form.Group controlId="form-loneevangelist">
                        <Form.Label>Lone evangelist</Form.Label>
                        <Form.Control as="select" value={le} onChange={e => setle(e.target.value)}/>
                        <Form.Text className="text-muted">
                            An agent that will act as an evangelist for a new idea
                        </Form.Text>
                    </Form.Group>
                    <Form.Group controlId="form-maxcolors">
                        <Form.Label>Maximum colors</Form.Label>
                        <Form.Control type="number" value={mc} onChange={e => setmc(e.target.valueAsNumber)}/>
                        <Form.Text className="text-muted">
                            Set the maximum number of colors representing competing ideas in the network
                        </Form.Text>
                    </Form.Group>
                </Modal.Body>
                <Modal.Footer>
                    <Button variant="success" onClick={handlesaveopt}>Save</Button>
                    <Button variant="secondary" onClick={handleoptclose}>Cancel</Button>
                </Modal.Footer>
            </Modal>
        </Card>
    )
}

export default NetworkOptionsCard