import React from 'react';
import {Link} from 'react-router-dom';
import {Card} from 'react-bootstrap';
import {Col} from 'react-bootstrap';
import {Trash} from 'react-bootstrap-icons';
import { SimInfo } from '../API/SimInfo';


type SimListProps = {
    deleteFunc(id: string, name: string): void;
    sims: Array<SimInfo>;
}

const SimList = (props:SimListProps) => {

    return props.sims.map(sim => {
            return(
                <Col className="mb-4">
                    <Card border="info" className="h-100" style={{ width: '20rem' }} key={sim.id}>
                        <Card.Header><Card.Title>{sim.name}<Trash className="float-right" onClick={() => props.deleteFunc(sim.id, sim.name)}/></Card.Title></Card.Header>
                        <Card.Body>
                            <Card.Text>{sim.description}</Card.Text>
                        </Card.Body>
                        <Card.Footer>
                            <Link className="btn btn-primary" role="button" to={'/simulation/' + sim.id}>Open</Link>
                        </Card.Footer>
                    </Card>
                </Col>
            );
        });
}

export default SimList