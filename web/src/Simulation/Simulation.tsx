import React, { useState, useEffect } from 'react';
import API from '../API/api';
import {Link} from 'react-router-dom';
import {ArrowLeftCircle} from 'react-bootstrap-icons';
import {Card, Button, Row, Col, Container} from 'react-bootstrap'
import EditNameDescModal from '../Home/EditNameDescModal'
import NetworkOptionsCard from './NetworkOptionsCard'
import NetworkCard from './NetworkCard'
import StepsCard from './StepsCard'
import { SimInfo } from '../API/SimInfo';
import { Step } from '../API/Step';

type SimulationProps = {
    match: {
        params: {
            id: string;
        }
    }
}

const Simulation = (props:SimulationProps) => {
    const [sim, setSim] = useState<SimInfo>(API.emptySim());
    const [simsteps, setSimsteps] = useState<Array<Step>>([]);
    const [showsimeditmodal, setShowsimeditmodal] = useState<boolean>(false);

    useEffect(() => {
        readsim(props.match.params.id);
      }, [props.match.params.id]);

    const handlesimeditshow = () => setShowsimeditmodal(true);
    const handlesimeditclose = () => setShowsimeditmodal(false);

    const readsim = (id:string) => {
        API.get(id).then(sresp => {
            setSim(sresp);
            API.getSteps(sresp).then(steps => {
                setSimsteps(steps);
            });
        })
    };
    
    const updatesim = (simtosave:SimInfo) => {
        API.update(simtosave).then(response => {
            setSim(response);
        })
    };

    return (
        <Container fluid>
            <h1 className="px-2 py-2 bg-light rounded"><Link className="btn btn-outline-secondary mr-3 mt-n2 mb-2" to='/' role="button"><ArrowLeftCircle className="mt-n1" /></Link>{sim.name}<Button size="sm" className="btn btn-primary mt-2 mr-2 float-right" onClick={handlesimeditshow}>Edit</Button></h1>
            <Container fluid>
                <Row>
                    <Col sm={8}>
                        <NetworkCard sim={sim} steps={simsteps} readsim={readsim}/>
                        <StepsCard sim={sim} steps={simsteps} readsim={readsim}/>
                    </Col>
                    <Col sm={4}>
                        <Card className="mb-2 mx-n2">
                            <Card.Header><Card.Title>Description</Card.Title></Card.Header>
                            <Card.Body><Card.Text>{sim.description}</Card.Text></Card.Body>
                        </Card>
                        <NetworkOptionsCard sim={sim} setsim={setSim}/>
                    </Col>
                </Row>
            </Container>
            <EditNameDescModal sim={sim} show={showsimeditmodal} saveFunc={updatesim} closeFunc={handlesimeditclose}/>
        </Container>
    )
}

export default Simulation