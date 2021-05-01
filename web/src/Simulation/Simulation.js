import React, { useState, useEffect } from 'react';
import API from '../api';
import {Link} from 'react-router-dom';
import {ArrowLeftCircle} from 'react-bootstrap-icons';
import {Card, Form} from 'react-bootstrap'
import {Button} from 'react-bootstrap'
import {OverlayTrigger} from 'react-bootstrap'
import {Tooltip} from 'react-bootstrap'
import {Row} from 'react-bootstrap'
import {Col} from 'react-bootstrap'
import {Container} from 'react-bootstrap'
import {Modal} from 'react-bootstrap'

const Simulation = (props) => {
    const [sim, setSim] = useState({options:{}});
    const [showoptmodal, setShowoptmodal] = useState(false);

    const [awm, setawm] = useState(false);
    const [el, setel] = useState([]);
    const [ic, setic] = useState([0]);
    const [ltp, setltp] = useState(false);
    const [ltl, setltl] = useState([]);
    const [le, setle] = useState("");
    const [mc, setmc] = useState(2);
    
    const setOptions = (options) => {
        setawm(options['agentsWithMemory'] === true);
        setel(options['evangelistList'] || []);
        setic(options['initColors'] || []);
        setltp(options['linkTeamPeers'] === true);
        setltl(options['linkedTeamList'] || []);
        setle(options['loneEvangelist']);
        setmc(options['maxColors']);
    }

    useEffect(() => {
        API.get(props.match.params.id).then(response => {
            setSim(response);
            setOptions(response.options);
        })
      }, [props.match.params.id]);

    const handleoptshow = () => setShowoptmodal(true);

    const handleoptclose = () => {
        setShowoptmodal(false);
        setOptions(sim.options);
    }

    return (
        <Container>
            <h1 class="px-2 bg-light rounded"><Link className="btn btn-outline-secondary mr-3 mt-n2" to='/' role="button"><ArrowLeftCircle className="mt-n1" /></Link>{sim.name}</h1>
            <Container>
                <Row>
                    <Col>
                        <Card>
                            <Card.Header><Card.Title>Description</Card.Title></Card.Header>
                            <Card.Body><Card.Text>{sim.description}</Card.Text></Card.Body>
                        </Card>
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
                                    <OverlayTrigger overlay={<Tooltip>Select the colours that the agents will be randomly assigned to</Tooltip>}>
                                        <span>
                                            <Form.Label className="pt-2">Initial colours</Form.Label>
                                            <Form.Control size="sm" type="text" value={ic.join("; ")}/>
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
                                    <OverlayTrigger overlay={<Tooltip>Set the maximum number of colours representing competing ideas in the network</Tooltip>}>
                                        <span>
                                            <Form.Label className="pt-2">Maximum colours</Form.Label>
                                            <Form.Control size="sm" type="text" value={mc}/>
                                        </span>
                                    </OverlayTrigger>
                                </fieldset>
                            </Card.Body>
                        </Card>
                    </Col>
                    <Col>
                        <Card>
                            <Card.Header><Card.Title>Network</Card.Title></Card.Header>
                        </Card>
                        <Card className="mt-3">
                            <Card.Header><Card.Title>Steps</Card.Title></Card.Header>
                        </Card>
                    </Col>
                </Row>
            </Container>
            
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
                        <Form.Check type="checkbox" label="Use agents with memory" checked={awm} onChange={e => setawm(e.target.value)}/>
                        <Form.Text className="text-muted">
                            Use agents that have memory in the network simulation
                        </Form.Text>
                    </Form.Group>
                    <Form.Group controlId="form-evangelist">
                        <Form.Label>Evangelist list</Form.Label>
                        <Form.Control as="select" value={el} onChange={e => setel(e.target.value)} multiple/>
                        <Form.Text className="text-muted">
                            List the agents that are evangelists for a new idea
                        </Form.Text>
                    </Form.Group>
                    <Form.Group controlId="form-initcolours">
                        <Form.Label>Initialisation colours</Form.Label>
                        <Form.Control as="select" value={ic} onChange={e => setic(Array.from(e.target.selectedOptions).filter(sel => sel.value).map(sel => parseInt(sel.value)))} multiple>
                            <option></option>
                            <option value={0}>Grey</option>
                            <option value={1}>Blue</option>
                            <option value={2}>Red</option>
                            <option value={3}>Green</option>
                            <option value={4}>Yellow</option>
                            <option value={5}>Orange</option>
                            <option value={6}>Purple</option>
                        </Form.Control>
                        <Form.Text className="text-muted">
                            Select the colours that the agents will be randomly assigned to
                        </Form.Text>
                    </Form.Group>
                    <Form.Group controlId="form-ltp">
                        <Form.Check type="checkbox" label="Link team peers" checked={ltp} onChange={e => setltp(e.target.value)}/>
                        <Form.Text className="text-muted">
                            Generate links between all members of a team
                        </Form.Text>
                    </Form.Group>
                    <Form.Group controlId="form-linkedteam">
                        <Form.Label>Linked team list</Form.Label>
                        <Form.Control as="select" value={ltl} onChange={e => setltl(e.target.value)} multiple/>
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
                    <Form.Group controlId="form-maxcolours">
                        <Form.Label>Maximum colours</Form.Label>
                        <Form.Control type="number" value={mc} onChange={e => setmc(e.target.value)}/>
                        <Form.Text className="text-muted">
                            Set the maximum number of colours representing competing ideas in the network
                        </Form.Text>
                    </Form.Group>
                </Modal.Body>
                <Modal.Footer>
                    <Button variant="success" >Save</Button>
                    <Button variant="secondary" onClick={handleoptclose}>Cancel</Button>
                </Modal.Footer>
            </Modal>
        </Container>
    )
}

export default Simulation